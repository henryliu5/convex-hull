package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Worker to process a subhull and send back to manager
func subhull_worker(work_ch chan [][2]float32, result_ch chan [][2]float32, wg *sync.WaitGroup) {
	for points := range work_ch {
		result_ch <- parallel_graham_scan(points)
	}
	wg.Done()
}

// Parallel subhull comptuation that uses a thread-pool like thing to limit the number of goroutines to MAX_SUBHULL_WORKERS
func thread_pool_subhull(points [][2]float32, group_size, n int, subhulls [][2]float32, subhull_sizes []int) ([][2]float32, []int) {
	MAX_SUBHULL_WORKERS := 300

	var subhull_compute time.Duration
	var subhull_append time.Duration

	work_ch := make(chan [][2]float32)
	result_ch := make(chan [][2]float32)
	wg := sync.WaitGroup{}

	num_workers := (n + group_size - 1) / group_size
	if num_workers > MAX_SUBHULL_WORKERS {
		num_workers = MAX_SUBHULL_WORKERS
	}

	for i := 0; i < num_workers; i++ {
		wg.Add(1)
		go subhull_worker(work_ch, result_ch, &wg)
	}

	num_subhulls := (n + group_size - 1) / group_size

	manager_wg := sync.WaitGroup{}
	manager_wg.Add(1)
	manager := func() {
		// Central manager to collect subhulls from workers
		for i := 0; i < num_subhulls; i++ {
			subhull := <-result_ch
			// for subhull := range result_ch {
			subhull_start := time.Now()
			subhull_compute += time.Since(subhull_start)

			append_start := time.Now()
			// Aggregrate into a single array for performance
			subhulls = append(subhulls, subhull...)
			// but track sizes to know offset
			subhull_sizes = append(subhull_sizes, len(subhull))
			subhull_append += time.Since(append_start)
		}
		manager_wg.Done()
	}

	go manager()

	// Run graham scan on subgroups of points
	for start := 0; start < n; start += group_size {
		end := start + group_size
		if n < end {
			end = n
		}
		// Compute convex hull of subgroup using graham-scan
		work_ch <- points[start:end]
	}

	// Finished sending work
	close(work_ch)
	// Wait for workers to complete
	wg.Wait()
	manager_wg.Wait()

	// fmt.Println("-------", num_subhulls, (n+group_size-1)/group_size)
	debug("subhull point count", len(subhulls))
	debug("chan's subgroups compute", subhull_compute)
	debug("chan's subgroups append", subhull_append)
	return subhulls, subhull_sizes
}

// Parallel subhull computation that spawns 1 goroutine for each subhull
func basic_par_subhull(points [][2]float32, group_size, n int, subhulls [][2]float32, subhull_sizes []int) ([][2]float32, []int) {
	var subhull_compute time.Duration
	var subhull_append time.Duration
	num_subhulls := 0
	ch := make(chan [][2]float32)
	worker := func(points [][2]float32, ch chan [][2]float32) {
		ch <- parallel_graham_scan(points)
	}
	// Run graham scan on subgroups of points
	for start := 0; start < n; start += group_size {
		end := start + group_size
		if n < end {
			end = n
		}
		// Compute convex hull of subgroup using graham-scan
		go worker(points[start:end], ch)
		num_subhulls += 1
	}
	// Central manager to collect subhulls from workers
	for i := 0; i < num_subhulls; i++ {
		subhull_start := time.Now()
		subhull := <-ch
		subhull_compute += time.Since(subhull_start)

		append_start := time.Now()
		// Aggregrate into a single array for performance
		subhulls = append(subhulls, subhull...)
		// but track sizes to know offset
		subhull_sizes = append(subhull_sizes, len(subhull))
		subhull_append += time.Since(append_start)
	}

	// debug("-------", num_subhulls)
	debug("subhull point count", len(subhulls))
	debug("chan's subgroups compute", subhull_compute)
	debug("chan's subgroups append", subhull_append)

	return subhulls, subhull_sizes
}

// Parallel Jarvis march to wrap subhulls
// 	Goes both left and right from start position in parallel
// 	Tracks global state w atomic counter and concurrent hash map
func parallel_subhull_jarvis(points [][2]float32, subhull_sizes []int, group_size int) [][2]float32 {
	var hull [][2]float32
	n_subhulls := len(subhull_sizes)

	// Use a concurrent hashmap to track which points have been selected
	selected := SafeMap{}
	selected.init(len(points) / 4)

	// Leftmost point starts as first point on hull
	left := leftmost(points)
	right := rightmost(points)

	wg := sync.WaitGroup{}
	failed := false

	// Track total steps taken by all workers with atomic counter
	var steps uint64

	// Worker to wrap subhulls from start point
	find_hull := func(start int, order float32) {
		candidates := make([][2]float32, n_subhulls)
		cur_p := points[start]
		// If need to add more points than the group size, retry with different group size
		for {
			// fmt.Println("cur", order, cur_p)
			selected.put(cur_p, [][2]float32{[2]float32{1.0, 1.0}})
			subhull_index := 0
			// Compute the tangent point for each of the subhulls
			for i := 0; i < n_subhulls; i++ {
				start := subhull_index
				end := subhull_index + subhull_sizes[i]
				// candidates[i] = find_tangent(points[start:end], cur_p, order)
				candidates[i] = find_tangent_bsearch(points[start:end], cur_p, order)
				subhull_index += subhull_sizes[i]
			}

			// Find leftmost endpoint out of all of the subhulls
			endpoint := 0
			leftmost_search_start := time.Now()
			for candidate := range candidates {
				cross := cross_prod(cur_p, candidates[endpoint], candidates[candidate]) * order
				if (candidates[endpoint][0] == cur_p[0] && candidates[endpoint][1] == cur_p[1]) || cross > 0 {
					// New point is to the left of current endpoint
					endpoint = candidate
				} else if cross == 0 && dist(cur_p, candidates[candidate]) >= dist(cur_p, candidates[endpoint]) {
					if dist(cur_p, points[candidate]) == dist(cur_p, points[endpoint]) && endpoint != left {
						// New point is collinear but further than current endpoint
						endpoint = candidate
					}
				}
			}
			leftmost_time += time.Since(leftmost_search_start)

			cur_p = candidates[endpoint]

			// Someone already found this point
			if selected.get(cur_p) != nil {
				failed = failed || false
				break
			}
			if steps == uint64(group_size-1) {
				failed = true
				break
			}
			atomic.AddUint64(&steps, 1)
		}
		wg.Done()
	}

	count := 0
	for i := 0; i < len(points); i++ {
		if selected.get(points[i]) != nil {
			count += 1
		}
	}
	// Now traverse the subhulls both left and right
	wg.Add(2)
	go find_hull(left, 1.0)
	go find_hull(right, 1.0)
	wg.Wait()

	if failed {
		// need to retry
		return nil
	}

	for i := 0; i < len(points); i++ {
		if selected.get(points[i]) != nil {
			hull = append(hull, points[i])
		}
	}
	return hull
}

// Enable subhull coaslecing from previous iterations
var USE_COALESCE bool
var SIMUL_ITERS int

// Parallel Chan's algorithm O(nlogh)
func parallel_chans(points [][2]float32) [][2]float32 {
	n := len(points)
	tangent_time = 0

	if USE_COALESCE {
		// Initialize size of concurrent map
		global_subhulls.init(n / 4)
	}

	var t uint
	t = 3 // Init group size as 2^2^3 = 256
	start := time.Now()

	debug("initial allocation time", time.Since(start))

	ch := make(chan [][2]float32)
	for {
		iteration := func(cur_t uint) {
			subhulls := make([][2]float32, 0, n>>1) // n>>1 Just a guess on how many points will be in subhulls
			subhull_sizes := make([]int, 0, n/(1<<(1<<t))+1)
			// Try out group size
			group_size := 1 << (1 << cur_t)
			if group_size == 0 {
				fmt.Println("chan's failed, too many iterations")
				ch <- nil
			}
			debug("current group size", group_size)
			if n < group_size {
				group_size = n
			}

			/***********************
			 * Subhull computation *
			 ***********************/
			subhull_start := time.Now()

			if USE_COALESCE {
				subhulls, subhull_sizes = coalesce_subhull(points, group_size, n, subhulls, subhull_sizes)
				// subhulls, subhull_sizes = thread_pool_coalesce_subhull(points, group_size, n, subhulls, subhull_sizes)
				debug("POINTS SAVED", points_saved)
			} else {
				copy_points := make([][2]float32, len(points))
				copy(copy_points, points)
				// subhulls, subhull_sizes = basic_par_subhull(points, group_size, n, subhulls, subhull_sizes)
				subhulls, subhull_sizes = thread_pool_subhull(copy_points, group_size, n, subhulls, subhull_sizes)
			}
			debug("chan's subgroups total", time.Since(subhull_start))

			/****************************************
			 * Jarvis March (gift wrap) of subhulls *
			 ****************************************/
			march_start := time.Now()
			var hull [][2]float32
			// Jarvis march meant for Chan's algorithm, operates on subhulls
			// hull = subhull_jarvis(subhulls, subhull_sizes, group_size)
			hull = parallel_subhull_jarvis(subhulls, subhull_sizes, group_size)
			debug("march time", time.Since(march_start))

			ch <- hull
		}

		// Conduct multiple iterations at once
		for i := 0; i < SIMUL_ITERS; i++ {
			go iteration(t + uint(i))
		}

		// See if any of the iterations were succesful
		for i := 0; i < SIMUL_ITERS; i++ {
			hull := <-ch
			if hull != nil {
				debug("tangent time", tangent_time)
				debug("leftmost time", leftmost_time)
				return hull
			}
		}

		// Group size too small
		t += uint(SIMUL_ITERS)
	}
}

/********************
 * Coalescing stuff *
 ********************/

var global_subhulls SafeMap
var points_saved int

type Subhull struct {
	points [][2]float32
	start  int
	end    int
}

// Parallel subhull computation that coalesces previously computed subhulls as well
func coalesce_subhull(points [][2]float32, group_size, n int, subhulls [][2]float32, subhull_sizes []int) ([][2]float32, []int) {
	var subhull_compute time.Duration
	var subhull_append time.Duration
	num_subhulls := 0
	ch := make(chan Subhull)
	worker := func(points [][2]float32, start, end int, ch chan Subhull) {
		res := parallel_graham_scan(points)
		ch <- Subhull{res, start, end}
	}
	// Run graham scan on subgroups of points
	for start := 0; start < n; start += group_size {
		end := start + group_size
		if n < end {
			end = n
		}

		// Try to build subhull using old results
		local_points := make([][2]float32, 0)
		if group_size > 256 {
			// See if a smaller iteration already computed the subhull
			for local_start := start; local_start < end; local_start += 256 {
				local_end := local_start + 256
				if n < local_end {
					local_end = n
				}
				previous_result := global_subhulls.get([2]float32{float32(local_start), float32(local_end)})
				if previous_result != nil {
					local_points = append(local_points, previous_result...)
					points_saved += (local_end - local_start) - len(previous_result)
				} else {
					local_points = append(local_points, points[local_start:local_end]...)
				}
			}
		} else {
			local_points = points[start:end]
		}

		// Compute convex hull of subgroup using graham-scan
		go worker(local_points, start, end, ch)
		num_subhulls += 1
	}
	// Central manager to collect subhulls from workers
	for i := 0; i < num_subhulls; i++ {
		subhull_start := time.Now()
		res := <-ch
		subhull := res.points
		subhull_compute += time.Since(subhull_start)

		// Update the result so future iterations can use
		global_subhulls.put([2]float32{float32(res.start), float32(res.end)}, res.points)

		append_start := time.Now()
		// Aggregrate into a single array for performance
		subhulls = append(subhulls, subhull...)
		// but track sizes to know offset
		subhull_sizes = append(subhull_sizes, len(subhull))
		subhull_append += time.Since(append_start)

	}

	// debug("-------", num_subhulls)
	debug("subhull point count", len(subhulls))
	debug("chan's subgroups compute", subhull_compute)
	debug("chan's subgroups append", subhull_append)

	return subhulls, subhull_sizes
}

// // Worker to process a subhull and send back to manager
// func coalesce_worker(work_ch chan Subhull, result_ch chan Subhull, wg *sync.WaitGroup) {
// 	for subhull := range work_ch {
// 		res := parallel_graham_scan(subhull.points)
// 		result_ch <- Subhull{res, subhull.start, subhull.end}
// 	}
// 	wg.Done()
// }

// // Parallel subhull comptuation that that coalesces previously computed subhulls as well, uses a thread-pool like thing to limit the number of goroutines to MAX_SUBHULL_WORKERS
// func thread_pool_coalesce_subhull(points [][2]float32, group_size, n int, subhulls [][2]float32, subhull_sizes []int) ([][2]float32, []int) {
// 	MAX_SUBHULL_WORKERS := 300

// 	var subhull_compute time.Duration
// 	var subhull_append time.Duration

// 	work_ch := make(chan Subhull)
// 	result_ch := make(chan Subhull)
// 	wg := sync.WaitGroup{}

// 	num_workers := (n + group_size - 1) / group_size
// 	if num_workers > MAX_SUBHULL_WORKERS {
// 		num_workers = MAX_SUBHULL_WORKERS
// 	}

// 	for i := 0; i < num_workers; i++ {
// 		wg.Add(1)
// 		go coalesce_worker(work_ch, result_ch, &wg)
// 	}

// 	num_subhulls := (n + group_size - 1) / group_size

// 	manager_wg := sync.WaitGroup{}
// 	manager_wg.Add(1)
// 	manager := func() {
// 		// Central manager to collect subhulls from workers
// 		for i := 0; i < num_subhulls; i++ {
// 			res := <-result_ch
// 			subhull := res.points
// 			// for subhull := range result_ch {
// 			subhull_start := time.Now()
// 			subhull_compute += time.Since(subhull_start)

// 			// Update the result so future iterations can use
// 			global_subhulls.put([2]float32{float32(res.start), float32(res.end)}, res.points)

// 			append_start := time.Now()
// 			// Aggregrate into a single array for performance
// 			subhulls = append(subhulls, subhull...)
// 			// but track sizes to know offset
// 			subhull_sizes = append(subhull_sizes, len(subhull))
// 			subhull_append += time.Since(append_start)
// 		}
// 		manager_wg.Done()
// 	}

// 	go manager()

// 	// Run graham scan on subgroups of points
// 	for start := 0; start < n; start += group_size {
// 		end := start + group_size
// 		if n < end {
// 			end = n
// 		}

// 		// Try to build subhull using old results
// 		local_points := make([][2]float32, 0)
// 		if group_size > 256 {
// 			// See if a smaller iteration already computed the subhull
// 			for local_start := start; local_start < end; local_start += 256 {
// 				local_end := local_start + 256
// 				if n < local_end {
// 					local_end = n
// 				}
// 				previous_result := global_subhulls.get([2]float32{float32(local_start), float32(local_end)})
// 				if previous_result != nil {
// 					local_points = append(local_points, previous_result...)
// 					points_saved += (local_end - local_start) - len(previous_result)
// 				} else {
// 					local_points = append(local_points, points[local_start:local_end]...)
// 				}
// 			}

// 		} else {
// 			local_points = points[start:end]
// 		}
// 		// Compute convex hull of subgroup using graham-scan
// 		work_ch <- Subhull{local_points, start, end}
// 	}

// 	// Finished sending work
// 	close(work_ch)
// 	// Wait for workers to complete
// 	wg.Wait()
// 	manager_wg.Wait()

// 	// fmt.Println("-------", num_subhulls, (n+group_size-1)/group_size)
// 	debug("subhull point count", len(subhulls))
// 	debug("chan's subgroups compute", subhull_compute)
// 	debug("chan's subgroups append", subhull_append)
// 	return subhulls, subhull_sizes
// }

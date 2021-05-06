package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Parallel subhull computation that uses a thread-pool like thing to limit the number of goroutines to MAX_SUBHULL_WORKERS
func thread_pool_subhull(points [][2]float32, group_size, n int, subhulls [][2]float32, subhull_sizes []int) ([][2]float32, []int) {
	const MAX_SUBHULL_WORKERS int = 300

	var subhull_compute time.Duration
	var subhull_append time.Duration

	// Sync primitives to manage worker/manager communication + lifecycle
	work_ch := make(chan [][2]float32)
	result_ch := make(chan [][2]float32)
	wg := sync.WaitGroup{}

	// Worker function to process a subhull and send back to manager
	subhull_worker := func() {
		defer wg.Done()
		for points := range work_ch {
			result_ch <- parallel_graham_scan(points)
		}
	}

	// Try to have 1 worker goroutine per subhull, or cap to MAX_SUBHULL_WORKERS
	num_workers := (n + group_size - 1) / group_size
	if num_workers > MAX_SUBHULL_WORKERS {
		num_workers = MAX_SUBHULL_WORKERS
	}

	for i := 0; i < num_workers; i++ {
		wg.Add(1)
		go subhull_worker()
	}

	// Create and launch manager goroutine to coalesce worker output
	manager_wg := sync.WaitGroup{}
	manager_wg.Add(1)
	manager := func() {
		defer manager_wg.Done()
		// Central manager to collect subhulls from workers
		for i := 0; i < num_workers; i++ {
			subhull_start := time.Now()
			subhull := <-result_ch
			subhull_compute += time.Since(subhull_start)

			append_start := time.Now()
			// Aggregate into a single slice
			subhulls = append(subhulls, subhull...)
			// but track sizes to know offset/index in later
			subhull_sizes = append(subhull_sizes, len(subhull))
			subhull_append += time.Since(append_start)
		}
	}

	// Launch manager before using this thread to send data
	// (manager can start aggregating before all work is sent)
	go manager()

	// Run graham scan on subgroups of points
	// Send groups of points to worker goroutines
	for start := 0; start < n; start += group_size {
		end := start + group_size
		if n < end {
			end = n
		}
		work_ch <- points[start:end]
	}

	// Finished sending work
	close(work_ch)
	// Wait for workers to complete
	wg.Wait()
	manager_wg.Wait()

	// Report compute metrics
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
		// Aggregate into a single array for performance
		subhulls = append(subhulls, subhull...)
		// but track sizes to know offset
		subhull_sizes = append(subhull_sizes, len(subhull))
		subhull_append += time.Since(append_start)
	}

	debug("subhull point count", len(subhulls))
	debug("chan's subgroups compute", subhull_compute)
	debug("chan's subgroups append", subhull_append)

	return subhulls, subhull_sizes
}

// Parallel Jarvis march to wrap subhulls
// 	Goes both left and right from start position in parallel
// 	Tracks global state w atomic counter and concurrent hash map
//  Note that subhull_points is points on subhulls concatenated together, can be split with subhull_sizes
//  i.e. subhull_points = [subhull1_point1, subhull1_point2, ..., subhull2_point1, ...  subhullN_pointK]
func parallel_subhull_jarvis(subhull_points [][2]float32, subhull_sizes []int, group_size int) [][2]float32 {
	var hull [][2]float32
	n_subhulls := len(subhull_sizes)

	// Use a concurrent hashmap to track which points have been selected (used as a set here)
	selected := SafeMap{}
	selected.init(len(subhull_points) / 4)

	// Leftmost point starts as first point on hull
	left := leftmost(subhull_points)
	right := rightmost(subhull_points)

	wg := sync.WaitGroup{}
	failed := false

	// Track total steps taken by all workers with atomic counter
	var steps uint64

	// Worker to wrap subhulls from start point
	find_hull := func(start int, order float32) {
		defer wg.Done()
		tangents := make([][2]float32, n_subhulls)
		cur_p := subhull_points[start]
		// If need to add more points than the group size, retry with different group size
		for {
			selected.put(cur_p, [][2]float32{[2]float32{1.0, 1.0}})
			subhull_index := 0
			// Compute the tangent point for each of the subhulls
			for i := 0; i < n_subhulls; i++ {
				start := subhull_index
				end := subhull_index + subhull_sizes[i]
				tangents[i] = find_tangent_bsearch(subhull_points[start:end], cur_p, order)
				subhull_index += subhull_sizes[i]
			}

			// Find leftmost endpoint out of all of the subhulls
			endpoint := 0
			for tangent := range tangents {
				cross := cross_prod(cur_p, tangents[endpoint], tangents[tangent]) * order
				if (tangents[endpoint][0] == cur_p[0] && tangents[endpoint][1] == cur_p[1]) || cross > 0 {
					// New point is to the left of current endpoint
					endpoint = tangent
				} else if cross == 0 && dist(cur_p, tangents[tangent]) >= dist(cur_p, tangents[endpoint]) {
					if dist(cur_p, subhull_points[tangent]) == dist(cur_p, subhull_points[endpoint]) && endpoint != left {
						// New point is collinear but further than current endpoint
						endpoint = tangent
					}
				}
			}

			cur_p = tangents[endpoint]

			// Someone already found this point
			if selected.get(cur_p) != nil {
				break
			}
			if steps == uint64(group_size-1) {
				failed = true
				break
			}
			atomic.AddUint64(&steps, 1)
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

	// Get hull points from map
	for i := 0; i < len(subhull_points); i++ {
		if selected.get(subhull_points[i]) != nil {
			hull = append(hull, subhull_points[i])
		}
	}
	return hull
}

// Enable subhull coalescing from previous iterations
var USE_COALESCE bool

// Number of iterations to run in parallel
var SIMUL_ITERS int

// Controls initial group size (2^2^3) = 256
const INITIAL_T uint = 3

// Parallel Chan's algorithm O(nlogh)
func parallel_chans(points [][2]float32) [][2]float32 {
	n := len(points)
	var global_subhulls SafeMap
	if USE_COALESCE {
		// Initialize size of concurrent map
		global_subhulls.init(n / 4)
	}
	// Initialize group size as 2^2^3 = 256
	var t uint = INITIAL_T

	// Channel to receive results from each iteration
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

			/*************************************
			 * Subhull Computation - Graham Scan *
			 *************************************/
			subhull_start := time.Now()

			if USE_COALESCE {
				subhulls, subhull_sizes = coalesce_subhull(&global_subhulls, points, group_size, n, subhulls, subhull_sizes)
				debug("coalesce points saved", points_saved)
			} else {
				// Graham scan modifies in place so need to be careful to copy as input slice will be modified
				copy_points := make([][2]float32, len(points))
				copy(copy_points, points)
				subhulls, subhull_sizes = thread_pool_subhull(copy_points, group_size, n, subhulls, subhull_sizes)
			}
			debug("chan's subgroups total", time.Since(subhull_start))

			/****************************************
			 * Jarvis March (gift wrap) of subhulls *
			 ****************************************/
			march_start := time.Now()
			var hull [][2]float32
			// Jarvis march meant for Chan's algorithm, operates on subhulls
			hull = parallel_subhull_jarvis(subhulls, subhull_sizes, group_size)
			debug("march time", time.Since(march_start))

			ch <- hull
		}

		// Conduct multiple iterations at once
		for i := 0; i < SIMUL_ITERS; i++ {
			go iteration(t + uint(i))
		}

		// See if any of the iterations were successful
		for i := 0; i < SIMUL_ITERS; i++ {
			hull := <-ch
			if hull != nil {
				return hull
			}
		}

		// Group size too small
		t += uint(SIMUL_ITERS)
	}
}

/*********************************************************
 * Coalescing Subhulls Between Iterations / Work Sharing *
 *********************************************************
When performing multiple iterations (different estimates for # of hull points) in parallel,
we can attempt to share work between these iterations by creating a mapping between ranges
in the original "points" array and the points that belong on the subhull corresponding to that range.

The intuition is if points are not in the subhull of points[i:j], then they will not be in the
subhull of points[a:b], where a <= i <= j <= b. Thus we can prune based on calculations from
iterations with smaller subhull sizes.
*/

var points_saved int

// Represents convex hull of range [start:end] in original input points
type Subhull struct {
	// points on convex hull of this range
	points [][2]float32
	start  int
	end    int
}

// Parallel subhull computation that coalesces previously computed subhulls as well
func coalesce_subhull(global_subhulls *SafeMap, points [][2]float32, group_size, n int, subhulls [][2]float32, subhull_sizes []int) ([][2]float32, []int) {
	num_subhulls := 0
	// Channel to collect subhulls with metadata (range info) from workers
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

		// local_points: Points to consider for this iteration, may consist of just points from
		// original input, or points known to be on convex hull of some subset of the pointsv
		local_points := make([][2]float32, 0)
		smallest_group_size := (int)((INITIAL_T >> 1) >> 1) // 256 if t=3

		// Try to build subhull using old results (from first iteration, from empirical results other
		// iterations are not as helpful, likely because results take longer to update and will be "missed"
		if group_size > smallest_group_size {
			// See if a smaller iteration already computed any subhulls for this range
			for local_start := start; local_start < end; local_start += smallest_group_size {
				local_end := local_start + smallest_group_size
				if n < local_end {
					local_end = n
				}
				// Check if map contains this range
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

	var subhull_compute time.Duration
	var subhull_append time.Duration

	// Collect subhulls from workers
	for i := 0; i < num_subhulls; i++ {
		subhull_start := time.Now()
		res := <-ch
		subhull := res.points
		subhull_compute += time.Since(subhull_start)

		// Update the result so future iterations can use
		global_subhulls.put([2]float32{float32(res.start), float32(res.end)}, res.points)

		append_start := time.Now()
		// Aggregate into a single array for performance
		subhulls = append(subhulls, subhull...)
		// but track sizes to know offset
		subhull_sizes = append(subhull_sizes, len(subhull))
		subhull_append += time.Since(append_start)

	}

	// Metrics logging
	debug("subhull point count", len(subhulls))
	debug("chan's subgroups compute", subhull_compute)
	debug("chan's subgroups append", subhull_append)

	return subhulls, subhull_sizes
}

package main

import (
	"fmt"
	"sync"
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
		// go subhull_worker(points[start:end], ch)
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

// Parallel Chan's algorithm O(nlogh)
func parallel_chans(points [][2]float32) [][2]float32 {
	n := len(points)

	// TODO try running different iterations in parallel
	var t uint
	t = 3 // Init group size as 2^2^3 = 256
	start := time.Now()
	subhulls := make([][2]float32, 0, n>>1) // n>>1 Just a guess on how many points will be in subhulls
	subhull_sizes := make([]int, 0, n/(1<<(1<<t))+1)
	debug("initial allocation time", time.Since(start))

	ch := make(chan [][2]float32)
	for {
		iteration := func(cur_t uint) {
			subhulls = subhulls[:0]
			subhull_sizes = subhull_sizes[:0]

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
			// subhulls, subhull_sizes = basic_par_subhull(points, group_size, n, subhulls, subhull_sizes)
			subhulls, subhull_sizes = thread_pool_subhull(points, group_size, n, subhulls, subhull_sizes)
			debug("chan's subgroups total", time.Since(subhull_start))

			/****************************************
			 * Jarvis March (gift wrap) of subhulls *
			 ****************************************/
			march_start := time.Now()
			var hull [][2]float32
			// Jarvis march meant for Chan's algorithm, should use bsearch
			// TODO parallelize this somehow
			hull = subhull_jarvis(subhulls, subhull_sizes, group_size)
			debug("march time", time.Since(march_start))

			ch <- hull
		}

		// Conduct multiple iterations at once
		simul_iters := 2
		for i := 0; i < simul_iters; i++ {
			go iteration(t + uint(i))
		}

		// See if any of the iterations were succesful
		for i := 0; i < simul_iters; i++ {
			hull := <-ch
			if hull != nil {
				return hull
			}
		}

		// Group size too small
		t += uint(simul_iters)
	}
}

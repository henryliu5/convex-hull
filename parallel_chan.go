package main

import (
	"fmt"
	"time"
)

// TODO make this properly limited.... thread pool?
func subhull_worker(points [][2]float32, ch chan [][2]float32) {
	ch <- parallel_graham_scan(points)
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

	for {
		subhulls = subhulls[:0]
		subhull_sizes = subhull_sizes[:0]

		// Try out group size
		group_size := 1 << (1 << t)
		if group_size == 0 {
			fmt.Println("chan's failed, too many iterations")
			return nil
		}
		debug("current group size", group_size)
		if n < group_size {
			group_size = n
		}

		/***********************
		 * Subhull computation *
		 ***********************/
		subhull_start := time.Now()
		var subhull_compute time.Duration
		var subhull_append time.Duration

		ch := make(chan [][2]float32)
		num_subhulls := 0
		// Run graham scan on subgroups of points
		for start := 0; start < n; start += group_size {
			end := start + group_size
			if n < end {
				end = n
			}
			// Compute convex hull of subgroup using graham-scan
			go subhull_worker(points[start:end], ch)
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
		debug("-------", num_subhulls)
		debug("subhull point count", len(subhulls))
		debug("chan's subgroups compute", subhull_compute)
		debug("chan's subgroups append", subhull_append)
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

		if hull != nil {
			return hull
		}
		// Group size too small
		t++
	}
}

package main

import (
	"time"
)

// Find point on this convex subhull that is as left as possible from point p (p must not be inside the subhull)
func basic_tangent(subhull [][]float32, p []float32) []float32 {
	endpoint := 0
	// Look through this subhull
	for i := 1; i < len(subhull); i++ {
		cross := cross_prod(p, subhull[endpoint], subhull[i])
		if (subhull[endpoint][0] == p[0] && subhull[endpoint][1] == p[1]) || cross > 0 {
			// New point is to the left of current endpoint
			endpoint = i
		} else if cross == 0 && dist(p, subhull[i]) > dist(p, subhull[endpoint]) {
			// New point is collinear but further than current endpoint
			endpoint = i
		}
	}
	return subhull[endpoint]
}

// TODO Use binary search to find tangent point of a subhull instead of using O(n) scan above ^^
// p: reference point for tangent
func binary_search_tangent(subhull [][]float32, p []float32) int {
	return -1
}

// Jarvis march on subhulls
func subhull_jarvis(points [][]float32, subhull_sizes []int, group_size int) [][]float32 {

	var hull [][]float32
	n_subhulls := len(subhull_sizes)
	candidates := make([][]float32, n_subhulls)

	// Leftmost point starts as first point on hull
	left := leftmost(points)
	cur_p := points[left]

	// If need to add more points than the group size, retry with different group size
	for step := 0; step < group_size; step++ {
		hull = append(hull, cur_p)
		subhull_index := 0
		// Compute the tangent point for each of the subhulls
		for i := 0; i < n_subhulls; i++ {
			start := subhull_index
			end := subhull_index + subhull_sizes[i]
			// TODO replace this with binary search for tangent
			candidates[i] = basic_tangent(points[start:end], cur_p)
			subhull_index += subhull_sizes[i]
		}

		// Find leftmost endpoint out of all of the subhulls
		endpoint := 0
		for candidate := range candidates {
			cross := cross_prod(cur_p, candidates[endpoint], candidates[candidate])
			if (candidates[endpoint][0] == cur_p[0] && candidates[endpoint][1] == cur_p[1]) || cross > 0 {
				// New point is to the left of current endpoint
				endpoint = candidate
			} else if cross == 0 && dist(cur_p, candidates[candidate]) > dist(cur_p, candidates[endpoint]) {
				// New point is collinear but further than current endpoint
				endpoint = candidate
			}
		}
		cur_p = candidates[endpoint]

		// Circled back to original point
		if cur_p[0] == points[left][0] && cur_p[1] == points[left][1] {
			return hull
		}
	}
	// Need to retry
	return nil
}

// Sequential Chan's algorithm O(nlogh)
func seq_chans(points [][]float32) [][]float32 {
	n := len(points)

	var t uint
	t = 3 // Init group size as 2^2^3 = 256
	start := time.Now()
	subhulls := make([][]float32, 0, n>>1) // n>>1 Just a guess on how many points will be in subhulls
	subhull_sizes := make([]int, 0, n/(1<<(1<<t))+1)
	debug("initial allocation time", time.Since(start))

	for {
		subhulls = subhulls[:0]
		subhull_sizes = subhull_sizes[:0]

		// Try out group size
		group_size := 1 << (1 << t)
		debug("current group size", group_size)
		if n < group_size {
			group_size = n
		}

		/***********************
		 * Subhull computation *
		 ***********************/

		var subhull_compute time.Duration
		var subhull_append time.Duration
		// Run graham scan on subgroups of points
		for start := 0; start < n; start += group_size {
			end := start + group_size
			if n < end {
				end = n
			}
			// Compute convex hull of subgroup
			// TODO change to graham scan
			subhull_start := time.Now()
			subhull := seq_jarvis(points[start:end])
			subhull_compute += time.Since(subhull_start)

			// Add subhull points
			append_start := time.Now()
			subhulls = append(subhulls, subhull...)
			subhull_sizes = append(subhull_sizes, len(subhull))
			subhull_append += time.Since(append_start)
		}
		debug("subhull point count", len(subhulls))
		debug("chan's subgroups compute", subhull_compute)
		debug("chan's subgroups append", subhull_append)
		debug("chan's subgroups total", subhull_append+subhull_compute)

		/****************************************
		 * Jarvis March (gift wrap) of subhulls *
		 ****************************************/
		march_start := time.Now()
		var hull [][]float32
		// Basic way, worse than O(nlogh)
		// hull = seq_jarvis(subhulls)
		// Jarvis march meant for Chan's algorithm, should use bsearch
		hull = subhull_jarvis(subhulls, subhull_sizes, group_size)
		debug("march time", time.Since(march_start))

		if hull != nil {
			return hull
		}
		// Group size too small
		t++
	}
}

package main

import (
	"fmt"
	"sync"
	"time"
)

// Find leftmost point relative to p
func find_single_left(points [][2]float32, p int, left_map []int, wg *sync.WaitGroup) {
	endpoint := 0
	for candidate := range points {
		cross := cross_prod(points[p], points[endpoint], points[candidate])
		if endpoint == p || cross > 0 {
			// New point is to the left of current endpoint
			endpoint = candidate
		} else if cross == 0 && dist(points[p], points[candidate]) >= dist(points[p], points[endpoint]) {
			// New point is collinear but further than current endpoint
			endpoint = candidate
		}
	}
	left_map[p] = endpoint
	fmt.Println(p)
	wg.Done()
}

// Naively precompute all leftmost points, O(N^2) no bueno
func find_all_lefts(points [][2]float32) []int {
	left_map := make([]int, len(points))
	wg := sync.WaitGroup{}
	for i := 0; i < len(points); i++ {
		wg.Add(1)
		go find_single_left(points, i, left_map, &wg)
	}
	wg.Wait()
	return left_map
}

// Parallel Jarvis March
func naive_parallel_jarvis(points [][2]float32) [][2]float32 {
	if len(points) < 3 {
		return points
	}
	var hull [][2]float32
	left := leftmost(points)
	// Last selected point on hull
	p := left

	left_map := find_all_lefts(points)
	for {
		hull = append(hull, points[p])
		// Find leftmost endpoint
		p = left_map[p]
		// Circled back to original point
		if p == left {
			break
		}
	}
	return hull
}

// Parallel Jarvis March - searches in parallel both left and right starting at multiple points on hull
func parallel_jarvis(points [][2]float32) [][2]float32 {
	if len(points) < 3 {
		return points
	}
	selected := make([]bool, len(points))
	var hull [][2]float32

	find_extremes_start := time.Now()
	left := leftmost(points)
	right := rightmost(points)
	down := lowest(points)
	up := highest(points)
	debug("parallel jarvis find extremes", time.Since(find_extremes_start))

	wg := sync.WaitGroup{}
	do_left := func(start int) {
		left_p := start
		for {
			selected[left_p] = true
			// Find leftmost endpoint
			endpoint := 0
			for candidate := range points {
				cross := cross_prod(points[left_p], points[endpoint], points[candidate])
				if endpoint == left_p || cross < 0 {
					// New point is to the left of current endpoint
					endpoint = candidate
				} else if cross == 0 && dist(points[left_p], points[candidate]) >= dist(points[left_p], points[endpoint]) {
					// Really annoying edge case when jarvis won't converge b/c left is collinear with something "same distance"
					if dist(points[left_p], points[candidate]) == dist(points[left_p], points[endpoint]) && endpoint != left {
						// New point is collinear but further than current endpoint
						endpoint = candidate
					}
				}
			}
			left_p = endpoint
			// Circled back to original point
			if selected[left_p] {
				break
			}
		}
		wg.Done()
	}
	do_right := func(start int) {
		right_p := start
		for {
			selected[right_p] = true
			// Find leftmost endpoint
			endpoint := 0
			for candidate := range points {
				cross := cross_prod(points[right_p], points[endpoint], points[candidate])
				if endpoint == right_p || cross > 0 {
					// New point is to the left of current endpoint
					endpoint = candidate
				} else if cross == 0 && dist(points[right_p], points[candidate]) >= dist(points[right_p], points[endpoint]) {
					// Really annoying edge case when jarvis won't converge b/c left is collinear with something "same distance"
					if dist(points[right_p], points[candidate]) == dist(points[right_p], points[endpoint]) && endpoint != left {
						// New point is collinear but further than current endpoint
						endpoint = candidate
					}
				}
			}
			right_p = endpoint
			// Circled back to original point
			if selected[right_p] {
				break
			}
		}
		wg.Done()
	}
	wg.Add(8)
	// Can go CCW or CCW from max X, max Y, min X, min Y
	go do_left(left)
	go do_right(left)
	go do_left(right)
	go do_right(right)
	go do_left(up)
	go do_right(up)
	go do_left(down)
	go do_right(down)
	wg.Wait()

	for i := 0; i < len(points); i++ {
		if selected[i] {
			hull = append(hull, points[i])
		}
	}

	return hull
}

// Sequential Jarvis March
func seq_jarvis(points [][2]float32) [][2]float32 {
	if len(points) < 3 {
		return points
	}
	var hull [][2]float32
	left := leftmost(points)
	// Last selected point on hull
	p := left
	for {
		hull = append(hull, points[p])
		// Find leftmost endpoint
		endpoint := 0
		for candidate := range points {
			cross := cross_prod(points[p], points[endpoint], points[candidate])
			if endpoint == p || cross < 0 {
				// New point is to the left of current endpoint
				endpoint = candidate
			} else if cross == 0 && dist(points[p], points[candidate]) >= dist(points[p], points[endpoint]) {
				// Really annoying edge case when jarvis won't converge b/c left is collinear with something "same distance"
				if dist(points[p], points[candidate]) == dist(points[p], points[endpoint]) && endpoint != left {
					// New point is collinear but further than current endpoint
					endpoint = candidate
				}
			}
		}
		p = endpoint
		// Circled back to original point
		if endpoint == left {
			break
		}
	}

	return hull
}

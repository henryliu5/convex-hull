package main

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

// Parallel quicksort
func parallel_qsort(a [][]float32, cmp func([]float32, []float32) bool, wg *sync.WaitGroup) [][]float32 {
	if len(a) < 2 {
		wg.Done()
		return a
	}
	left, right := 0, len(a)-1
	pivot := 0

	a[pivot], a[right] = a[right], a[pivot]
	for i := 0; i < right; i++ {
		if cmp(a[i], a[right]) {
			a[left], a[i] = a[i], a[left]
			left++
		}
	}
	a[left], a[right] = a[right], a[left]
	if len(a) > 20000 {
		wg.Add(2)
		go parallel_qsort(a[:left], cmp, wg)
		go parallel_qsort(a[left+1:], cmp, wg)

	} else {
		qsort(a[:left], cmp)
		qsort(a[left+1:], cmp)
	}
	wg.Done()
	return a
}

// adapted from https://stackoverflow.com/a/55267961/15471686
func qsort(a [][]float32, cmp func([]float32, []float32) bool) [][]float32 {
	if len(a) < 2 {
		return a
	}
	left, right := 0, len(a)-1
	pivot := 0

	a[pivot], a[right] = a[right], a[pivot]
	for i := 0; i < right; i++ {
		if cmp(a[i], a[right]) {
			a[left], a[i] = a[i], a[left]
			left++
		}
	}
	a[left], a[right] = a[right], a[left]

	qsort(a[:left], cmp)
	qsort(a[left+1:], cmp)

	return a
}

// Sort by polar angle using custom quicksort (faster)
func custom_sort(a [][]float32, bot_point []float32, order float32) {
	fmt.Println("started sort")
	start := time.Now()
	// Sort by polar angle to bottom most point
	cmp := func(a, b []float32) bool {
		cross := cross_prod(bot_point, a, b)
		// Break colinear ties using distance
		if cross == 0 {
			// Point j should be further than point i
			return dist(bot_point, a) < dist(bot_point, b)
		} else {
			// Point j should be left of point i
			return order*cross > 0
		}
	}
	qsort(a, cmp)
	// wg := new(sync.WaitGroup)
	// wg.Add(1)
	// parallel_qsort(a, cmp, wg)
	// wg.Wait()

	fmt.Println("finished sort", time.Since(start))
}

// Sort by polar angle using Go builtin slice sort (slow)
func go_sort(a [][]float32, bot_point []float32, order float32) {
	fmt.Println("started sort")
	start := time.Now()
	// Sort by polar angle to bottom most point
	sort.Slice(a, func(i, j int) bool {
		cross := cross_prod(bot_point, a[i], a[j])
		// Break colinear ties using distance
		if cross == 0 {
			// Point j should be further than point i
			return dist(bot_point, a[i]) <= dist(bot_point, a[j])
		} else {
			// Point j should be left of point i
			return order*cross > 0
		}
	})
	fmt.Println("finished sort", time.Since(start))
}

// Sequential Graham Scan
func seq_graham_scan(points [][]float32) [][]float32 {
	fn_start := time.Now()
	// Set -1 for CW hull, 1 for CCW
	var order float32 = -1
	var float_error float32 = 0.00000001

	// Get index of bottommost & leftmost point
	var miny float32 = math.MaxFloat32
	var minx float32 = math.MaxFloat32
	bottom := -1
	for i, point := range points {
		if point[1] <= miny && point[0] < minx {
			bottom = i
			minx = point[0]
			miny = point[1]
		}
	}
	bot_point := points[bottom]
	points[0], points[bottom] = points[bottom], points[0]

	// fmt.Println("points pre-sort", points)

	sort_points := points[1:]
	custom_sort(sort_points, bot_point, order)
	// go_sort(sort_points, bot_point, order)

	// fmt.Println("points post-sort", points)

	// Remove collinear points
	new_index := 1
	for i := 1; i < len(points); i++ {
		for i < len(points)-1 && cross_prod(bot_point, points[i], points[i+1]) == 0 {
			i++
		}
		points[new_index] = points[i]
		new_index++
	}
	points = points[:new_index]

	// Iterate through points with Graham scan
	stack := make([]int, 0, len(points)/4)
	for i := 1; i < len(points); i++ {
		// Pop off stack if new point makes a clockwise turn
		for len(stack) > 1 && order*cross_prod(points[stack[len(stack)-2]], points[stack[len(stack)-1]], points[i]) <= float_error {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	// Stack now contains indices of convex hull
	var hull [][]float32 = [][]float32{bot_point}
	for _, i := range stack {
		hull = append(hull, points[i])
	}

	fmt.Println("graham scan", time.Since(fn_start))
	return hull
}

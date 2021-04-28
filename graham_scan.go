package main

import (
	"math"
	"sort"
	"sync"
	"time"
)

var PAR_QUICKSORT_LIMIT int = 2000

// Parallel quicksort
func parallel_qsort(a [][]float32, cmp func([]float32, []float32) bool, wg *sync.WaitGroup) {
	if len(a) >= 2 {
		left, right := 0, len(a)-1
		pivot := 0

		// Parititon
		a[pivot], a[right] = a[right], a[pivot]
		for i := 0; i < right; i++ {
			if cmp(a[i], a[right]) {
				a[left], a[i] = a[i], a[left]
				left++
			}
		}
		a[left], a[right] = a[right], a[left]

		// Spawn new goroutines if enough elements remaining
		if len(a) > PAR_QUICKSORT_LIMIT {
			wg.Add(2)
			go parallel_qsort(a[:left], cmp, wg)
			go parallel_qsort(a[left+1:], cmp, wg)

		} else {
			qsort(a[:left], cmp)
			qsort(a[left+1:], cmp)
		}
	}
	wg.Done()
}

// Quicksort adapted from https://stackoverflow.com/a/55267961/15471686
func qsort(a [][]float32, cmp func([]float32, []float32) bool) {
	if len(a) < 2 {
		return
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
}

// Iterative quicksort
func qsort2(a [][]float32, l, r int, cmp func([]float32, []float32) bool) {

	stack := make([]int, 0, 2*len(a))
	stack = append(stack, l)
	stack = append(stack, r)

	for len(stack) > 1 {
		back := len(stack) - 1
		r = stack[back]
		l = stack[back-1]

		stack = stack[:back-1]

		if l >= r || r-l+1 < 2 {
			continue
		}

		left, right := l, r
		pivot := left
		// fmt.Println(left, right)
		a[pivot], a[right] = a[right], a[pivot]
		for i := l; i < right; i++ {
			if cmp(a[i], a[right]) {
				a[left], a[i] = a[i], a[left]
				left++
			}
		}
		a[left], a[right] = a[right], a[left]

		stack = append(stack, l)
		stack = append(stack, left-1)
		stack = append(stack, left+1)
		stack = append(stack, r)
	}
}

// Sort by polar angle using custom quicksort (faster)
func custom_sort(a [][]float32, bot_point []float32, order float32) {
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
	// qsort(a, cmp)
	// qsort2(a, 0, len(a)-1, cmp)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	parallel_qsort(a, cmp, wg)
	wg.Wait()
}

// Sort by polar angle using Go builtin slice sort (slow)
func go_sort(a [][]float32, bot_point []float32, order float32) {
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
}

// Sequential Graham Scan
func seq_graham_scan_run(points [][]float32, clockwise bool) [][]float32 {
	// Set -1 for CW hull, 1 for CCW
	var order float32
	if clockwise {
		order = -1
	} else {
		order = 1
	}

	// Get index of bottommost & leftmost point
	var miny float32 = math.MaxFloat32
	var minx float32 = order * math.MaxFloat32
	bottom := -1
	for i, point := range points {
		if point[1] <= miny && order*point[0] < order*minx {
			bottom = i
			minx = point[0]
			miny = point[1]
		}
	}
	bot_point := points[bottom]
	points[0], points[bottom] = points[bottom], points[0]

	// Sort points based on angle to the bottom point
	sort_points := points[1:]
	debug("started sort")
	sort_start := time.Now()
	custom_sort(sort_points, bot_point, order)
	// go_sort(sort_points, bot_point, order)
	debug("finished sort", time.Since(sort_start))

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
	for i := 0; i < len(points); i++ {
		// Pop off stack if new point makes a clockwise turn
		for len(stack) > 1 && order*cross_prod(points[stack[len(stack)-2]], points[stack[len(stack)-1]], points[i]) <= 0 {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	// Stack now contains indices of convex hull
	var hull [][]float32 = [][]float32{}
	for _, i := range stack {
		hull = append(hull, points[i])
	}

	return hull
}

// Run graham scan, use two passes because of float associativity
func seq_graham_scan(points [][]float32) [][]float32 {
	if len(points) < 3 {
		return points
	}
	hull := seq_graham_scan_run(points, false)
	hull = seq_graham_scan_run(hull, true)
	return hull
}

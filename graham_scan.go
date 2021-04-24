package main

import (
	"math"
	"sort"
)

// Sequential Graham Scan
func seq_graham_scan(points [][]float32) [][]float32 {
	// Set -1 for CW hull, 1 for CCW
	var order float32 = -1

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

	points[0], points[bottom] = points[bottom], points[0]

	// fmt.Println("points pre-sort", points)
	sort_points := points[1:]
	// Sort by polar angle to bottom most point
	sort.SliceStable(sort_points, func(i, j int) bool {
		cross := cross_prod(points[bottom], sort_points[i], sort_points[j])
		// Break colinear ties using distance
		if cross == 0 {
			// Point j should be further than point i
			return dist(points[bottom], sort_points[i]) <= dist(points[bottom], sort_points[j])
		} else {
			// Point j should be left of point i
			return order*cross > 0
		}
	})
	// fmt.Println("points post-sort", points)
	// TODO remove colinear points

	stack := make([]int, 0, len(points)/4)

	for i, p := range points {
		// Pop off stack if new point makes a clockwise turn
		for len(stack) > 1 && order*cross_prod(points[stack[len(stack)-2]], points[stack[len(stack)-1]], p) <= 0 {
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}

	// Stack now contains indices of convex hull
	var hull [][]float32
	for _, i := range stack {
		hull = append(hull, points[i])
	}

	return hull
}

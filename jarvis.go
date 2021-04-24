package main

import (
	"math"
)

// Computes ab x ac
func ccw(a, b, c []float32) bool {
	x1 := a[0] - b[0]
	x2 := a[0] - c[0]
	y1 := a[1] - b[1]
	y2 := a[1] - c[1]
	cross := y2*x1 - y1*x2
	// c is counterclockwise of vector ab iff cross >0
	return cross > 0
}

// Get index of leftmost point
func leftmost(points [][]float32) int {
	var min float32 = math.MaxFloat32
	index := -1
	for i, point := range points {
		if point[0] < min {
			index = i
			min = point[0]
		}
	}
	return index
}

// Sequential Jarvis March
func seq_jarvis(points [][]float32) [][]float32 {
	// fmt.Println("points", points)
	// fmt.Println("leftmost", leftmost(points))
	var hull [][]float32

	left := leftmost(points)
	// Last selected point on hull
	p := left
	for {
		hull = append(hull, points[p])
		// Find leftmost endpoint
		endpoint := 0
		for candidate := range points {
			// TODO handle colinear points
			if endpoint == p || ccw(points[p], points[endpoint], points[candidate]) {
				endpoint = candidate
			}
		}
		p = endpoint
		if endpoint == left {
			break
		}
	}

	return hull
}

package main

import (
	"fmt"
	"math"
	"os"
)

// Computes ab x ac
func cross_prod(a, b, c []float32) float32 {
	x1 := a[0] - b[0]
	x2 := a[0] - c[0]
	y1 := a[1] - b[1]
	y2 := a[1] - c[1]

	// c is counterclockwise of vector ab iff > 0
	return y2*x1 - y1*x2
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

// Distance between two points
func dist(a, b []float32) float64 {
	x := a[0] - b[0]
	y := a[1] - b[1]
	return math.Sqrt(float64(x*x + y*y))
}

func debug(a ...interface{}) {
	// change this to switch on/off
	if false {
		fmt.Println(a...)
	}
}

// Output points to file
func output_points(filename string, points [][]float32) {
	f, _ := os.Create(filename)
	defer f.Close()

	for i := 0; i < len(points); i++ {
		fmt.Fprintf(f, "%f,%f\n", points[i][0], points[i][1])
	}
}

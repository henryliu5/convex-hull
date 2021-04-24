package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// actual hull
	a := [][]float32{
		{-1, -1},
		{-1, 1},
		{1, 1},
		{1, -1},
	}

	// random points inside
	for i := 0; i < 10; i++ {
		x := rand.Float32()*2 - 1
		y := rand.Float32()*2 - 1
		a = append(a, []float32{x, y})
	}

	// throw in some collinear points
	a = append(a, []float32{1, 0})

	fmt.Println(seq_jarvis(a))
	fmt.Println(seq_graham_scan(a))

}

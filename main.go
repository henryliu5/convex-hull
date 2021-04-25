package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func main() {
	n := 10000000
	start := time.Now()
	a := make([][]float32, 0, n+4)
	fmt.Println("allocation time", time.Since(start))

	// actual hull
	a = append(a, [][]float32{
		{-1, -1},
		{-1, 1},
		{1, 1},
		{1, -1},
	}...,
	)

	start = time.Now()
	// random points inside
	for i := 0; i < n; i++ {
		// x := rand.Float32()*2 - 1
		// y := rand.Float32()*2 - 1
		// a = append(a, []float32{x, y})

		// Integer test due to floating point error
		x := rand.Int() - math.MaxInt32/2
		y := rand.Int() - math.MaxInt32/2
		a = append(a, []float32{float32(x), float32(y)})
	}
	fmt.Println("creation time", time.Since(start))

	fn_start := time.Now()
	fmt.Println(len(seq_jarvis(a)))
	fmt.Println("jarvis march", time.Since(fn_start))

	a_1 := make([][]float32, len(a))
	copy(a_1, a)
	fn_start = time.Now()
	hull := seq_graham_scan(a_1)
	fmt.Println(len(hull))
	fmt.Println("graham scan", time.Since(fn_start))

	// output_points("serial_graham.txt", hull)
	// output_points("input.txt", a)

	fn_start = time.Now()
	fmt.Println(len(seq_chans(a)))
	fmt.Println("chan's", time.Since(fn_start))
}

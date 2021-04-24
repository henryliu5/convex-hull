package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	n := 1000000
	start := time.Now()
	a := make([][]float32, 0, n+4)
	fmt.Println("allocation time", time.Since(start))

	// // actual hull
	// a = append(a, [][]float32{
	// 	{-1, -1},
	// 	{-1, 1},
	// 	{1, 1},
	// 	{1, -1},
	// }...,
	// )

	start = time.Now()
	// random points inside
	for i := 0; i < n; i++ {
		// x := rand.Float32()*2 - 1
		// y := rand.Float32()*2 - 1
		// a = append(a, []float32{x, y})

		// Integer test due to floating point error
		r := 10000000
		x := rand.Intn(r) - r/2
		y := rand.Intn(r) - r/2
		a = append(a, []float32{float32(x), float32(y)})
	}
	fmt.Println("creation time", time.Since(start))

	fmt.Println(len(seq_jarvis(a)))
	fmt.Println(len(seq_graham_scan(a)))

}

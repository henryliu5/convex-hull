package main

import (
	"fmt"
	"time"
)

func main() {
	/**********************
	 * Points from memory *
	 **********************/
	// n := 10000000
	// start := time.Now()
	// a := make([][]float32, 0, n+4)
	// fmt.Println("allocation time", time.Since(start))

	// // actual hull
	// a = append(a, [][]float32{
	// 	{-1, -1},
	// 	{-1, 1},
	// 	{1, 1},
	// 	{1, -1},
	// }...,
	// )

	// start = time.Now()
	// // random points inside
	// for i := 0; i < n; i++ {
	// 	// x := rand.Float32()*2 - 1
	// 	// y := rand.Float32()*2 - 1
	// 	// a = append(a, []float32{x, y})

	// 	// Integer test due to floating point error
	// 	x := rand.Int() - math.MaxInt32/2
	// 	y := rand.Int() - math.MaxInt32/2
	// 	a = append(a, []float32{float32(x), float32(y)})
	// }
	// fmt.Println("creation time", time.Since(start))

	/********************
	 * Points from file *
	 ********************/

	points := parse_file("./serial_quickhull/input_points.txt")
	a := make([][]float32, 0, 10)
	for i := range points {
		a = append(a, points[i][0:2])
	}

	// Run jarvis march
	fn_start := time.Now()
	hull1 := seq_jarvis(a)
	fmt.Println("points on hull:", len(hull1))
	fmt.Println("jarvis march", time.Since(fn_start))

	output_points("serial_jarvis.txt", hull1)

	// Run graham scan
	a_1 := make([][]float32, len(a))
	copy(a_1, a)
	fn_start = time.Now()
	hull2 := seq_graham_scan(a_1)
	fmt.Println("points on hull:", len(hull2))
	fmt.Println("graham scan", time.Since(fn_start))

	output_points("serial_graham.txt", hull2)

	// Run chan's
	fn_start = time.Now()
	hull3 := seq_chans(a)
	fmt.Println("points on hull:", len(hull3))
	fmt.Println("chan's", time.Since(fn_start))
	output_points("serial_chan.txt", hull3)

	output_points("input.txt", a)
}

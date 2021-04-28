package main

import (
	"fmt"
	"time"
)

// Run convex hull using algorithm: method
func run_hull(points [][2]float32, method func([][2]float32) [][2]float32, name string) {
	fn_start := time.Now()
	hull := method(points)
	fmt.Println("points on hull:", len(hull))
	fmt.Println(name, time.Since(fn_start))
	// Write hull to output
	output_points(fmt.Sprintf("%s.txt", name), hull)
}

func main() {
	// /**********************
	//  * Points from memory *
	//  **********************/
	// n := 2000000
	// start := time.Now()
	// points := make([][2]float32, 0, n+4)
	// fmt.Println("allocation time", time.Since(start))

	// // actual hull
	// points = append(points, [][2]float32{
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
	// 	// a = append(a, [2]float32{x, y})

	// 	// Integer test due to floating point error
	// 	x := rand.Int31() - math.MaxInt32/2
	// 	y := rand.Int31() - math.MaxInt32/2
	// 	points = append(points, [2]float32{float32(x), float32(y)})
	// }
	// fmt.Println("creation time", time.Since(start))

	/********************
	 * Points from file *
	 ********************/
	points := parse_file("./serial_quickhull/input_points.txt")

	// Run jarvis march
	run_hull(points, seq_jarvis, "serial_jarvis")

	// Run graham scan
	a := make([][2]float32, len(points))
	copy(a, points)
	run_hull(a, seq_graham_scan, "serial_graham")

	// Run chan's
	run_hull(points, seq_chans, "serial_chans")

	// Run quickhull
	run_hull(points, quickhull, "serial_qh")

	output_points("input.txt", points)
}

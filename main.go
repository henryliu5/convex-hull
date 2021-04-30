package main

import (
	"fmt"
	"time"
	"os"
	"flag"
)

// Run convex hull using algorithm: method
func run_hull(points [][2]float32, method func([][2]float32) [][2]float32, name string, trials int, save_time bool, result_file string) {
	time_total := int64(0)
	for i := 0; i < trials; i++{
		fn_start := time.Now()
		hull := method(points)
		//fmt.Println("points on hull:", len(hull))
		//fmt.Println(name, time.Since(fn_start))
		ns_elap := time.Since(fn_start).Nanoseconds()
		time_total += (ns_elap)
		// Write hull to output
		output_points(fmt.Sprintf("%s.txt", name), hull)
	}
	avg_time := float64(time_total) / float64(trials)

	if (save_time){
		fmt.Println("Saving at ", result_file)
		f, _ := os.Create(result_file)
    	defer f.Close()
		fmt.Fprintf(f, "%f\n", avg_time)
	}
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

	result_dir_ptr := flag.String("result_dir", "", "result directory location")
	inputPtr := flag.String("input", "./serial_quickhull/input_points.txt", "input file location")
	num_trials_ptr := flag.Int("trials", 1, "number of trials")
	flag.Parse()

	save_time := (*result_dir_ptr != "")
	points := parse_file(*inputPtr)

	// Run jarvis march
	run_hull(points, seq_jarvis, "serial_jarvis",*num_trials_ptr, save_time, *result_dir_ptr + "/serial_jarvis.txt")

	// Run graham scan
	a := make([][2]float32, len(points))
	copy(a, points)
	run_hull(a, seq_graham_scan, "serial_graham",*num_trials_ptr, save_time, *result_dir_ptr + "/serial_graham.txt")

	// Run chan's
	run_hull(points, seq_chans, "serial_chans",*num_trials_ptr, save_time, *result_dir_ptr + "/serial_chan.txt")

	// Run quickhull
	run_hull(points, quickhull, "serial_qh",*num_trials_ptr, save_time, *result_dir_ptr + "/serial_qh.txt")

	output_points("input.txt", points)
}

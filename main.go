package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

// Run convex hull using algorithm: method
func run_hull(points [][2]float32, method func([][2]float32) [][2]float32, name string, trials int, save_time bool, result_file string, variable_of_interest string, do_output bool) {
	time_total := int64(0)
	points_copy := make([][2]float32, len(points))

	for i := 0; i < trials; i++ {
		// Copy because graham's will modify - next run will be O(N^2) quicksort otherwise
		copy(points_copy, points)

		fmt.Println()
		fn_start := time.Now()
		hull := method(points_copy)
		elapsed := time.Since(fn_start)
		fmt.Println(fmt.Sprintf("%s points on hull:", name), len(hull))
		fmt.Println(name, elapsed)

		ns_elap := time.Since(fn_start).Nanoseconds()
		time_total += (ns_elap)
		// Write hull to output
		if do_output {
			output_points(fmt.Sprintf("%s.txt", name), hull)
		}
	}
	avg_time := float64(time_total) / float64(trials)

	if save_time {
		fmt.Println("Saving at ", result_file)
		f, _ := os.OpenFile(result_file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		defer f.Close()
		fmt.Fprintf(f, "%s\t%d\t%f\t%s\n", name, trials, avg_time, variable_of_interest)
	}
}

func main() {
	result_file_ptr := flag.String("result_file", "", "result file location")
	inputPtr := flag.String("input", "./serial_quickhull/input_points.txt", "input file location")
	num_trials_ptr := flag.Int("trials", 1, "number of trials")
	//Pass something like number of points if you want it to be recorded in the data for later visualization
	variable_of_interest := flag.String("voi", "", "variable of interest to be recorded in data")
	do_output_ptr := flag.Bool("do_output", true, "output hull")
	go_maxprocs := flag.Int("procs", runtime.NumCPU(), "set runtime.GOMAXPROCS aka how many OS threads")
	do_coalesce := flag.Bool("coalesce", false, "enable coalescing of subhulls from chan's iterations")

	flag.Parse()

	// Set # OS threads
	runtime.GOMAXPROCS(*go_maxprocs)
	// Enable coalescing of subhulls thru iterations
	USE_COALESCE = *do_coalesce

	do_output := *do_output_ptr

	save_time := (*result_file_ptr != "")
	points := parse_file(*inputPtr)

	if len(points) == 0 {
		fmt.Printf("file: %s not found!\n", *inputPtr)
		return
	}

	// Run jarvis march
	run_hull(points, seq_jarvis, "serial_jarvis", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
	// run_hull(points, naive_parallel_jarvis, "naive_parallel_jarvis", *num_trials_ptr, save_time, *result_dir_ptr+"/naive_parallel_jarvis.txt")
	run_hull(points, parallel_jarvis, "parallel_jarvis", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)

	// Run graham scan
	run_hull(points, seq_graham_scan, "serial_graham", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
	run_hull(points, parallel_graham_scan, "parallel_graham", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)

	// Run chan's
	run_hull(points, seq_chans, "serial_chans", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
	run_hull(points, parallel_chans, "parallel_chans", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)

	// Run quickhull
	run_hull(points, quickhull_serial, "serial_qh", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
	run_hull(points, quickhull_parallel, "parallel_qh", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)

	// output_points("input.txt", points)
}

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

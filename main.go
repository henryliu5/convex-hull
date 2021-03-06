package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
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

	// Pass something like number of points if you want it to be recorded in the data for later visualization
	variable_of_interest := flag.String("voi", "", "variable of interest to be recorded in data")
	do_output_ptr := flag.Bool("do_output", true, "output hull")
	go_maxprocs := flag.Int("procs", runtime.NumCPU(), "set runtime.GOMAXPROCS aka how many OS threads")
	do_coalesce := flag.Bool("coalesce", false, "enable coalescing of subhulls from chan's iterations")
	simul_iters := flag.Int("simul_iters", 2, "how many iterations of chan's to run simultaneously")
	impl_ptr := flag.String("impl", "", "which implementation to run")

	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile := flag.String("memprofile", "", "write memory profile to file")

	flag.Parse()

	// CPU/mem profiling w/ pproc - https://blog.golang.org/pprof
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Set # OS threads
	runtime.GOMAXPROCS(*go_maxprocs)
	// Enable coalescing of subhulls thru iterations
	USE_COALESCE = *do_coalesce
	// Set number of iterations of chan's to run simultaneously
	SIMUL_ITERS = *simul_iters

	do_output := *do_output_ptr
	impl := *impl_ptr

	save_time := (*result_file_ptr != "")
	points := parse_file(*inputPtr)

	if len(points) == 0 {
		fmt.Printf("file: %s not found!\n", *inputPtr)
		return
	}

	// Run jarvis march
	if impl == "" || impl == "jarv" {
		run_hull(points, seq_jarvis, "serial_jarvis", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
		run_hull(points, parallel_jarvis, "parallel_jarvis", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
	}

	// Run graham scan
	if impl == "" || impl == "grah" {
		run_hull(points, seq_graham_scan, "serial_graham", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
		run_hull(points, parallel_graham_scan, "parallel_graham", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
	}

	// Run chan's
	if impl == "" || impl == "chan" {
		run_hull(points, seq_chans, "serial_chans", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
		run_hull(points, parallel_chans, "parallel_chans", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
	}

	// Run quickhull
	if impl == "" || impl == "quic" {
		run_hull(points, quickhull_serial, "serial_qh", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
		run_hull(points, quickhull_parallel, "parallel_qh", *num_trials_ptr, save_time, *result_file_ptr, *variable_of_interest, do_output)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}
}

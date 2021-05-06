# Parallel Convex Hull Algorithms
Implemenation of Chan's Algorithm (Jarvis March + Graham Scan) and Quickhull in Golang.

Run with
```
make
./runner [flags]

Usage of ./runner:
  -coalesce
        enable coalescing of subhulls from chan's iterations
  -cpuprofile string
        write cpu profile to file
  -do_output
        output hull (default true)
  -impl string
        which implementation to run (jarv/grah/chan/quic for Jarvis March, Graham Scan,
        Chan's Algorithm, Quickhull, respectively)
  -input string
        input file location (default "./serial_quickhull/input_points.txt")
  -procs int
        set runtime.GOMAXPROCS aka how many OS threads (default 20)
  -result_file string
        result file location
  -simul_iters int
        how many iterations of chan's to run simultaneously (default 2)
  -trials int
        number of trials (default 1)
  -voi string
        variable of interest to be recorded in data
```

### File Structure
* main.go - Main driver that handles file I/O and end-to-end timing
* chan.go - Sequential implementation of Chan's algorithm, including modified Jarvis march
    * Uses functions in graham_scan.go to compute subhulls
* concurrent_map.go - Implementation of a custom concurrent hash map for [2]float32 (used by parallel_chan.go)
* graham_scan.go - Parallel and sequential implementations of the Graham scan algorithm
* jarvis.go - Parallel and sequential implementations of the Jarvis march algorithm
* parallel_chan.go - Parallel implementation of Chan's algorithm, includes management of parallel subhull computation communication
    * Uses parallel graham scan implemented in graham_scan.go
* quickhull.go - Parallel and sequential implementations of Quickhull
* util.go - Functions for file I/O and common geometric calculations like cross product
* test_case_generation - Scripts to generate test cases - see gen_tests.sh for usage example
* verification - Visualization tool to see convex hull points

Created by Henry Liu and Raymond Hong
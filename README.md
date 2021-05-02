# Parallel Convex Hull Algorithms
Implemenation of Chan's Algorithm (Jarvis March + Graham Scan) and Quickhull in Golang.

Run with
```bash
make
./runner [flags]

Usage of ./runner:
  -input string
        input file location (default "./serial_quickhull/input_points.txt")
  -result_dir string
        result directory location
  -trials int
        number of trials (default 1)
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

Created by Henry Liu and Raymond Hong
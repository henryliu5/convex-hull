package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

// Computes ab x ac
func cross_prod(a, b, c [2]float32) float32 {
	x1 := a[0] - b[0]
	x2 := a[0] - c[0]
	y1 := a[1] - b[1]
	y2 := a[1] - c[1]

	// c is counterclockwise of vector ab iff > 0
	return y2*x1 - y1*x2
}

// Get index of leftmost point
func leftmost(points [][2]float32) int {
	var min float32 = math.MaxFloat32
	index := -1
	for i, point := range points {
		if point[0] < min {
			index = i
			min = point[0]
		}
	}
	return index
}

// Distance between two points
func dist(a, b [2]float32) float64 {
	x := a[0] - b[0]
	y := a[1] - b[1]
	return math.Sqrt(float64(x*x + y*y))
}

func debug(a ...interface{}) {
	// change this to switch on/off
	if false {
		fmt.Println(a...)
	}
}

// Output points to file
func output_points(filename string, points [][2]float32) {
	f, _ := os.Create(filename)
	defer f.Close()

	for i := 0; i < len(points); i++ {
		fmt.Fprintf(f, "%f,%f\n", points[i][0], points[i][1])
	}
}

// Parse file contents into memory
func parse_file(filename string) [][2]float32 {
	//Counts line in a file
	count_lines := func(file_str string) int {
		file, _ := os.Open(file_str)
		scanner := bufio.NewScanner(file)
		count := 0

		for scanner.Scan() {
			_ = scanner.Text()
			count = count + 1
		}
		return count
	}

	file, _ := os.Open(filename)

	num_lines := count_lines(filename)
	scanner := bufio.NewScanner(file)

	points := make([][2]float32, num_lines, num_lines)

	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		lst := strings.Split(line, ",")

		x_64, _ := strconv.ParseFloat(lst[0], 32)
		y_64, _ := strconv.ParseFloat(lst[1], 32)

		x := float32(x_64)
		y := float32(y_64)

		points[i][0] = x
		points[i][1] = y

		i = i + 1
	}

	return points
}

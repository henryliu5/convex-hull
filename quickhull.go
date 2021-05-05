package main

import (
	"bufio"
	"math"
	"os"
	"sync"
)

var convex_hull map[[2]float32]bool
var hull_lock sync.Mutex

//Counts line in a file
func count_lines(file_str string) int {
	file, _ := os.Open(file_str)
	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		_ = scanner.Text()
		count = count + 1
	}
	return count
}

func getMaxMinPt(points [][2]float32) [2]int {

	max_pt_x := float32(-math.MaxFloat32)
	max_pt_y := float32(-math.MaxFloat32)
	max_pt_ind := -1

	min_pt_x := float32(math.MaxFloat32)
	min_pt_y := float32(math.MaxFloat32)
	min_pt_ind := -1

	for i := 0; i < len(points); i++ {
		pt := points[i]
		if pt[0] > max_pt_x || (pt[0] == max_pt_x && pt[1] > max_pt_y) {
			max_pt_x = pt[0]
			max_pt_y = pt[1]
			max_pt_ind = i
		}

		if pt[0] < min_pt_x || (pt[0] == min_pt_x && pt[1] < min_pt_y) {
			min_pt_x = pt[0]
			min_pt_y = pt[1]
			min_pt_ind = i
		}
	}

	var res [2]int
	res[0] = max_pt_ind
	res[1] = min_pt_ind

	return res
}

func is_above(l1 [2]float32, l2 [2]float32, p [2]float32) float32 {
	AB := []float32{l2[0] - l1[0], l2[1] - l1[1]}
	AX := []float32{p[0] - l1[0], p[1] - l1[1]}
	cross := AB[0]*AX[1] - AB[1]*AX[0]
	return cross
}

func point_line_dist(l1 [2]float32, l2 [2]float32, pt [2]float32) float32 {
	x0 := pt[0]
	x1 := l1[0]
	x2 := l2[0]
	y0 := pt[1]
	y1 := l1[1]
	y2 := l2[1]

	num := math.Abs(float64((x2-x1)*(y1-y0) - (x1-x0)*(y2-y1)))
	den := math.Sqrt(float64(x2-x1)*float64(x2-x1) + float64(y2-y1)*float64(y2-y1))
	return float32(num / den)
}

func getSide(l1 [2]float32, l2 [2]float32, p [2]float32) int {
	AB := []float32{l2[0] - l1[0], l2[1] - l1[1]}
	AX := []float32{p[0] - l1[0], p[1] - l1[1]}
	cross := AB[0]*AX[1] - AB[1]*AX[0]
	if cross > 0 {
		return 1
	} else if cross < 0 {
		return -1
	}
	return 0
}

func hull(points [][2]float32, min_pt [2]float32, max_pt [2]float32, side int) {
	max_dist := float32(0.0)
	ind := -1

	new_points := make([][2]float32, 0)

	for i := 0; i < len(points); i++ {
		pt := points[i]
		dist := point_line_dist(min_pt, max_pt, pt)

		correct_side := getSide(min_pt, max_pt, pt) == side
		if correct_side && dist > max_dist {
			ind = i
			max_dist = dist
		}

		if correct_side {
			new_points = append(new_points, pt)
		}
	}

	if ind == -1 {
		//Add max, min
		convex_hull[min_pt] = true
		convex_hull[max_pt] = true
	} else {
		hull(new_points, points[ind], min_pt, -getSide(points[ind], min_pt, max_pt))
		hull(new_points, points[ind], max_pt, -getSide(points[ind], max_pt, min_pt))
	}
}

func quickhull(points [][2]float32) {
	res := getMaxMinPt(points)

	var min_pt [2]float32
	var max_pt [2]float32

	min_pt[0] = points[res[0]][0]
	min_pt[1] = points[res[0]][1]

	max_pt[0] = points[res[1]][0]
	max_pt[1] = points[res[1]][1]

	hull(points, max_pt, min_pt, 1)
	hull(points, max_pt, min_pt, -1)
}

func hull_p(points [][2]float32, min_pt [2]float32, max_pt [2]float32, side int, c chan int) {
	max_dist := float32(0.0)
	ind := -1

	new_points := make([][2]float32, 0)

	for i := 0; i < len(points); i++ {
		pt := points[i]
		dist := point_line_dist(min_pt, max_pt, pt)

		correct_side := getSide(min_pt, max_pt, pt) == side
		if correct_side && dist > max_dist {
			ind = i
			max_dist = dist
		}

		if correct_side {
			new_points = append(new_points, pt)
		}
	}

	if ind == -1 {
		//Add max, min
		hull_lock.Lock()
		convex_hull[min_pt] = true
		convex_hull[max_pt] = true
		hull_lock.Unlock()
		c <- 1
	} else {

		leftChan := make(chan int, 1)
		rightChan := make(chan int, 1)

		go hull_p(new_points, points[ind], min_pt, -getSide(points[ind], min_pt, max_pt), leftChan)
		go hull_p(new_points, points[ind], max_pt, -getSide(points[ind], max_pt, min_pt), rightChan)

		_ = <-leftChan
		_ = <-rightChan
		c <- 1
	}
}

func quickhull_p(points [][2]float32) {
	res := getMaxMinPt(points)

	var min_pt [2]float32
	var max_pt [2]float32

	min_pt[0] = points[res[0]][0]
	min_pt[1] = points[res[0]][1]

	max_pt[0] = points[res[1]][0]
	max_pt[1] = points[res[1]][1]

	leftChan := make(chan int, 1)
	rightChan := make(chan int, 1)

	go hull_p(points, max_pt, min_pt, 1, leftChan)
	go hull_p(points, max_pt, min_pt, -1, rightChan)

	_ = <-leftChan
	_ = <-rightChan
}

func quickhull_serial(points [][2]float32) [][2]float32 {

	convex_hull = make(map[[2]float32]bool)

	quickhull(points)

	hull_res := make([][2]float32, 0)

	for key, _ := range convex_hull {
		hull_res = append(hull_res, key)
	}
	return hull_res
}

func quickhull_parallel(points [][2]float32) [][2]float32 {
	convex_hull = make(map[[2]float32]bool)

	quickhull_p(points)

	hull_res := make([][2]float32, 0)

	for key, _ := range convex_hull {
		hull_res = append(hull_res, key)
	}
	return hull_res
}

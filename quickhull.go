package main

import (
	"math"
)

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
	den := math.Sqrt((math.Pow(float64(x2-x1), 2) + math.Pow(float64(y2-y1), 2)))
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

func hull(points [][2]float32, min_pt_index int, max_pt_index int, side int, convex_hull map[int]bool) {
	max_dist := float32(0.0)
	ind := -1

	max_pt := points[max_pt_index]
	min_pt := points[min_pt_index]

	for i := 0; i < len(points); i++ {
		pt := points[i]
		dist := point_line_dist(min_pt, max_pt, pt)

		if getSide(min_pt, max_pt, pt) == side && dist > max_dist {
			ind = i
			max_dist = dist
		}
	}

	if ind == -1 {
		//Add max, min
		// fmt.Println(max_pt)
		// fmt.Println(min_pt)
		convex_hull[min_pt_index] = true
		convex_hull[max_pt_index] = true
	} else {
		hull(points, ind, min_pt_index, -getSide(points[ind], min_pt, max_pt), convex_hull)
		hull(points, ind, max_pt_index, -getSide(points[ind], max_pt, min_pt), convex_hull)
	}
}

func quickhull(points [][2]float32) [][2]float32 {
	res := getMaxMinPt(points)

	var min_pt [2]float32
	var max_pt [2]float32
	convex_hull := make(map[int]bool)

	min_pt[0] = points[res[0]][0]
	min_pt[1] = points[res[0]][1]

	max_pt[0] = points[res[1]][0]
	max_pt[1] = points[res[1]][1]

	max_pt_index := res[0]
	min_pt_index := res[1]

	hull(points, max_pt_index, min_pt_index, 1, convex_hull)
	hull(points, max_pt_index, min_pt_index, -1, convex_hull)

	result := make([][2]float32, 0, len(points)>>1)
	for i := 0; i < len(points); i++ {
		if convex_hull[i] {
			result = append(result, points[i])
		}
	}
	return result
}

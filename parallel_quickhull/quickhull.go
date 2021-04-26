package main
import (
    "flag"
	"os"
	"bufio"
    "fmt"
	"strings"
	"strconv"
	"math"
)

var convex_hull [][2]float32

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

func getMaxMinPt(points [][2]float32) [2][2]float32 {

	max_pt_x := float32(-math.MaxFloat32)
	max_pt_y := float32(-math.MaxFloat32)

	min_pt_x := float32(math.MaxFloat32)
	min_pt_y := float32(math.MaxFloat32)

	for i := 0; i < len(points); i++{
		pt := points[i]
		if (pt[0] > max_pt_x || (pt[0] == max_pt_x && pt[1] > max_pt_y)){
			max_pt_x = pt[0];
			max_pt_y = pt[1];
		}

		if (pt[0] < min_pt_x || (pt[0] == min_pt_x && pt[1] < min_pt_y)){
			min_pt_x = pt[0];
			min_pt_y = pt[1];
		}
	}

	var res [2][2]float32
	res[0][0] = max_pt_x
	res[0][1] = max_pt_y

	res[1][0] = min_pt_x
	res[1][1] = min_pt_y

	return res
}

func is_above(l1 [2]float32, l2 [2]float32, p [2]float32) float32{
	AB := []float32{l2[0]-l1[0], l2[1]-l1[1]}
	AX := []float32{p[0]-l1[0], p[1]-l1[1]}
	cross := AB[0] * AX[1] - AB[1] * AX[0]
	return cross
}

func point_line_dist(l1 [2]float32, l2 [2]float32, pt [2]float32) float32{
	x0 := pt[0]
	x1 := l1[0]
	x2 := l2[0]
	y0 := pt[1]
	y1 := l1[1]
	y2 := l2[1]

	num := math.Abs(float64((x2-x1) * (y1-y0) - (x1-x0) * (y2-y1)))
	den := math.Sqrt((math.Pow(float64(x2-x1),2) + math.Pow(float64(y2-y1),2)))
    return float32(num/den)
}

func side_of_line_point_is(l1 [2]float32, l2 [2]float32, p [2]float32) float32{
	AB := []float32{l2[0]-l1[0], l2[1]-l1[1]}
	AX := []float32{p[0]-l1[0], p[1]-l1[1]}
	cross := AB[0] * AX[1] - AB[1] * AX[0]
	return cross
}

func isInsideTriangle(A [2]float32, B [2]float32, C [2]float32, pt [2]float32) bool{
	if ((side_of_line_point_is(A,B,pt) < 0) && (side_of_line_point_is(B,C,pt) < 0) && (side_of_line_point_is(C,A,pt) < 0)){
		return true
	}
	return false
}

func hull(points[][2]float32, min_pt [2]float32, max_pt [2]float32, c chan int){
	fmt.Println(points)
	if (len(points) == 0){
		c <- 1
		return
	}
	max_dist := float32(-1.0)
	var furthest_pt [2]float32;
	furthest_index := -1

	for i := 0; i < len(points); i++{
		pt := points[i]
		dist := point_line_dist(min_pt, max_pt, pt)
		if (dist > max_dist){
			max_dist = dist
			furthest_pt[0] = pt[0]
			furthest_pt[1] = pt[1]
			furthest_index = i
		}
	}

	points = append(points[:furthest_index], points[furthest_index + 1:]...)

	convex_hull = append(convex_hull, furthest_pt)

	L1 := make([][2]float32, 0, 0)
	L2 := make([][2]float32, 0, 0)

	for i := 0; i < len(points); i++{
		var pt [2]float32
		pt[0] = points[i][0]
		pt[1] = points[i][1]

		if (isInsideTriangle(min_pt, furthest_pt, max_pt, pt)){
			//Cannot be part of hull
		} else if (side_of_line_point_is(min_pt, furthest_pt, pt) > 0){
			L1 = append(L1, pt)
		} else if (side_of_line_point_is(min_pt, furthest_pt, pt) < 0){
			L2 = append(L2, pt)
		}
	}

	leftChan := make(chan int, 1)
	rightChan := make(chan int, 1)

	go hull(L1, furthest_pt,max_pt,leftChan)
    go hull(L2, min_pt, furthest_pt,rightChan)

	_ = <-leftChan
	_ = <-rightChan

	c <- 1
}

func quickhull(points [][2]float32){
	res := getMaxMinPt(points)

	var min_pt [2]float32
	var max_pt [2]float32

	min_pt[0] = res[1][0]
	min_pt[1] = res[1][1]

	max_pt[0] = res[0][0]
	max_pt[1] = res[0][1]

	convex_hull = append(convex_hull, min_pt)
	convex_hull = append(convex_hull, max_pt)

	above_pts := make([][2]float32, 0, 0)
	below_pts := make([][2]float32, 0, 0)

	for i := 0; i < len(points); i++{
		var pt [2]float32
		pt[0] = points[i][0]
		pt[1] = points[i][1]

		is_abv := is_above(min_pt, max_pt, pt)
		if is_abv > 0 {
			above_pts = append(above_pts, pt)
		}else if is_abv < 0 {
			below_pts = append(below_pts, pt)
		}
	}

	leftChan := make(chan int, 1)
	rightChan := make(chan int, 1)

	go hull(above_pts, min_pt, max_pt,leftChan)
	go hull(below_pts, max_pt, min_pt,rightChan)

	_ = <-leftChan
	_ = <-rightChan
}

func main() {
	inputPtr := flag.String("input", "", "input file location")
	flag.Parse()

	file, _ := os.Open(*inputPtr)

	num_lines := count_lines(*inputPtr)
	scanner := bufio.NewScanner(file)

	points := make([][2]float32, num_lines, num_lines)
	
	i:=0
	for scanner.Scan() {
        line := scanner.Text()
		lst := strings.Split(line, ",")

		x_64, _ := strconv.ParseFloat(lst[0],32)
		y_64, _ :=strconv.ParseFloat(lst[1],32)

		x := float32(x_64)
		y := float32(y_64)

		points[i][0] = x
		points[i][1] = y

		i=i+1
	}

	quickhull(points)

	f, _ := os.Create("parallel_qh_hull.txt")
    defer f.Close()

	for i := 0; i < len(convex_hull); i++{
		fmt.Fprintf(f, "%f,%f\n", convex_hull[i][0], convex_hull[i][1])
	}
}
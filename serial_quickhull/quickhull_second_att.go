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

var convex_hull map[int]bool

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

	for i := 0; i < len(points); i++{
		pt := points[i]
		if (pt[0] > max_pt_x || (pt[0] == max_pt_x && pt[1] > max_pt_y)){
			max_pt_x = pt[0];
			max_pt_y = pt[1];
			max_pt_ind = i
		}

		if (pt[0] < min_pt_x || (pt[0] == min_pt_x && pt[1] < min_pt_y)){
			min_pt_x = pt[0];
			min_pt_y = pt[1];
			min_pt_ind = i
		}
	}

	var res [2]int
	res[0] = max_pt_ind
	res[1] = min_pt_ind

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

func getSide(l1 [2]float32, l2 [2]float32, p [2]float32) int {
	AB := []float32{l2[0]-l1[0], l2[1]-l1[1]}
	AX := []float32{p[0]-l1[0], p[1]-l1[1]}
	cross := AB[0] * AX[1] - AB[1] * AX[0]
	if (cross > 0){
		return 1
	} else if (cross < 0){
		return -1
	}
	return 0
}


func hull(points[][2]float32, min_pt_index int, max_pt_index int, side int){
	max_dist := float32(0.0)
	ind := -1

	max_pt := points[max_pt_index]
	min_pt := points[min_pt_index]

	for i := 0; i < len(points); i++{
		pt := points[i]
		dist := point_line_dist(min_pt, max_pt, pt)

		if (getSide(min_pt, max_pt, pt) == side && dist > max_dist){
			ind = i
			max_dist = dist
		}
	}

	if (ind == -1){
		//Add max, min
		fmt.Println(max_pt)
		fmt.Println(min_pt)
		convex_hull[min_pt_index]=true
		convex_hull[max_pt_index]=true
	} else{
		
	}
}

func quickhull(points [][2]float32){
	res := getMaxMinPt(points)

	var min_pt [2]float32
	var max_pt [2]float32

	min_pt[0] = points[res[0]][0]
	min_pt[1] = points[res[0]][1]

	max_pt[0] = points[res[1]][0]
	max_pt[1] = points[res[1]][1]

	max_pt_index := res[0]
	min_pt_index := res[1]

	hull(points, max_pt_index, min_pt_index, 1)
	hull(points, max_pt_index, min_pt_index, -1)
}

func main() {
	inputPtr := flag.String("input", "", "input file location")
	flag.Parse()

	file, _ := os.Open(*inputPtr)

	num_lines := count_lines(*inputPtr)
	scanner := bufio.NewScanner(file)

	points := make([][2]float32, num_lines, num_lines)
	
	i:=0
	convex_hull = make(map[int]bool)
	for scanner.Scan() {
        line := scanner.Text()
		lst := strings.Split(line, ",")

		x_64, _ := strconv.ParseFloat(lst[0],32)
		y_64, _ :=strconv.ParseFloat(lst[1],32)

		x := float32(x_64)
		y := float32(y_64)

		points[i][0] = x
		points[i][1] = y

		convex_hull[i] = false
		i=i+1
	}

	quickhull(points)

	f, _ := os.Create("serial_qh_hull.txt")
    defer f.Close()
	
	for i := 0; i < len(convex_hull); i++{
		if (convex_hull[i]){
			fmt.Fprintf(f, "%f,%f\n", points[i][0], points[i][1])
		}
	}
}
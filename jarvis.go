package main

// Sequential Jarvis March
func seq_jarvis(points [][]float32) [][]float32 {
	// fmt.Println("points", points)
	// fmt.Println("leftmost", leftmost(points))
	var hull [][]float32

	left := leftmost(points)
	// Last selected point on hull
	p := left
	for {
		hull = append(hull, points[p])
		// Find leftmost endpoint
		endpoint := 0
		for candidate := range points {
			// TODO handle colinear points
			if endpoint == p || cross_prod(points[p], points[endpoint], points[candidate]) > 0 {
				endpoint = candidate
			}
		}
		p = endpoint
		if endpoint == left {
			break
		}
	}

	return hull
}

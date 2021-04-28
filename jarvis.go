package main

// Sequential Jarvis March
func seq_jarvis(points [][]float32) [][]float32 {
	if len(points) < 3 {
		return points
	}
	var hull [][]float32
	left := leftmost(points)
	// Last selected point on hull
	p := left
	for {
		hull = append(hull, points[p])
		// Find leftmost endpoint
		endpoint := 0
		for candidate := range points {
			cross := cross_prod(points[p], points[endpoint], points[candidate])
			if endpoint == p || cross > 0 {
				// New point is to the left of current endpoint
				endpoint = candidate
			} else if cross == 0 && dist(points[p], points[candidate]) >= dist(points[p], points[endpoint]) {
				// Really annoying edge case when jarvis won't converge b/c left is collinear with something "same distance"
				if dist(points[p], points[candidate]) == dist(points[p], points[endpoint]) && endpoint != left {
					// New point is collinear but further than current endpoint
					endpoint = candidate
				}
			}
		}
		p = endpoint
		// Circled back to original point
		if endpoint == left {
			break
		}
	}

	return hull
}

package knn

import (
	"math"
)

// ManhattanDistance measure Manhattan Distance between p1 and p2
// |x2 - x1| + |y2 - y1|
func ManhattanDistance(point1, point2 []int) float64 {
	if len(point1) != len(point2) {
		panic("len point1 and point2 should be equal")
	}

	var dist float64
	for i, p2 := range point2 {
		disI := p2 - point1[i]
		dist += math.Abs(float64(disI))
	}

	return dist
}

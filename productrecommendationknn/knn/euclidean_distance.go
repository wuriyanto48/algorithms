package knn

import (
	"math"
)

// EuclideanDistance measure Euclidean Distance between point1 and point1
func EuclideanDistance(point1, point2 []int) float64 {
	if len(point1) != len(point2) {
		panic("len point1 and point2 should be equal")
	}

	var dbs float64
	for i, p2 := range point2 {
		disI := p2 - point1[i]
		dbs += math.Pow(float64(disI), 2)
	}

	dis := math.Sqrt(dbs)
	return dis
}

package main

import (
	"bytes"
	"fmt"
)

func main() {
	points := Points{
		{X: 5, Y: 99},
		{X: 7, Y: 86},
		{X: 8, Y: 87},
		{X: 7, Y: 88},
		{X: 2, Y: 111},
		{X: 17, Y: 86},
		{X: 2, Y: 103},
		{X: 9, Y: 87},
		{X: 4, Y: 94},
		{X: 11, Y: 78},
		{X: 12, Y: 77},
		{X: 9, Y: 85},
		{X: 6, Y: 86},
	}

	var plane bytes.Buffer

	canvasDot := "* "
	width := 120
	height := 120

	cartesian := NewCartesianPlane(canvasDot, width, height)

	// draw plane
	cartesian.Draw()

	for i := len(points) - 1; i >= 0; i-- {
		ds := points[i]

		cartesian.Plot("0", &ds)
	}

	cartesian.WriteTo(&plane)
	fmt.Print(plane.String())

	rSquared := points.RSquared()
	fmt.Println("R squared: ", rSquared)
	fmt.Println("std err: ", points.StdError())
	fmt.Println("slope: ", points.Slope())
	fmt.Println("y intercept: ", points.GetYIntercept())

	fmt.Println(points.Predict(6))
}

// package main

// import (
// 	"fmt"
// )

// func main() {
// 	points := Points{
// 		{X: 5, Y: 99},
// 		{X: 7, Y: 86},
// 		{X: 8, Y: 87},
// 		{X: 7, Y: 88},
// 		{X: 2, Y: 111},
// 		{X: 17, Y: 86},
// 		{X: 2, Y: 103},
// 		{X: 9, Y: 87},
// 		{X: 4, Y: 94},
// 		{X: 11, Y: 78},
// 		{X: 12, Y: 77},
// 		{X: 9, Y: 85},
// 		{X: 6, Y: 86},
// 	}

// 	rSquared := points.RSquared()
// 	fmt.Println(rSquared)
// 	fmt.Println(points.Slope())
// 	fmt.Println(points.GetYIntercept())

// 	fmt.Println(points.Predict(10))
// }

package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	fmt.Println(CosineDistance(Point{X: 8, Y: 21}, Point{X: 8, Y: 22}))
	fmt.Println(EuclideanDistance(Point{X: 0, Y: 0}, Point{X: 3, Y: 2}))
	fmt.Println(ManhattanDistance(Point{X: 1, Y: 1}, Point{X: 4, Y: 3}))
	// k
	k := 5

	var plane bytes.Buffer

	canvasDot := "* "
	width := 40
	height := 40

	cartesian := NewCartesianPlane(canvasDot, width, height)

	// draw plane
	cartesian.Draw()

	dataSets := Points{
		{X: 4, Y: 21},
		{X: 5, Y: 19},
		{X: 10, Y: 24},
		{X: 4, Y: 17},
		{X: 3, Y: 16},
		{X: 11, Y: 25},
		{X: 14, Y: 24},
		{X: 8, Y: 22},
		{X: 10, Y: 21},
		{X: 12, Y: 21},
	}

	classes := []int{
		0, 0, 1, 0, 0, 1, 1, 0, 1, 1,
	}

	// construct labelled
	// var labelledDatasets []LabelledDataset

	dataTest := Point{X: 8, Y: 21}

	// plot the data test
	cartesian.Plot("@ ", &dataTest)

	// KNN construct
	knn := NewKNN(k)

	// fit data
	if err := knn.FitData(dataSets, classes); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// plotting
	knn.Plot(cartesian)

	fmt.Println("---------------------------------")

	cartesian.WriteTo(&plane)
	fmt.Print(plane.String())

	classifiedDataSets := knn.Predict(&dataTest)

	for i := 0; i < len(classifiedDataSets); i++ {
		dataSet := classifiedDataSets[i]
		fmt.Printf("(x1: %d - x2:%d) | (y1: %d - y2: %d) | dist: %f | label: %d\n",
			dataTest.X, dataSet.Point.X,
			dataTest.Y, dataSet.Point.Y,
			dataSet.Distance,
			dataSet.Class)
	}

}

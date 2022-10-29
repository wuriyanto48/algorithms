package main

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strings"
)

// Point represent cartesian Point
type Point struct {
	X int
	Y int
}

func (p Point) Draw(s string) string {
	return fmt.Sprintf("%s", s)
}

// Points represent collections of point
type Points []Point

// LabelledDataset represent labelled dataset
type LabelledDataset struct {
	Class    int
	Point    Point
	Distance float64
}

// CosineSimilarity will measure Cosine Similarity between p1 and p2
func CosineSimilarity(p1, p2 Point) float64 {
	// product of p1 and p2 (a.b)
	p1Dotp2 := (p1.X * p2.X) + (p1.Y * p2.Y)

	// p1 and p2 magnitude |a||b|
	p1Magnitude := math.Sqrt(math.Pow(float64(p1.X), 2) + math.Pow(float64(p1.Y), 2))
	p2Magnitude := math.Sqrt(math.Pow(float64(p2.X), 2) + math.Pow(float64(p2.Y), 2))
	// fmt.Println(p1Dotp2, " ", p1Magnitude, " ", p2Magnitude)

	res := float64(p1Dotp2) / (p1Magnitude * p2Magnitude)
	// fmt.Println(math.Acos(res)* (180 / math.Pi))
	return res
}

// CosineDistance will measure Cosine Similarity between p1 and p2
func CosineDistance(p1, p2 Point) float64 {
	cosSim := CosineSimilarity(p1, p2)
	return 1 - cosSim
}

// EuclideanDistance will measure Euclidean Distance between p1 and p2
func EuclideanDistance(p1, p2 Point) float64 {
	disX := p2.X - p1.X
	disY := p2.Y - p1.Y
	dbs := math.Pow(float64(disX), 2) + math.Pow(float64(disY), 2)
	dis := math.Sqrt(dbs)
	return dis
}

// ManhattanDistance will measure Manhattan Distance between p1 and p2
// |x2 - x1| + |y2 - y1|
func ManhattanDistance(p1, p2 Point) float64 {
	disX := p2.X - p1.X
	disY := p2.Y - p1.Y

	return math.Abs(float64(disX)) + math.Abs(float64(disY))
}

// KNN the k nearest neighbors
type KNN struct {
	k        int
	dataSets []LabelledDataset
	classes  []int
}

// NewKNN the KNN constructor
func NewKNN(k int) *KNN {
	return &KNN{
		k:        k,
		dataSets: nil,
	}
}

func (k *KNN) sort() {
	for i := 0; i < len(k.dataSets)-1; i++ {

		minIdx := i
		for j := i + 1; j < len(k.dataSets); j++ {
			if k.dataSets[j].Distance < k.dataSets[minIdx].Distance {
				minIdx = j
			}
		}

		if minIdx != i {
			temp := k.dataSets[minIdx]
			k.dataSets[minIdx] = k.dataSets[i]
			k.dataSets[i] = temp

		}
	}
}

func (k *KNN) measureDistance(p *Point) {
	for i := 0; i < len(k.dataSets); i++ {
		dataSet := k.dataSets[i]
		distance := EuclideanDistance(*p, dataSet.Point)
		manhattanDistance := ManhattanDistance(*p, dataSet.Point)
		cosDis := CosineDistance(*p, dataSet.Point)
		fmt.Printf("(x1: %d - x2:%d) | (y1: %d - y2: %d) | euc dist: %f | man dist: %f | cos dist: %f | class: %d\n",
			p.X, dataSet.Point.X,
			p.Y, dataSet.Point.Y,
			distance,
			manhattanDistance,
			cosDis,
			dataSet.Class)

		// set Distance to Log 10 scale
		k.dataSets[i].Distance = math.Log10(distance)
	}
}

// FitData will fit data
func (k *KNN) FitData(dataSets Points, classes []int) error {
	if len(dataSets) != len(classes) {
		return errors.New("error fitting data: datasets size must equal to classes size")
	}

	var (
		labelledDatasets []LabelledDataset
	)

	for i := len(dataSets) - 1; i >= 0; i-- {
		p := dataSets[i]
		class := classes[i]

		labelledData := LabelledDataset{Class: class, Point: p}
		labelledDatasets = append(labelledDatasets, labelledData)
	}

	k.dataSets = labelledDatasets
	k.classes = classes
	return nil
}

// Plot will plot data to cartesian plane
func (k *KNN) Plot(cartesian *Cartesian) {
	var (
		clasessMap = make(map[int]int)
	)

	for _, class := range k.classes {
		clasessMap[class] = class
	}

	for i := len(k.dataSets) - 1; i >= 0; i-- {
		ds := k.dataSets[i]

		cartesian.Plot(fmt.Sprintf("%d ", clasessMap[ds.Class]), &ds.Point)
	}
}

// Predict will classify the p
func (k *KNN) Predict(p *Point) []LabelledDataset {
	fmt.Println("---------------- sorted dataset in logarithmic scale -----------------")

	// measure distance
	k.measureDistance(p)

	// sort
	k.sort()

	for i := 0; i < len(k.dataSets); i++ {
		dataSet := k.dataSets[i]
		fmt.Printf("(x1: %d - x2:%d) | (y1: %d - y2: %d) | dist: %f | class: %d\n",
			p.X, dataSet.Point.X,
			p.Y, dataSet.Point.Y,
			dataSet.Distance,
			dataSet.Class)
	}

	fmt.Printf("----------------- classifed data with %d K ----------------\n", k.k)
	classifiedDatas := k.dataSets[:k.k]

	return classifiedDatas
}

// Cartesian represent cartesian plane
type Cartesian struct {
	plane  [][]string
	width  int
	height int
	s      string
}

// NewCartesianPlane the Cartesian constructor
func NewCartesianPlane(s string, width, height int) *Cartesian {
	plane := make([][]string, height+1)
	for y := 0; y < height; y++ {
		plane[y] = make([]string, width)
	}

	return &Cartesian{
		plane:  plane,
		width:  width,
		height: height,
		s:      s,
	}
}

// Draw will draw the plane
func (c *Cartesian) Draw() {
	for y := c.height; y >= 0; y-- {
		var xCoord []string
		for x := 0; x <= c.width; x++ {
			xCoord = append(xCoord, c.s)
		}
		c.plane[y] = xCoord
	}

}

// Plot will draw w to the point p
func (c *Cartesian) Plot(w string, p *Point) {
	c.plane[p.Y][p.X] = w
}

// WriteTo will write cartesian plane to the dst
func (c *Cartesian) WriteTo(dst io.Writer) (int64, error) {
	var written int64 = 0
	for i := len(c.plane)-1; i >= 0 ; i-- {
		ss := strings.Join(c.plane[i], " ")

		w, err := dst.Write([]byte(ss))
		if err != nil {
			return written, err
		}

		written = written + int64(w)
		w, err = dst.Write([]byte{0xD, 0xA})
		if err != nil {
			return written, err
		}

		written = written + int64(w)
	}

	return written, nil
}

// Plane will return plane
func (c *Cartesian) Plane() [][]string {
	return c.plane
}

package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

// Cluster represent kmean clusters
type Cluster map[int]Points

// Point represent cartesian Point
type Point struct {
	X int
	Y int
}

// Points represent collections of point
type Points []Point

// MeanPoints will return Mean of multiple point
func (ps Points) MeanPoints() Point {
	length := float64(len(ps))
	var x float64 = 0.0
	var y float64 = 0.0

	for _, p := range ps {
		x += float64(p.X)
		y += float64(p.Y)
	}

	return Point{X: int(math.Round(x / length)), Y: int(math.Round(y / length))}
}

func (p Point) Draw(s string) string {
	return fmt.Sprintf("%s", s)
}

// Number the base number type
type Number interface {
	int |
		int8 |
		int16 |
		int32 |
		int64 |
		uint |
		uint8 |
		uint16 |
		uint32 |
		uint64 |
		~float32 |
		~float64
}

// Argmax
func Argmax[T Number](ts []T) int {
	index := 0
	value := ts[0]
	for i, v := range ts {
		if v > value {
			index = i
			value = v
		}
	}

	return index
}

// Argmin
func Argmin[T Number](ts []T) int {
	index := 0
	value := ts[0]
	for i, v := range ts {
		if v < value {
			index = i
			value = v
		}
	}

	return index
}

// EuclideanDistance measure Euclidean Distance between p1 and p2
func EuclideanDistance(p1, p2 Point) float64 {
	disX := p2.X - p1.X
	disY := p2.Y - p1.Y
	dbs := math.Pow(float64(disX), 2) + math.Pow(float64(disY), 2)
	dis := math.Sqrt(dbs)
	return dis
}

// ManhattanDistance measure Manhattan Distance between p1 and p2
// |x2 - x1| + |y2 - y1|
func ManhattanDistance(p1, p2 Point) float64 {
	disX := p2.X - p1.X
	disY := p2.Y - p1.Y

	return math.Abs(float64(disX)) + math.Abs(float64(disY))
}

// KMean the kmeans
type KMean struct {
	k          int
	maxCompute int
	clusters   map[int]Points
	centroids  Points
}

// NewKMean the KMean constructor
func NewKMean(k int, maxCompute int) *KMean {
	clusters := make(Cluster)

	return &KMean{
		k:          k,
		maxCompute: maxCompute,
		clusters:   clusters,
	}
}

// initCentroids will choose random point centroid for each cluster
func (kmean *KMean) initCentroids(dataSets Points) {
	rand.Seed(time.Now().UnixNano())

	// var (
	// 	minX int = dataSets[0].X
	// 	maxX int
	// 	minY int = dataSets[0].Y
	// 	maxY int
	// )

	// for _, d := range dataSets {
	// 	if d.X < minX {
	// 		minX = d.X
	// 	}

	// 	if d.X > maxX {
	// 		maxX = d.X
	// 	}

	// 	if d.Y < minY {
	// 		minY = d.Y
	// 	}

	// 	if d.Y > maxY {
	// 		maxY = d.Y
	// 	}
	// }

	var (
		minX int = 0
		maxX int = 60
		minY int = 0
		maxY int = 60
	)

	for i := 0; i < kmean.k; i++ {
		point := Point{
			X: int(math.Abs(float64(rand.Intn(maxX-minX+1) - minX))),
			Y: int(math.Abs(float64(rand.Intn(maxY-minY+1) - minY))),
		}

		kmean.centroids = append(kmean.centroids, point)
	}

	log.Println(kmean.centroids)
}

func (kmean *KMean) reassignCentroids(it int) {
	for k, _ := range kmean.centroids {
		log.Println("iteration: ", it)
		var (
			countPoint float64

			newCentroidX float64
			newCentroidY float64
		)

		clusters := kmean.clusters[k]
		for _, point := range clusters {
			newCentroidX += float64(point.X)
			newCentroidY += float64(point.Y)

			countPoint += 1
		}

		if countPoint > 0 {
			newCentroidX = newCentroidX / countPoint
			newCentroidY = newCentroidY / countPoint
		}

		log.Println("newCentroidX: ", newCentroidX, "countPoint: ", countPoint)
		log.Println("newCentroidY: ", newCentroidY, "countPoint: ", countPoint)

		kmean.centroids[k] = Point{X: int(math.Round(newCentroidX)), Y: int(math.Round(newCentroidY))}
	}
}

// Compute will compute the cluster
func (kmean *KMean) Compute(dataSets Points) {
	kmean.initCentroids(dataSets)

	it := 0
	for it < kmean.maxCompute {

		oldCentroids := kmean.centroids
		for j := 0; j < len(dataSets); j++ {
			p := dataSets[j]

			// to initialize the minimum distance:
			// - set minimumIndex to 0
			// - set minDist to oldCentroids at index 0
			minimumIndex := 0
			minDist := EuclideanDistance(p, oldCentroids[0])

			for i := 1; i < len(oldCentroids); i++ {
				centroid := oldCentroids[i]
				nextOldCentroid := EuclideanDistance(p, centroid)

				if nextOldCentroid < minDist {
					minimumIndex = i
					minDist = nextOldCentroid
				}
			}

			kmean.clusters[minimumIndex] = append(kmean.clusters[minimumIndex], p)

		}

		// reassign centroids
		kmean.reassignCentroids(it)

		if reflect.DeepEqual(oldCentroids, kmean.centroids) {
			log.Println("solution found")
			log.Println("old centroids: ", oldCentroids, " | new centroids: ", kmean.centroids)
			break
		}

		it++
	}

}

// WCSS Within Cluster Sum of Squares used by Elbow Method
func (kmean *KMean) WCSS() float64 {
	var clusterSum float64
	for k, centroid := range kmean.centroids {
		var sumOfSquare float64
		clusters := kmean.clusters[k]
		for _, point := range clusters {
			distanceSquare := math.Pow(EuclideanDistance(point, centroid), 2)
			sumOfSquare += distanceSquare
		}

		clusterSum += sumOfSquare
	}

	return clusterSum
}

// Predict will predict labels of p
func (kmean *KMean) Predict(p *Point) int {
	var distances []int
	// int(math.Round(newCentroidX))
	for _, centroid := range kmean.centroids {
		distance := EuclideanDistance(*p, centroid)
		distances = append(distances, int(math.Round(distance)))
	}

	log.Println("Predict| distances: ", distances)

	return Argmin(distances)
}

// Clusters will return kmean.cluster
func (kmean *KMean) Clusters() Cluster {
	return kmean.clusters
}

// Centroids will return kmean.centroids
func (kmean *KMean) Centroids() Points {
	return kmean.centroids
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

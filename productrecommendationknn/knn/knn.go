package knn

import (
	"log"
	"math"
)

// LabelledDataset represent labelled dataset
type LabelledDataset struct {
	Product  *Product
	Distance float64
}

// LabelledDatasets represent labelled dataset collections
type LabelledDatasets []LabelledDataset

// KNN the k nearest neighbors
type KNN struct {
	dataSets LabelledDatasets
}

// NewKNN the KNN constructor
func NewKNN(dataSets LabelledDatasets) *KNN {
	return &KNN{
		dataSets: dataSets,
	}
}

// Train will train the KNN
func (k *KNN) Train() {

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

func (k *KNN) measureDistance(kSize int, product *Product) {
	for i := 0; i < len(k.dataSets); i++ {
		dataSet := k.dataSets[i]
		distance := EuclideanDistance(product.CategoryNum, dataSet.Product.CategoryNum)
		// manhattanDistance := ManhattanDistance(product.CategoryNum, dataSet.Product.CategoryNum)
		// log.Printf("product name: %s | euc dist: %f | man dist: %f\n",
		// 	dataSet.Product.ProductName,
		// 	distance,
		// 	manhattanDistance)

		// set Distance to Log 10 scale
		if distance == 0.0 {
			k.dataSets[i].Distance = 0.0
		} else {
			k.dataSets[i].Distance = math.Log10(distance)
		}

		// k.dataSets[i].Distance = distance
	}
}

// Classify will classify the p
func (k *KNN) Classify(kSize int, product *Product) LabelledDatasets {
	if kSize <= 0 {
		kSize = 1
	}
	// measure distance
	k.measureDistance(kSize, product)

	log.Println("---------------- sorted dataset in logarithmic scale -----------------")

	// sort
	k.sort()

	// for i := 0; i < len(k.dataSets); i++ {
	// 	dataSet := k.dataSets[i]
	// 	log.Printf("dist: %f\n", dataSet.Distance)
	// }

	log.Printf("----------------- classifed data with %d K ----------------\n", kSize)
	classifiedDatas := k.dataSets[:kSize]

	return classifiedDatas
}

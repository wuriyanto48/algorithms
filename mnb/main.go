package main

import (
	"fmt"
	"regexp"
	"strings"
	"math"
)

type Dataset struct {
	c    int
	docs string
}

type Mnb struct {
}

func Split(text string) []string {
	pattern := "[^(){}|,.&%#@~!?=_<>/\\-\\^\\[\\]\\+\\*\\s\\t\\n]+"
	re := regexp.MustCompile(pattern)
	return re.FindAllString(text, -1)
}

// https://www.mathsisfun.com/data/probability-false-negatives-positives.html
func main() {

	trainingDatas := []Dataset{
		{
			c:    1,
			docs: "chinese chinese beijing",
		},
		{
			c:    1,
			docs: "chinese, chinese sanghai",
		},
		{
			c:    1,
			docs: "chinese macao",
		},

		{
			c:    2,
			docs: "tokyo japan chinese",
		},
	}

	testData := "chinese chinese chinese tokyo japan indonesia"

	labelClass := []string{"chinese", "japan"}

	// training
	// size of the document training
	sizeOfDocumentTraining := float64(len(trainingDatas))
	var (
		alpha            float64 = 1.0
		priorLabelClass1 float64
		priorLabelClass2 float64
		numOfVocabClass1 float64
		numOfVocabClass2 float64
	)

	vocabsLabelClass1 := make(map[string]float64)
	vocabsLabelClass2 := make(map[string]float64)
	for _, td := range trainingDatas {
		docs := td.docs
		docsSplits := Split(docs)
		for _, ds := range docsSplits {
			vocabsLabelClass1[ds] = alpha
			vocabsLabelClass2[ds] = alpha
		}
	}

	for _, td := range trainingDatas {
		docs := td.docs
		docsSplits := Split(docs)

		if td.c == 1 {
			priorLabelClass1 = priorLabelClass1 + 1

			for _, ds := range docsSplits {
				vocabsLabelClass1[ds]++
				numOfVocabClass1++
			}

		}

		if td.c == 2 {
			priorLabelClass2 = priorLabelClass2 + 1

			for _, ds := range docsSplits {
				vocabsLabelClass2[ds]++
				numOfVocabClass2++
			}
		}
	}

	// priors probability
	priorLabelClass1 = priorLabelClass1 / sizeOfDocumentTraining
	priorLabelClass2 = priorLabelClass2 / sizeOfDocumentTraining

	// compute Conditional Probabilities
	
	for k, _ := range vocabsLabelClass1 {
		vocabsLabelClass1[k] = vocabsLabelClass1[k] / (numOfVocabClass1 + float64(len(vocabsLabelClass1)))
	}

	for k, _ := range vocabsLabelClass2 {
		vocabsLabelClass2[k] = vocabsLabelClass2[k] / (numOfVocabClass2 + float64(len(vocabsLabelClass2)))
	}

	fmt.Println("prior for label class 1: ", priorLabelClass1)
	fmt.Println("prior for label class 2: ", priorLabelClass2)

	fmt.Println("num of vocabs class 1: ", numOfVocabClass1, " ", len(vocabsLabelClass1))
	fmt.Println("num of vocabs class 2: ", numOfVocabClass2, " ", len(vocabsLabelClass2))

	fmt.Println(vocabsLabelClass1)
	fmt.Println(vocabsLabelClass2)

	fmt.Println(testData)

	// predict

	var (
		resultLabel1 float64
		resultLabel2 float64
	)

	resultLabel1 = priorLabelClass1
	resultLabel2 = priorLabelClass2

	testDataSplit := Split(testData)

	// predict label 1
	for _, w := range testDataSplit {
		w = strings.ToLower(w)

		resultLabel1 = resultLabel1 * vocabsLabelClass1[w]
	}

	for _, w := range testDataSplit {
		w = strings.ToLower(w)

		resultLabel2 = resultLabel2 * vocabsLabelClass2[w]
	}

	fmt.Println("result label 1: ", resultLabel1, " ", math.Log(resultLabel1))
	fmt.Println("result label 2: ", resultLabel2, " ", math.Log(resultLabel2))

	maxIndex := Argmax[float64]([]float64{math.Log(resultLabel1), math.Log(resultLabel2)})

	fmt.Println("classified label: ", labelClass[maxIndex])

}

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

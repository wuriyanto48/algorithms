package main

import (
	"errors"
	"fmt"
	"math"
)

// K K-factor
const K = 30

func main() {
	dataStore := NewHMDataStore()
	chelsea := &Data{
		Name:          "Chelsea",
		CurrentRating: 1000,
		Win:           0,
		Lose:          0,
		Play:          0,
	}

	mu := &Data{
		Name:          "MU",
		CurrentRating: 1000,
		Win:           0,
		Lose:          0,
		Play:          0,
	}

	arsenal := &Data{
		Name:          "Arsenal",
		CurrentRating: 1000,
		Win:           0,
		Lose:          0,
		Play:          0,
	}

	liverpool := &Data{
		Name:          "Liverpool",
		CurrentRating: 1000,
		Win:           0,
		Lose:          0,
		Play:          0,
	}

	dataStore.Add(chelsea)
	dataStore.Add(mu)
	dataStore.Add(arsenal)
	dataStore.Add(liverpool)

	// match
	fmt.Println("------ Match -------")
	dataStore.Match("MU", "Chelsea", "Chelsea")
	dataStore.Match("MU", "Arsenal", "MU")
	dataStore.Match("MU", "Liverpool", "Liverpool")

	dataStore.Match("Arsenal", "Liverpool", "Liverpool")
	dataStore.Match("Arsenal", "Chelsea", "Chelsea")

	dataStore.Match("Liverpool", "Chelsea", "Liverpool")

	dataStore.Match("Chelsea", "Arsenal", "Arsenal")
	dataStore.Match("MU", "Arsenal", "Arsenal")

	fmt.Println()
	fmt.Println("------ Rating -------")
	for _, v := range dataStore.GetCurrentList() {
		fmt.Println("Team: ", v.Name, " | ", "Rating: ", math.Round(v.CurrentRating), " | ", "Play Count: ", v.Play, " | ", "Win: ", v.Win, " | ", "Lose: ", v.Lose)
	}

}

// Probability function
// r1 = rating 1
// r2 = rating 2
func Probability(r1, r2 float64) float64 {
	return 1.0 / (1.0 + math.Pow(10, (r2-r1)/400))
}

// GetNewRating function
// oldRating = the old rating
// prob = probability of the rating
func GetNewRating(oldRating, prob float64, sa int) float64 {
	return oldRating + K*(float64(sa)-prob)
}

type Data struct {
	Name          string
	CurrentRating float64
	Win           int
	Lose          int
	Play          int
}

type DataStore[K any, V any] interface {
	Add(V) V
	Update(V) V
	Find(K) V
}

type HMDataStore struct {
	data map[string]*Data
}

func NewHMDataStore() *HMDataStore {
	data := make(map[string]*Data)
	return &HMDataStore{
		data: data,
	}
}

func (ds *HMDataStore) Add(data *Data) *Data {
	ds.data[data.Name] = data
	return data
}

func (ds *HMDataStore) Update(data *Data) *Data {
	_, ok := ds.data[data.Name]
	if !ok {
		return ds.Add(data)
	}

	ds.data[data.Name] = data
	return data
}

func (ds *HMDataStore) Find(key string) *Data {
	data, ok := ds.data[key]
	if !ok {
		return nil
	}

	return data
}

func (ds *HMDataStore) Match(a string, b string, winner string) error {
	dataA := ds.Find(a)
	if dataA == nil {
		return errors.New("data A not found")
	}

	dataB := ds.Find(b)
	if dataB == nil {
		return errors.New("data B not found")
	}

	// probability
	var pA float64 = Probability(dataA.CurrentRating, dataB.CurrentRating)
	var pB float64 = Probability(dataB.CurrentRating, dataA.CurrentRating)

	if winner == dataA.Name {
		var rnA float64 = GetNewRating(dataA.CurrentRating, pA, 1)
		var rnB float64 = GetNewRating(dataB.CurrentRating, pB, 0)

		dataA.CurrentRating = rnA
		dataB.CurrentRating = rnB

		dataA.Win = dataA.Win + 1
		dataB.Lose = dataB.Lose + 1

		dataA.Play = dataA.Play + 1
		dataB.Play = dataB.Play + 1

		ds.Update(dataA)
		ds.Update(dataB)
	} else if winner == dataB.Name {
		var rnA float64 = GetNewRating(dataA.CurrentRating, pA, 0)
		var rnB float64 = GetNewRating(dataB.CurrentRating, pB, 1)

		dataA.CurrentRating = rnA
		dataB.CurrentRating = rnB

		dataA.Lose = dataA.Lose + 1
		dataB.Win = dataB.Win + 1

		dataA.Play = dataA.Play + 1
		dataB.Play = dataB.Play + 1

		ds.Update(dataA)
		ds.Update(dataB)
	} else {
		return errors.New("invalid winner input")
	}

	return nil
}

func (ds *HMDataStore) GetCurrentList() []*Data {
	var datas []*Data

	for _, v := range ds.data {
		datas = append(datas, v)
	}

	for i := 0; i < len(datas)-1; i++ {
		maxIdx := i
		for j := i + 1; j < len(datas); j++ {
			if datas[i].CurrentRating < datas[j].CurrentRating {
				maxIdx = j
			}
		}

		temp := datas[i]
		datas[i] = datas[maxIdx]
		datas[maxIdx] = temp
	}

	return datas
}
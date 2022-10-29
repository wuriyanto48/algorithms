package main

import (
	"fmt"
	"math"
)

// Point represent cartesian Point
type Point struct {
	X float64
	Y float64
}

// Points represent collections of point
type Points []Point

// MeanPoints will return Mean of multiple point
func (ps Points) MeanXY() (float64, float64) {
	length := float64(len(ps))
	var x float64 = 0.0
	var y float64 = 0.0

	for _, p := range ps {
		x += p.X
		y += p.Y
	}

	return (x / length), (y / length)
}

// SubXYByMean Substract X and Y with the Mean
// x - x̄
// y - ȳ
func (ps Points) SubXYByMean() ([]float64, []float64) {
	var (
		xSS []float64
		ySS []float64
	)

	xMean, yMean := ps.MeanXY()

	for _, p := range ps {
		xS := p.X - xMean
		yS := p.Y - yMean

		xSS = append(xSS, xS)
		ySS = append(ySS, yS)
	}

	return xSS, ySS
}

// XSubMeanMulYSubMean (x-x̄)(y-ȳ)
func (ps Points) XSubMeanMulYSubMean() []float64 {
	var (
		res []float64
	)

	xSm, ySm := ps.SubXYByMean()

	for i := 0; i < len(xSm); i++ {
		xyMulRes := xSm[i] * ySm[i]
		res = append(res, xyMulRes)
	}

	return res
}

// XSubMeanMulYSubMeanSum the sum of (x-x̄)(y-ȳ)
func (ps Points) XSubMeanMulYSubMeanSum() float64 {
	var res float64
	xMys := ps.XSubMeanMulYSubMean()

	for _, xy := range xMys {
		res += xy
	}

	return res
}

// SumOfSquareMeanXY the Sum of Square around the Mean X and Y
//  (x-x̄) ^ 2
// (y-ȳ) ^ 2
func (ps Points) SumOfSquareMeanXY() (float64, float64) {
	var xSSM float64 = 0.0
	var ySSM float64 = 0.0

	xSS, ySS := ps.SubXYByMean()
	for _, xS := range xSS {
		xSSM += math.Pow(xS, 2)
	}

	for _, yS := range ySS {
		ySSM += math.Pow(yS, 2)
	}

	return xSSM, ySSM
}

// GetYIntercept return Y intercept
func (ps Points) GetYIntercept() float64 {
	xMean, yMean := ps.MeanXY()

	return yMean - (ps.Slope() * xMean)
}

// Slope will return Slope
func (ps Points) Slope() float64 {
	ssX, _ := ps.SumOfSquareMeanXY()
	soXmulY := ps.XSubMeanMulYSubMeanSum()

	return soXmulY / ssX
}

// Predict will predict the future value
func (ps Points) Predict(x float64) float64 {
	yIntercept := ps.GetYIntercept()
	slope := ps.Slope()

	return yIntercept + (slope * x)
}

// RSquared will return R Squared (Σ(ŷ - ȳ) ^ 2 / Σ(y-ȳ) ^ 2)
func (ps Points) RSquared() float64 {
	var (
		sumOfSquareYAroundFit float64
	)

	_, yMean := ps.MeanXY()
	_, ssY := ps.SumOfSquareMeanXY()
	for _, p := range ps {
		yRes := ps.Predict(p.X) - yMean
		sumOfSquareYAroundFit += math.Pow(yRes, 2)
	}

	return (sumOfSquareYAroundFit / ssY)
}

// StdError the Standard Error of the Estimate (√ Σ(ŷ - y) ^ 2 / N - 2)
func (ps Points) StdError() float64 {
	var (
		lenOfObservations = float64(len(ps))
		sumOfSquareYAroundFit float64
	)

	for _, p := range ps {
		yRes := ps.Predict(p.X) - p.Y
		sumOfSquareYAroundFit += math.Pow(yRes, 2)
	}

	return math.Sqrt(sumOfSquareYAroundFit / (lenOfObservations - 2))
}

func (p Point) Draw(s string) string {
	return fmt.Sprintf("%s", s)
}

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	// "math"
	"math/cmplx"
	"os"
)

func main() {
	fileName := "out.png"

	var width int = 1080
	var height int = 1080

	var rl float64 = -0.74543 // real number
	var im float64 = 0.11301 // imaginary number

	outFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() { outFile.Close() }()

	err = GenerateJuliaSet(complex(rl, im), width, height, outFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GenerateJuliaSet(c complex128, width int, height int, out io.Writer) error {
	from := image.Point{0, 0}
	to := image.Point{width, height}

	imageOut := image.NewRGBA(image.Rectangle{from, to})

	backgroundColor := color.RGBA{84, 226, 247, 0xff}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			imageOut.Set(x, y, backgroundColor)
		}
	}

	// ----------------- draw circle -----------------

	// r := 150
	// centerPoint := width / 2

	// lineColor := color.RGBA{117, 5, 18, 0xff}

	// var xx float64 = 0.0
	// var yy float64 = 0.0

	// for t := 0; t <= 360; t++ {
	// 	xx = float64(centerPoint) + math.Round(float64(r)*math.Cos(float64(t)))
	// 	yy = float64(centerPoint) + math.Round(float64(r)*math.Sin(float64(t)))

	// 	imageOut.Set(int(math.Abs(xx)), int(math.Abs(yy)), lineColor)
	// }

	// ----------------- end draw circle -----------------

	// ----------------- draw julia fractal -----------------

	var scaleX float64 = 3.0 / float64(width)
	var scaleY float64 = 3.0 / float64(height)

	var escapeRadius float64 = 3.5
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			iteration := 1000

			var zx float64 = scaleX*float64(x) - 1.5 // zx represents the real part.
			var zy float64 = scaleY*float64(y) - 1.5 // zy represents the imaginary part

			i := 0

			var z = complex(zx, zy)
			// math.Sqrt(zx*zx+zy*zy) equivalent to cmplx.Abs()
			for cmplx.Abs(z) < escapeRadius && i < iteration {
				// xTemp := zx*zx - zy*zy

				// zy = 2.0*zx*zy + im
				// zx = xTemp + rl

				z = z*z + c
				i = i + 1
			}

			oR, oG, oB, oA := imageOut.At(x, y).RGBA()
			var _ uint8 = uint8(oR) // red
			var _ uint8 = uint8(oG) // green
			var _ uint8 = uint8(oB) // blue
			var a uint8 = uint8(oA) // alpha

			fColor := color.RGBA{uint8(i << 10), uint8(i << 50), uint8(i * 5), a}
			bColor := color.RGBA{13, 12, 13, a}

			if i == iteration {
				imageOut.Set(x, y, bColor)
			} else {
				imageOut.Set(x, y, fColor)
			}

			// imageOut.Set(x, y, fColor)
		}

	}

	// ----------------- end draw julia fractal -----------------

	err := png.Encode(out, imageOut)
	if err != nil {
		return err
	}

	return nil
}

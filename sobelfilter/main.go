package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"os"
)

func main() {
	file, err := os.Open("input.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() { file.Close() }()

	outFile, err := os.Create("out.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() { outFile.Close() }()

	grayscaleImage, err := GrayScaleFilter(256, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// img, format, err := image.Decode(file)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	// fmt.Println("format: ", format)
	bounds := grayscaleImage.Bounds()

	width := bounds.Max.X
	height := bounds.Max.Y

	fmt.Println("width: ", width)
	fmt.Println("height: ", height)

	sobelX := [][]int32{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}

	// sobelY := [][]int32{
	// 	{-1, -2, -1},
	// 	{0, 0, 0},
	// 	{1, 2, 1},
	// }

	imageOut := image.NewRGBA(image.Rectangle{bounds.Min, bounds.Max})

	// ignore last and first rows / x
	// ignore last and first column / y
	for x := 1; x < width-1; x++ {
		for y := 1; y < height-1; y++ {

			var (
				sumRed   int32 = 0
				sumGreen int32 = 0
				sumBlue  int32 = 0

				denominator int32 = 0
			)

			var xPos int = -1
			for ky := 0; ky < 3; ky++ {
				var yPos int = -1
				for kx := 0; kx < 3; kx++ {
					// oR, oG, oB, _ := img.At(x+xPos, y+yPos).RGBA()
					cc := color.RGBAModel.Convert(grayscaleImage.At(x+xPos, y+yPos)).(color.RGBA)
					// fmt.Println(cc.R, " ", cc.G, " ", cc.B)

					var r int32 = int32(cc.R) // red
					var g int32 = int32(cc.G) // green
					var b int32 = int32(cc.B) // blue

					// warning: this is will messed up your terminal
					// fmt.Println(r, " ", cc.R, " | ", g, " ",cc.G, " | ", b, " ", cc.B, "|")

					kernel := sobelX[ky][kx]
					// fmt.Print(kernel, " ")

					sumRed += kernel * r
					sumGreen += kernel * g
					sumBlue += kernel * b

					denominator += kernel
					yPos = yPos + 1
				}
				// fmt.Println()

				xPos = xPos + 1
			}

			if denominator <= 0 {
				denominator = 1
			}

			sumRed = sumRed / denominator
			sumGreen = sumGreen / denominator
			sumBlue = sumBlue / denominator

			if sumRed > 0xff || sumGreen > 0xff || sumBlue > 0xff {
				fmt.Println(sumRed, " ", sumGreen, " ", sumBlue)
			}

			if sumRed > 0xFF {
				sumRed = 0xFF
			} else if sumRed < 0 {
				sumRed = 0
			}

			if sumGreen > 0xFF {
				sumGreen = 0xFF
			} else if sumGreen < 0 {
				sumGreen = 0
			}

			if sumBlue > 0xFF {
				sumBlue = 0xFF
			} else if sumBlue < 0 {
				sumBlue = 0
			}

			newValue := color.RGBA{
				uint8(sumRed),
				uint8(sumGreen),
				uint8(sumBlue), 0xff,
			}

			imageOut.Set(x, y, newValue)
		}
	}

	err = png.Encode(outFile, imageOut)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GrayScaleFilter(numOfShades int, input io.Reader) (image.Image, error) {
	// decode image file to Image
	img, _, err := image.Decode(input)
	if err != nil {
		return nil, err
	}

	// get rectangle representation of image X = Width, Y = Height
	bounds := img.Bounds()
	point := bounds.Size()

	rect := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{point.X, point.Y},
	}

	imgOut := image.NewRGBA(rect)

	// todo
	//imgOut.Set(point.X, point.Y, color.Black)

	// convertionFactor := 255 / (numOfShades - 1)

	// loop to every pixel
	for x := 0; x < point.X; x++ {

		// X's Y
		for y := 0; y < point.Y; y++ {

			pixel := img.At(x, y)

			// get original red, green, blue and alpha
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)

			// http://en.wikipedia.org/wiki/Luma_%28video%29
			red := float64(originalColor.R) * 0.299
			green := float64(originalColor.G) * 0.587
			blue := float64(originalColor.B) * 0.114

			gray := red + green + blue

			// construct new color based on above calculation
			col := color.RGBA{
				R: uint8(gray),
				G: uint8(gray),
				B: uint8(gray),
				A: originalColor.A,
			}

			imgOut.Set(x, y, col)
		}
	}

	return imgOut, nil
}

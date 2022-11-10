package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
)

// https://tech-algorithm.com/articles/boxfiltering/
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

	img, format, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("format: ", format)
	bounds := img.Bounds()

	width := bounds.Max.X
	height := bounds.Max.Y

	fmt.Println("width: ", width)
	fmt.Println("height: ", height)

	// output
	// gaussianKernels := [][]int32{
	// 	{1, 2, 1},
	// 	{2, 4, 2},
	// 	{1, 2, 1},
	// }

	embossKernels := [][]int32{
		{-18, -9, 0},
		{-9, 9, 9},
		{0, 9, 18},
	}

	// lightenkernels := [][]int32{
	// 	{0, 0, 0},
	// 	{0, 12, 0},
	// 	{0, 0, 0},
	// }

	// darkenkernels := [][]int32{
	// 	{0, 0, 0},
	// 	{0, 6, 0},
	// 	{0, 0, 0},
	// }

	// testkernels := [][]int32{
	// 	{2, 22, 1},
	// 	{22, 1, -22},
	// 	{1, -22, -2},
	// }

	// sharpenkernels := [][]int32{
	// 	{-1, -1, -1},
	// 	{-1, 9, -1},
	// 	{-1, -1, -1},
	// }

	// raiesedkernels := [][]int32{
	// 	{0, 0, -2},
	// 	{0, 2, 0},
	// 	{1, 0, 0},
	// }

	// edgeDetections := [][]int32{
	// 	{-1, -1, -1},
	// 	{-1, 8, -1},
	// 	{-1, -1, -1},
	// }

	// edgeDetections2 := [][]int32{
	// 	{0, 9, 0},
	// 	{9, -36, 9},
	// 	{0, 9, 0},
	// }

	// ridgeDetections := [][]int32{
	// 	{-1, -1, -1},
	// 	{-1, 4, -1},
	// 	{-1, -1, -1},
	// }

	// kernels := [][]int32{
	// 	{0, -1, 0},
	// 	{-1, 5, -1},
	// 	{0, -1, 0},
	// }

	// motionKernels := [][]int32{
	// 	{0, 0, 1},
	// 	{0, 0, 0},
	// 	{1, 0, 0},
	// }

	// sobelX := [][]int32{
	// 	{-1, 0, 1},
	// 	{-2, 0, 2},
	// 	{-1, 0, 1},
	// }

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
					cc := color.RGBAModel.Convert(img.At(x+xPos, y+yPos)).(color.RGBA)
					// fmt.Println(cc.R, " ", cc.G, " ", cc.B)

					var r int32 = int32(cc.R) // red
					var g int32 = int32(cc.G) // green
					var b int32 = int32(cc.B) // blue

					// warning: this is will messed up your terminal
					// fmt.Println(r, " ", cc.R, " | ", g, " ",cc.G, " | ", b, " ", cc.B, "|")

					kernel := embossKernels[ky][kx]
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

			denominator = 9

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

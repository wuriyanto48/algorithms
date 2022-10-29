package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io"
	"math"
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

	err = BoxBlur(2, file, outFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func BoxBlur(kernelRadius int, in io.Reader, out io.Writer) error {
	areaSize := int32(math.Pow(float64(2*kernelRadius+1), 2))
	kernelSize := int(math.Sqrt(float64(areaSize)))

	fmt.Println("area size: ", areaSize)

	// generate kernels
	boxKernels := make([][]int32, kernelSize)

	for i := 0; i < kernelSize; i++ {
		boxKernels[i] = make([]int32, kernelSize)
	}

	for i := 0; i < kernelSize; i++ {
		for j := 0; j < kernelSize; j++ {
			boxKernels[i][j] = 1
		}
	}

	fmt.Println("kernel: ")
	for i := 0; i < kernelSize; i++ {
		fmt.Println(boxKernels[i])
	}

	// decode image
	img, format, err := image.Decode(in)
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
	fmt.Println("area size: ", areaSize)

	imageOut := image.NewRGBA(image.Rectangle{bounds.Min, bounds.Max})

	// ignore last and first rows / x
	// ignore last and first column / y
	for x := kernelRadius; x < width-kernelRadius; x++ {
		for y := kernelRadius; y < height-kernelRadius; y++ {

			var (
				sumRed   int32 = 0
				sumGreen int32 = 0
				sumBlue  int32 = 0
			)

			var xPos int = -int(kernelRadius)
			for ky := 0; ky < kernelSize; ky++ {
				var yPos int = -int(kernelRadius)
				for kx := 0; kx < kernelSize; kx++ {
					cc := color.RGBAModel.Convert(img.At(x+xPos, y+yPos)).(color.RGBA)
					// fmt.Println(cc.R, " ", cc.G, " ", cc.B)

					var r int32 = int32(cc.R) // red
					var g int32 = int32(cc.G) // green
					var b int32 = int32(cc.B) // blue

					// warning: this is well messed up your terminal
					// fmt.Print(r, " ", g, " ", b, "|")

					kernel := boxKernels[ky][kx]

					sumRed += kernel * r
					sumGreen += kernel * g
					sumBlue += kernel * b

					yPos = yPos + 1
				}

				xPos = xPos + 1
			}

			sumRed = sumRed / areaSize
			sumGreen = sumGreen / areaSize
			sumBlue = sumBlue / areaSize

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
				uint8(sumBlue),
				0xff,
			}

			imageOut.Set(x, y, newValue)
		}
	}

	err = png.Encode(out, imageOut)
	if err != nil {
		return err
	}

	return nil
}

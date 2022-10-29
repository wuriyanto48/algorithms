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

	resultImage, err := TemperatureFilter(35000, file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = png.Encode(outFile, resultImage)
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

func TemperatureFilter(temperature float64, input io.Reader) (image.Image, error) {
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

	redAdjustment, _, blueAdjustment := ComputeFilterFromKelvin(temperature)

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

			red := int32(originalColor.R)
			green := int32(originalColor.G)
			blue := int32(originalColor.B)

			red = red + int32(redAdjustment)
			// green = green + int32(greenAdjustment)
			blue = blue - int32(blueAdjustment)

			// fmt.Println(red, " ", green, " ", blue)

			if red > 0xFF {
				red = 0xFF
			} else if red < 0 {
				red = 0
			}

			if green > 0xFF {
				green = 0xFF
			} else if green < 0 {
				green = 0
			}

			if blue > 0xFF {
				blue = 0xFF
			} else if blue < 0 {
				blue = 0
			}

			// construct new color based on above calculation
			col := color.RGBA{
				R: uint8(red),
				G: uint8(green),
				B: uint8(blue),
				A: originalColor.A,
			}

			imgOut.Set(x, y, col)
		}
	}

	return imgOut, nil
}

func ComputeFilterFromKelvin(temperature float64) (uint8, uint8, uint8) {
	if temperature < 1000 {
		temperature = 1000
	}

	if temperature > 40000 {
		temperature = 40000
	}

	temperature = temperature / 100

	var (
		red, green, blue float64
	)

	if temperature <= 66 {
		red = 0xFF
	} else {
		temperatureCalcius := temperature - 60
		temperatureCalcius = 329.698727446 * math.Pow(temperatureCalcius, -0.1332047592)

		red = temperatureCalcius

		if red > 0xFF {
			red = 0xFF
		}

		if red < 0x0 {
			red = 0x0
		}
	}

	if temperature <= 66 {
		temperatureCalcius := temperature
		temperatureCalcius = 99.4708025861*math.Log(temperatureCalcius) - 161.1195681661

		green = temperatureCalcius

		if green > 0xFF {
			green = 0xFF
		}

		if green < 0x0 {
			green = 0x0
		}
	} else {
		temperatureCalcius := temperature - 60
		temperatureCalcius = 288.1221695283 * math.Pow(temperatureCalcius, -0.0755148492)

		green = temperatureCalcius

		if green > 0xFF {
			green = 0xFF
		}

		if green < 0x0 {
			green = 0x0
		}
	}

	if temperature >= 66 {
		blue = 0xFF
	} else if temperature <= 19 {
		blue = 0x0
	} else {
		temperatureCalcius := temperature - 10
		temperatureCalcius = 138.5177312231*math.Log(temperatureCalcius) - 305.0447927307

		blue = temperatureCalcius

		if blue > 0xFF {
			blue = 0xFF
		}

		if blue < 0x0 {
			blue = 0x0
		}
	}

	return uint8(red), uint8(green), uint8(blue)
}

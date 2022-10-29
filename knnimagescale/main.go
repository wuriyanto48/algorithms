package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

// https://tech-algorithm.com/articles/nearest-neighbor-image-scaling/
// https://tech-algorithm.com/articles/bilinear-image-scaling
// https://www.researchgate.net/publication/272092207_A_Novel_Visual_Cryptographic_Method_for_Color_Images
func main() {
	file, err := os.Open("input.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() { file.Close() }()

	outFile, err := os.Create("out.txt")
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

	imageOut, _ := KNNScale(250, 100, img)

	levels := [][]byte{
		{0x20},             // Space
		{0xE2, 0x96, 0x91}, // Light Shade
		{0xE2, 0x96, 0x92}, // Medium Shade
		{0xE2, 0x96, 0x93}, // Dark Shade
		{0xE2, 0x96, 0x88}, // Full Block
	}

	imageOutBounds := imageOut.Bounds()
	imageOutWidth := imageOutBounds.Max.X
	imageOutHeight := imageOutBounds.Max.Y
	for y := 0; y < imageOutHeight; y++ {
		for x := 0; x < imageOutWidth; x++ {
			c := color.GrayModel.Convert(imageOut.At(x, y)).(color.Gray)
			lv := c.Y / 51 // 255 / 51 = 5
			if lv == 5 {
				lv = lv - 1
			}
			_, err = outFile.Write(levels[lv])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		_, err = outFile.Write([]byte{10})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// err = png.Encode(outFile, imageOut)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

}

func KNNScale(wOut int, hOut int, in image.Image) (image.Image, error) {
	from := image.Point{0, 0}
	to := image.Point{wOut, hOut}

	bounds := in.Bounds()

	width := bounds.Max.X
	height := bounds.Max.Y

	out := image.NewRGBA(image.Rectangle{from, to})

	var (
		xRatio int = ((width << 16) / wOut) + 1
		yRatio int = ((height << 16) / hOut) + 1
	)

	var (
		px int = 0
		py int = 0
	)

	for y := 0; y < hOut; y++ {
		for x := 0; x < wOut; x++ {
			px = (x * xRatio) >> 16
			py = (y * yRatio) >> 16

			out.Set(x, y, in.At(px, py))
		}
	}

	return out, nil
}

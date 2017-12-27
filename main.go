package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"
)

const (
	width, height = 800, 400
)

func main() {
	start := time.Now()
	img := render()
	elapsed := time.Since(start)
	fmt.Println("Image rendering took:", elapsed)

	f, err := os.Create("out.png")
	defer f.Close()
	if err != nil {
		panic("cannot create out.png")
	}

	png.Encode(f, img)
}

func render() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			r := float64(px) / width
			g := float64(height-py) / height
			b := 0.2

			ir := uint8(255.99 * r)
			ig := uint8(255.99 * g)
			ib := uint8(255.99 * b)

			clr := color.RGBA{R: ir, G: ig, B: ib, A: 255}

			img.Set(px, py, clr)
		}
	}

	return img
}

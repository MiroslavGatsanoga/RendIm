package rendim

import (
	"image/color"
)

type Color struct {
	R, G, B float64
}

func (c Color) Add(c2 Color) Color {
	return Color{R: c.R + c2.R, G: c.G + c2.G, B: c.B + c2.B}
}

func (c Color) Subtract(c2 Color) Color {
	return Color{R: c.R - c2.R, G: c.G - c2.G, B: c.B - c2.B}
}

func (c Color) Multiply(c2 Color) Color {
	return Color{R: c.R * c2.R, G: c.G * c2.G, B: c.B * c2.B}
}

func (c Color) MultiplyScalar(s float64) Color {
	return Color{R: c.R * s, G: c.G * s, B: c.B * s}
}

func (c Color) DivideScalar(s float64) Color {
	return Color{R: c.R / s, G: c.G / s, B: c.B / s}
}

func (c Color) ToRGBA() color.RGBA {
	ir := uint8(255.99 * c.R)
	ig := uint8(255.99 * c.G)
	ib := uint8(255.99 * c.B)

	return color.RGBA{R: ir, G: ig, B: ib, A: 255}
}

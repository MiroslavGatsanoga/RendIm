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

func (c Color) Clamp() Color {
	if c.R < 0.0 {
		c.R = 0.0
	}
	if c.G < 0.0 {
		c.G = 0.0
	}
	if c.B < 0.0 {
		c.B = 0.0
	}

	if c.R > 1.0 {
		c.R = 1.0
	}
	if c.G > 1.0 {
		c.G = 1.0
	}
	if c.B > 1.0 {
		c.B = 1.0
	}

	return c
}

func (c Color) ToRGBA() color.RGBA {
	clamped := c.Clamp()

	ir := uint8(255.99 * clamped.R)
	ig := uint8(255.99 * clamped.G)
	ib := uint8(255.99 * clamped.B)

	return color.RGBA{R: ir, G: ig, B: ib, A: 255}
}

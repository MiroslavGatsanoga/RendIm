package rendim

import (
	"testing"
)

func TestColorAdd(t *testing.T) {
	c1 := Color{R: 0.2, G: 0.3, B: 0.4}
	c2 := Color{R: 0.1, G: 0.2, B: 0.3}
	result := c1.Add(c2)
	
	expected := Color{R: 0.3, G: 0.5, B: 0.7}
	if !colorEqual(result, expected) {
		t.Errorf("Add() = (%f, %f, %f), want (0.3, 0.5, 0.7)", result.R, result.G, result.B)
	}
}

func colorEqual(c1, c2 Color) bool {
	const epsilon = 1e-10
	return abs(c1.R-c2.R) < epsilon && abs(c1.G-c2.G) < epsilon && abs(c1.B-c2.B) < epsilon
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestColorSubtract(t *testing.T) {
	c1 := Color{R: 0.5, G: 0.6, B: 0.7}
	c2 := Color{R: 0.2, G: 0.3, B: 0.4}
	result := c1.Subtract(c2)
	
	expected := Color{R: 0.3, G: 0.3, B: 0.3}
	if !colorEqual(result, expected) {
		t.Errorf("Subtract() = (%f, %f, %f), want (0.3, 0.3, 0.3)", result.R, result.G, result.B)
	}
}

func TestColorMultiply(t *testing.T) {
	c1 := Color{R: 0.5, G: 0.6, B: 0.8}
	c2 := Color{R: 0.2, G: 0.5, B: 0.25}
	result := c1.Multiply(c2)
	
	if result.R != 0.1 || result.G != 0.3 || result.B != 0.2 {
		t.Errorf("Multiply() = (%f, %f, %f), want (0.1, 0.3, 0.2)", result.R, result.G, result.B)
	}
}

func TestColorMultiplyScalar(t *testing.T) {
	c := Color{R: 0.2, G: 0.4, B: 0.6}
	result := c.MultiplyScalar(2.0)
	
	if result.R != 0.4 || result.G != 0.8 || result.B != 1.2 {
		t.Errorf("MultiplyScalar(2.0) = (%f, %f, %f), want (0.4, 0.8, 1.2)", result.R, result.G, result.B)
	}
}

func TestColorDivideScalar(t *testing.T) {
	c := Color{R: 0.4, G: 0.8, B: 1.2}
	result := c.DivideScalar(2.0)
	
	if result.R != 0.2 || result.G != 0.4 || result.B != 0.6 {
		t.Errorf("DivideScalar(2.0) = (%f, %f, %f), want (0.2, 0.4, 0.6)", result.R, result.G, result.B)
	}
}

func TestColorClamp(t *testing.T) {
	tests := []struct {
		name     string
		input    Color
		expected Color
	}{
		{"All in range", Color{0.5, 0.6, 0.7}, Color{0.5, 0.6, 0.7}},
		{"All above 1.0", Color{1.5, 2.0, 3.0}, Color{1.0, 1.0, 1.0}},
		{"All below 0.0", Color{-0.5, -1.0, -2.0}, Color{0.0, 0.0, 0.0}},
		{"Mixed", Color{-0.1, 0.5, 1.5}, Color{0.0, 0.5, 1.0}},
		{"Exact bounds", Color{0.0, 0.5, 1.0}, Color{0.0, 0.5, 1.0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Clamp()
			if result.R != tt.expected.R || result.G != tt.expected.G || result.B != tt.expected.B {
				t.Errorf("Clamp() = (%f, %f, %f), want (%f, %f, %f)",
					result.R, result.G, result.B,
					tt.expected.R, tt.expected.G, tt.expected.B)
			}
		})
	}
}

func TestColorToRGBA(t *testing.T) {
	tests := []struct {
		name     string
		color    Color
		expectedR uint8
		expectedG uint8
		expectedB uint8
	}{
		{"Black", Color{0.0, 0.0, 0.0}, 0, 0, 0},
		{"White", Color{1.0, 1.0, 1.0}, 255, 255, 255},
		{"Red", Color{1.0, 0.0, 0.0}, 255, 0, 0},
		{"Half grey", Color{0.5, 0.5, 0.5}, 127, 127, 127},
		{"Clamped high", Color{2.0, 2.0, 2.0}, 255, 255, 255},
		{"Clamped low", Color{-1.0, -1.0, -1.0}, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rgba := tt.color.ToRGBA()
			if rgba.R != tt.expectedR || rgba.G != tt.expectedG || rgba.B != tt.expectedB {
				t.Errorf("ToRGBA() = (%d, %d, %d), want (%d, %d, %d)",
					rgba.R, rgba.G, rgba.B,
					tt.expectedR, tt.expectedG, tt.expectedB)
			}
			if rgba.A != 255 {
				t.Errorf("ToRGBA() alpha = %d, want 255", rgba.A)
			}
		})
	}
}

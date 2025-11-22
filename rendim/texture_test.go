package rendim

import (
	"image"
	"image/color"
	"math"
	"testing"
)

func TestConstantTextureValue(t *testing.T) {
	clr := Color{R: 0.5, G: 0.6, B: 0.7}
	tex := ConstantTexture{color: clr}
	
	result := tex.Value(0.0, 0.0, NewVec3d(0.0, 0.0, 0.0))
	
	if result.R != 0.5 || result.G != 0.6 || result.B != 0.7 {
		t.Errorf("Value() = (%f, %f, %f), want (0.5, 0.6, 0.7)", result.R, result.G, result.B)
	}
}

func TestCheckerTextureValue(t *testing.T) {
	even := ConstantTexture{color: Color{R: 1.0, G: 1.0, B: 1.0}}
	odd := ConstantTexture{color: Color{R: 0.0, G: 0.0, B: 0.0}}
	checker := CheckerTexture{even: even, odd: odd}
	
	result1 := checker.Value(0.0, 0.0, NewVec3d(0.0, 0.0, 0.0))
	if result1.R != 1.0 || result1.G != 1.0 || result1.B != 1.0 {
		t.Errorf("Checker at (0,0,0) should be even color")
	}
	
	result2 := checker.Value(0.0, 0.0, NewVec3d(0.5, 0.5, 0.5))
	if result2.R == result1.R && result2.G == result1.G && result2.B == result1.B {
		t.Error("Checker pattern should alternate")
	}
}

func TestNoiseTextureValue(t *testing.T) {
	noiseTex := NoiseTexture{scale: 1.0}
	
	result := noiseTex.Value(0.0, 0.0, NewVec3d(1.0, 1.0, 1.0))
	
	if result.R < 0.0 || result.R > 1.0 {
		t.Errorf("Noise texture R = %f, should be in [0, 1]", result.R)
	}
	if result.G < 0.0 || result.G > 1.0 {
		t.Errorf("Noise texture G = %f, should be in [0, 1]", result.G)
	}
	if result.B < 0.0 || result.B > 1.0 {
		t.Errorf("Noise texture B = %f, should be in [0, 1]", result.B)
	}
}

func TestPerlinNoise(t *testing.T) {
	p1 := NewVec3d(1.0, 2.0, 3.0)
	noise1 := perlinNoise.Noise(p1)
	
	if math.IsNaN(noise1) || math.IsInf(noise1, 0) {
		t.Error("Perlin noise should return valid float")
	}
	
	p2 := NewVec3d(1.0, 2.0, 3.0)
	noise2 := perlinNoise.Noise(p2)
	
	if noise1 != noise2 {
		t.Error("Perlin noise should be deterministic for same input")
	}
}

func TestPerlinTurbulence(t *testing.T) {
	p := NewVec3d(1.0, 2.0, 3.0)
	turb := perlinNoise.Turbulence(p)
	
	if turb < 0.0 {
		t.Errorf("Turbulence = %f, should be >= 0", turb)
	}
	
	if math.IsNaN(turb) || math.IsInf(turb, 0) {
		t.Error("Turbulence should return valid float")
	}
}

func TestPerlinGenerate(t *testing.T) {
	vectors := perlinGenerate()
	
	if len(vectors) != 256 {
		t.Errorf("perlinGenerate() returned %d vectors, want 256", len(vectors))
	}
	
	for i, v := range vectors {
		if v.X() < -1.0 || v.X() > 1.0 {
			t.Errorf("Vector %d X component %f out of range [-1, 1]", i, v.X())
		}
		if v.Y() < -1.0 || v.Y() > 1.0 {
			t.Errorf("Vector %d Y component %f out of range [-1, 1]", i, v.Y())
		}
		if v.Z() < -1.0 || v.Z() > 1.0 {
			t.Errorf("Vector %d Z component %f out of range [-1, 1]", i, v.Z())
		}
	}
}

func TestPerlinGeneratePerm(t *testing.T) {
	perm := perlinGeneratePerm()
	
	if len(perm) != 256 {
		t.Errorf("perlinGeneratePerm() returned %d values, want 256", len(perm))
	}
	
	seen := make(map[int]bool)
	for _, val := range perm {
		if val < 0 || val >= 256 {
			t.Errorf("Permutation value %d out of range [0, 256)", val)
		}
		if seen[val] {
			t.Errorf("Duplicate value %d in permutation", val)
		}
		seen[val] = true
	}
}

func TestPermute(t *testing.T) {
	p := []int{0, 1, 2, 3, 4, 5}
	permute(p)
	
	seen := make(map[int]bool)
	for _, val := range p {
		seen[val] = true
	}
	
	if len(seen) != 6 {
		t.Error("Permute should preserve all elements")
	}
}

func TestImageTextureValue(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 128, B: 64, A: 255})
		}
	}
	
	tex := ImageTexture{image: img}
	
	result := tex.Value(0.5, 0.5, NewVec3d(0.0, 0.0, 0.0))
	
	if result.R < 0.0 || result.R > 1.0 {
		t.Errorf("Image texture R = %f, should be in [0, 1]", result.R)
	}
	if result.G < 0.0 || result.G > 1.0 {
		t.Errorf("Image texture G = %f, should be in [0, 1]", result.G)
	}
	if result.B < 0.0 || result.B > 1.0 {
		t.Errorf("Image texture B = %f, should be in [0, 1]", result.B)
	}
}

func TestGetSphereUV(t *testing.T) {
	tests := []struct {
		name string
		p    Vec3d
	}{
		{"Top", NewVec3d(0.0, 1.0, 0.0)},
		{"Bottom", NewVec3d(0.0, -1.0, 0.0)},
		{"Front", NewVec3d(0.0, 0.0, 1.0)},
		{"Back", NewVec3d(0.0, 0.0, -1.0)},
		{"Right", NewVec3d(1.0, 0.0, 0.0)},
		{"Left", NewVec3d(-1.0, 0.0, 0.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, v := getSphereUV(tt.p)
			
			if u < 0.0 || u > 1.0 {
				t.Errorf("u = %f, should be in [0, 1]", u)
			}
			if v < 0.0 || v > 1.0 {
				t.Errorf("v = %f, should be in [0, 1]", v)
			}
		})
	}
}

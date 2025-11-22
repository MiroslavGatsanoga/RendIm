package rendim

import (
	"testing"
)

func TestFfMin(t *testing.T) {
	tests := []struct {
		a, b     float64
		expected float64
	}{
		{1.0, 2.0, 1.0},
		{2.0, 1.0, 1.0},
		{-1.0, 1.0, -1.0},
		{0.0, 0.0, 0.0},
	}

	for _, tt := range tests {
		result := ffMin(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("ffMin(%f, %f) = %f, want %f", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestFfMax(t *testing.T) {
	tests := []struct {
		a, b     float64
		expected float64
	}{
		{1.0, 2.0, 2.0},
		{2.0, 1.0, 2.0},
		{-1.0, 1.0, 1.0},
		{0.0, 0.0, 0.0},
	}

	for _, tt := range tests {
		result := ffMax(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("ffMax(%f, %f) = %f, want %f", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestAABBHit(t *testing.T) {
	box := AABB{
		Min: NewVec3d(-1.0, -1.0, -1.0),
		Max: NewVec3d(1.0, 1.0, 1.0),
	}

	tests := []struct {
		name     string
		ray      Ray
		tMin     float64
		tMax     float64
		expected bool
	}{
		{
			"Ray hits box center",
			NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			0.0,
			10.0,
			true,
		},
		{
			"Ray misses box",
			NewRay(NewVec3d(-5.0, 5.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			0.0,
			10.0,
			false,
		},
		{
			"Ray inside box",
			NewRay(NewVec3d(0.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			0.0,
			10.0,
			true,
		},
		{
			"Ray hits from opposite direction",
			NewRay(NewVec3d(5.0, 0.0, 0.0), NewVec3d(-1.0, 0.0, 0.0), 0.0),
			0.0,
			10.0,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := box.hit(tt.ray, tt.tMin, tt.tMax)
			if result != tt.expected {
				t.Errorf("hit() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSurroundingBox(t *testing.T) {
	box1 := AABB{
		Min: NewVec3d(0.0, 0.0, 0.0),
		Max: NewVec3d(1.0, 1.0, 1.0),
	}
	box2 := AABB{
		Min: NewVec3d(-1.0, -1.0, -1.0),
		Max: NewVec3d(0.5, 0.5, 0.5),
	}

	result := surroundingBox(box1, box2)

	expectedMin := NewVec3d(-1.0, -1.0, -1.0)
	expectedMax := NewVec3d(1.0, 1.0, 1.0)

	if result.Min.X() != expectedMin.X() || result.Min.Y() != expectedMin.Y() || result.Min.Z() != expectedMin.Z() {
		t.Errorf("surroundingBox().Min = (%f, %f, %f), want (%f, %f, %f)",
			result.Min.X(), result.Min.Y(), result.Min.Z(),
			expectedMin.X(), expectedMin.Y(), expectedMin.Z())
	}

	if result.Max.X() != expectedMax.X() || result.Max.Y() != expectedMax.Y() || result.Max.Z() != expectedMax.Z() {
		t.Errorf("surroundingBox().Max = (%f, %f, %f), want (%f, %f, %f)",
			result.Max.X(), result.Max.Y(), result.Max.Z(),
			expectedMax.X(), expectedMax.Y(), expectedMax.Z())
	}
}

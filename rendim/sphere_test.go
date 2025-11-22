package rendim

import (
	"math"
	"testing"
)

type mockMaterial struct{}

func (m mockMaterial) Scatter(rayIn Ray, rec HitRecord, attenuation *Color) (bool, Ray) {
	return false, Ray{}
}

func (m mockMaterial) Emitted(u, v float64, p Vec3d) Color {
	return Color{0, 0, 0}
}

func TestNewSphere(t *testing.T) {
	center := NewVec3d(1.0, 2.0, 3.0)
	radius := 5.0
	mat := mockMaterial{}
	
	sphere := NewSphere(center, radius, mat)
	
	if sphere.Center.X() != 1.0 || sphere.Center.Y() != 2.0 || sphere.Center.Z() != 3.0 {
		t.Errorf("Sphere center = (%f, %f, %f), want (1.0, 2.0, 3.0)",
			sphere.Center.X(), sphere.Center.Y(), sphere.Center.Z())
	}
	
	if sphere.Radius != 5.0 {
		t.Errorf("Sphere radius = %f, want 5.0", sphere.Radius)
	}
}

func TestSphereHit(t *testing.T) {
	mat := mockMaterial{}
	sphere := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	
	tests := []struct {
		name        string
		ray         Ray
		tMin        float64
		tMax        float64
		shouldHit   bool
	}{
		{
			"Ray hits sphere from front",
			NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			0.0,
			10.0,
			true,
		},
		{
			"Ray misses sphere",
			NewRay(NewVec3d(-5.0, 5.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			0.0,
			10.0,
			false,
		},
		{
			"Ray starts inside sphere",
			NewRay(NewVec3d(0.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			0.0,
			10.0,
			true,
		},
		{
			"Ray hits but outside tMax range",
			NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			0.0,
			1.0,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hit, rec := sphere.Hit(tt.ray, tt.tMin, tt.tMax)
			if hit != tt.shouldHit {
				t.Errorf("Hit() = %v, want %v", hit, tt.shouldHit)
			}
			
			if hit {
				length := rec.Normal.Length()
				if math.Abs(length-1.0) > 1e-10 {
					t.Errorf("Normal length = %f, want 1.0", length)
				}
			}
		})
	}
}

func TestSphereBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	center := NewVec3d(2.0, 3.0, 4.0)
	radius := 1.5
	sphere := NewSphere(center, radius, mat)
	
	var box AABB
	hasBox := sphere.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("BoundingBox() should return true for sphere")
	}
	
	expectedMin := NewVec3d(0.5, 1.5, 2.5)
	expectedMax := NewVec3d(3.5, 4.5, 5.5)
	
	if box.Min.X() != expectedMin.X() || box.Min.Y() != expectedMin.Y() || box.Min.Z() != expectedMin.Z() {
		t.Errorf("BoundingBox Min = (%f, %f, %f), want (%f, %f, %f)",
			box.Min.X(), box.Min.Y(), box.Min.Z(),
			expectedMin.X(), expectedMin.Y(), expectedMin.Z())
	}
	
	if box.Max.X() != expectedMax.X() || box.Max.Y() != expectedMax.Y() || box.Max.Z() != expectedMax.Z() {
		t.Errorf("BoundingBox Max = (%f, %f, %f), want (%f, %f, %f)",
			box.Max.X(), box.Max.Y(), box.Max.Z(),
			expectedMax.X(), expectedMax.Y(), expectedMax.Z())
	}
}

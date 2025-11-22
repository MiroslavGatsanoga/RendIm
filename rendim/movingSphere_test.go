package rendim

import (
	"math"
	"testing"
)

func TestNewMovingSphere(t *testing.T) {
	mat := mockMaterial{}
	cen0 := NewVec3d(0.0, 0.0, 0.0)
	cen1 := NewVec3d(10.0, 0.0, 0.0)
	
	ms := NewMovingSphere(cen0, cen1, 0.0, 1.0, 1.0, mat)
	
	if ms.center0.X() != 0.0 || ms.center1.X() != 10.0 {
		t.Error("MovingSphere centers not set correctly")
	}
	
	if ms.time0 != 0.0 || ms.time1 != 1.0 {
		t.Error("MovingSphere times not set correctly")
	}
	
	if ms.Radius != 1.0 {
		t.Errorf("MovingSphere radius = %f, want 1.0", ms.Radius)
	}
}

func TestMovingSphereCenter(t *testing.T) {
	mat := mockMaterial{}
	cen0 := NewVec3d(0.0, 0.0, 0.0)
	cen1 := NewVec3d(10.0, 0.0, 0.0)
	ms := NewMovingSphere(cen0, cen1, 0.0, 1.0, 1.0, mat)
	
	tests := []struct {
		time     float64
		expectedX float64
	}{
		{0.0, 0.0},
		{0.5, 5.0},
		{1.0, 10.0},
	}

	for _, tt := range tests {
		center := ms.Center(tt.time)
		if math.Abs(center.X()-tt.expectedX) > 1e-10 {
			t.Errorf("Center at time %f: X = %f, want %f", tt.time, center.X(), tt.expectedX)
		}
	}
}

func TestMovingSphereHit(t *testing.T) {
	mat := mockMaterial{}
	cen0 := NewVec3d(0.0, 0.0, 0.0)
	cen1 := NewVec3d(0.0, 0.0, 0.0)
	ms := NewMovingSphere(cen0, cen1, 0.0, 1.0, 1.0, mat)
	
	tests := []struct {
		name      string
		ray       Ray
		shouldHit bool
	}{
		{
			"Ray hits sphere",
			NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.5),
			true,
		},
		{
			"Ray misses sphere",
			NewRay(NewVec3d(-5.0, 5.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.5),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hit, rec := ms.Hit(tt.ray, 0.0, 10.0)
			if hit != tt.shouldHit {
				t.Errorf("Hit() = %v, want %v", hit, tt.shouldHit)
			}
			
			if hit {
				normalLength := rec.Normal.Length()
				if math.Abs(normalLength-1.0) > 1e-10 {
					t.Errorf("Normal length = %f, want 1.0", normalLength)
				}
			}
		})
	}
}

func TestMovingSphereHitMoving(t *testing.T) {
	mat := mockMaterial{}
	cen0 := NewVec3d(0.0, 0.0, 0.0)
	cen1 := NewVec3d(10.0, 0.0, 0.0)
	ms := NewMovingSphere(cen0, cen1, 0.0, 1.0, 1.0, mat)
	
	ray := NewRay(NewVec3d(5.0, 0.0, -5.0), NewVec3d(0.0, 0.0, 1.0), 0.5)
	hit, _ := ms.Hit(ray, 0.0, 10.0)
	
	if !hit {
		t.Error("Ray should hit moving sphere at its midpoint position")
	}
}

func TestMovingSphereBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	cen0 := NewVec3d(0.0, 0.0, 0.0)
	cen1 := NewVec3d(10.0, 0.0, 0.0)
	ms := NewMovingSphere(cen0, cen1, 0.0, 1.0, 2.0, mat)
	
	var box AABB
	hasBox := ms.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("MovingSphere should have bounding box")
	}
	
	if box.Min.X() > -2.0 {
		t.Errorf("Bounding box should encompass starting position")
	}
	
	if box.Max.X() < 12.0 {
		t.Errorf("Bounding box should encompass ending position")
	}
}

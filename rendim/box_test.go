package rendim

import (
	"testing"
)

func TestNewBox(t *testing.T) {
	mat := mockMaterial{}
	p0 := NewVec3d(0.0, 0.0, 0.0)
	p1 := NewVec3d(1.0, 1.0, 1.0)
	
	box := NewBox(p0, p1, mat)
	
	if box.pMin.X() != 0.0 || box.pMin.Y() != 0.0 || box.pMin.Z() != 0.0 {
		t.Errorf("Box pMin = (%f, %f, %f), want (0.0, 0.0, 0.0)",
			box.pMin.X(), box.pMin.Y(), box.pMin.Z())
	}
	
	if box.pMax.X() != 1.0 || box.pMax.Y() != 1.0 || box.pMax.Z() != 1.0 {
		t.Errorf("Box pMax = (%f, %f, %f), want (1.0, 1.0, 1.0)",
			box.pMax.X(), box.pMax.Y(), box.pMax.Z())
	}
	
	if len(box.faces) != 6 {
		t.Errorf("Box should have 6 faces, got %d", len(box.faces))
	}
}

func TestBoxHit(t *testing.T) {
	mat := mockMaterial{}
	box := NewBox(NewVec3d(-1.0, -1.0, -1.0), NewVec3d(1.0, 1.0, 1.0), mat)
	
	tests := []struct {
		name      string
		ray       Ray
		shouldHit bool
	}{
		{
			"Ray hits box from X direction",
			NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			true,
		},
		{
			"Ray misses box",
			NewRay(NewVec3d(-5.0, 5.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			false,
		},
		{
			"Ray hits from Y direction",
			NewRay(NewVec3d(0.0, -5.0, 0.0), NewVec3d(0.0, 1.0, 0.0), 0.0),
			true,
		},
		{
			"Ray hits from Z direction",
			NewRay(NewVec3d(0.0, 0.0, -5.0), NewVec3d(0.0, 0.0, 1.0), 0.0),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hit, _ := box.Hit(tt.ray, 0.001, 10.0)
			if hit != tt.shouldHit {
				t.Errorf("Hit() = %v, want %v", hit, tt.shouldHit)
			}
		})
	}
}

func TestBoxBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	p0 := NewVec3d(1.0, 2.0, 3.0)
	p1 := NewVec3d(4.0, 5.0, 6.0)
	box := NewBox(p0, p1, mat)
	
	var aabb AABB
	hasBox := box.BoundingBox(0.0, 1.0, &aabb)
	
	if !hasBox {
		t.Error("Box should have bounding box")
	}
	
	if aabb.Min.X() != 1.0 || aabb.Min.Y() != 2.0 || aabb.Min.Z() != 3.0 {
		t.Errorf("BoundingBox Min = (%f, %f, %f), want (1.0, 2.0, 3.0)",
			aabb.Min.X(), aabb.Min.Y(), aabb.Min.Z())
	}
	
	if aabb.Max.X() != 4.0 || aabb.Max.Y() != 5.0 || aabb.Max.Z() != 6.0 {
		t.Errorf("BoundingBox Max = (%f, %f, %f), want (4.0, 5.0, 6.0)",
			aabb.Max.X(), aabb.Max.Y(), aabb.Max.Z())
	}
}

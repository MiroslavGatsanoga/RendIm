package rendim

import (
	"math"
	"testing"
)

func TestXYRectHit(t *testing.T) {
	mat := mockMaterial{}
	rect := XYRect{x0: -1.0, x1: 1.0, y0: -1.0, y1: 1.0, k: 0.0, material: mat}
	
	tests := []struct {
		name      string
		ray       Ray
		shouldHit bool
	}{
		{
			"Ray hits from front",
			NewRay(NewVec3d(0.0, 0.0, -5.0), NewVec3d(0.0, 0.0, 1.0), 0.0),
			true,
		},
		{
			"Ray hits from back",
			NewRay(NewVec3d(0.0, 0.0, 5.0), NewVec3d(0.0, 0.0, -1.0), 0.0),
			true,
		},
		{
			"Ray misses (outside bounds)",
			NewRay(NewVec3d(5.0, 5.0, -5.0), NewVec3d(0.0, 0.0, 1.0), 0.0),
			false,
		},
		{
			"Ray parallel to rectangle (in plane)",
			NewRay(NewVec3d(-2.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hit, rec := rect.Hit(tt.ray, 0.001, 10.0)
			if hit != tt.shouldHit {
				t.Errorf("Hit() = %v, want %v", hit, tt.shouldHit)
			}
			
			if hit {
				if rec.Normal.X() != 0.0 || rec.Normal.Y() != 0.0 || math.Abs(rec.Normal.Z()) != 1.0 {
					t.Errorf("XYRect normal should be (0,0,±1), got (%f,%f,%f)",
						rec.Normal.X(), rec.Normal.Y(), rec.Normal.Z())
				}
			}
		})
	}
}

func TestXYRectBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	rect := XYRect{x0: 1.0, x1: 2.0, y0: 3.0, y1: 4.0, k: 5.0, material: mat}
	
	var box AABB
	hasBox := rect.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("XYRect should have bounding box")
	}
	
	if box.Min.X() != 1.0 || box.Min.Y() != 3.0 {
		t.Errorf("Bounding box Min X,Y = (%f, %f), want (1.0, 3.0)", box.Min.X(), box.Min.Y())
	}
	
	if box.Max.X() != 2.0 || box.Max.Y() != 4.0 {
		t.Errorf("Bounding box Max X,Y = (%f, %f), want (2.0, 4.0)", box.Max.X(), box.Max.Y())
	}
}

func TestXZRectHit(t *testing.T) {
	mat := mockMaterial{}
	rect := XZRect{x0: -1.0, x1: 1.0, z0: -1.0, z1: 1.0, k: 0.0, material: mat}
	
	tests := []struct {
		name      string
		ray       Ray
		shouldHit bool
	}{
		{
			"Ray hits from above",
			NewRay(NewVec3d(0.0, 5.0, 0.0), NewVec3d(0.0, -1.0, 0.0), 0.0),
			true,
		},
		{
			"Ray hits from below",
			NewRay(NewVec3d(0.0, -5.0, 0.0), NewVec3d(0.0, 1.0, 0.0), 0.0),
			true,
		},
		{
			"Ray misses",
			NewRay(NewVec3d(5.0, 5.0, 0.0), NewVec3d(0.0, -1.0, 0.0), 0.0),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hit, rec := rect.Hit(tt.ray, 0.001, 10.0)
			if hit != tt.shouldHit {
				t.Errorf("Hit() = %v, want %v", hit, tt.shouldHit)
			}
			
			if hit {
				if rec.Normal.X() != 0.0 || math.Abs(rec.Normal.Y()) != 1.0 || rec.Normal.Z() != 0.0 {
					t.Errorf("XZRect normal should be (0,±1,0), got (%f,%f,%f)",
						rec.Normal.X(), rec.Normal.Y(), rec.Normal.Z())
				}
			}
		})
	}
}

func TestXZRectBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	rect := XZRect{x0: 1.0, x1: 2.0, z0: 3.0, z1: 4.0, k: 5.0, material: mat}
	
	var box AABB
	hasBox := rect.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("XZRect should have bounding box")
	}
	
	if box.Min.X() != 1.0 || box.Min.Z() != 3.0 {
		t.Errorf("Bounding box Min X,Z = (%f, %f), want (1.0, 3.0)", box.Min.X(), box.Min.Z())
	}
	
	if box.Max.X() != 2.0 || box.Max.Z() != 4.0 {
		t.Errorf("Bounding box Max X,Z = (%f, %f), want (2.0, 4.0)", box.Max.X(), box.Max.Z())
	}
}

func TestYZRectHit(t *testing.T) {
	mat := mockMaterial{}
	rect := YZRect{y0: -1.0, y1: 1.0, z0: -1.0, z1: 1.0, k: 0.0, material: mat}
	
	tests := []struct {
		name      string
		ray       Ray
		shouldHit bool
	}{
		{
			"Ray hits from right",
			NewRay(NewVec3d(5.0, 0.0, 0.0), NewVec3d(-1.0, 0.0, 0.0), 0.0),
			true,
		},
		{
			"Ray hits from left",
			NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			true,
		},
		{
			"Ray misses",
			NewRay(NewVec3d(5.0, 5.0, 5.0), NewVec3d(-1.0, 0.0, 0.0), 0.0),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hit, rec := rect.Hit(tt.ray, 0.001, 10.0)
			if hit != tt.shouldHit {
				t.Errorf("Hit() = %v, want %v", hit, tt.shouldHit)
			}
			
			if hit {
				if math.Abs(rec.Normal.X()) != 1.0 || rec.Normal.Y() != 0.0 || rec.Normal.Z() != 0.0 {
					t.Errorf("YZRect normal should be (±1,0,0), got (%f,%f,%f)",
						rec.Normal.X(), rec.Normal.Y(), rec.Normal.Z())
				}
			}
		})
	}
}

func TestYZRectBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	rect := YZRect{y0: 1.0, y1: 2.0, z0: 3.0, z1: 4.0, k: 5.0, material: mat}
	
	var box AABB
	hasBox := rect.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("YZRect should have bounding box")
	}
	
	if box.Min.Y() != 1.0 || box.Min.Z() != 3.0 {
		t.Errorf("Bounding box Min Y,Z = (%f, %f), want (1.0, 3.0)", box.Min.Y(), box.Min.Z())
	}
	
	if box.Max.Y() != 2.0 || box.Max.Z() != 4.0 {
		t.Errorf("Bounding box Max Y,Z = (%f, %f), want (2.0, 4.0)", box.Max.Y(), box.Max.Z())
	}
}

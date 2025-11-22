package rendim

import (
	"testing"
)

func TestConstantMediumHit(t *testing.T) {
	mat := mockMaterial{}
	boundary := NewSphere(NewVec3d(0.0, 0.0, 0.0), 2.0, mat)
	phaseFunction := Isotropic{albedo: ConstantTexture{color: Color{R: 1.0, G: 1.0, B: 1.0}}}
	
	cm := ConstantMedium{
		boundary:      boundary,
		density:       0.01,
		phaseFunction: phaseFunction,
	}
	
	ray := NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	
	hitCount := 0
	for i := 0; i < 100; i++ {
		hit, _ := cm.Hit(ray, 0.0, 10.0)
		if hit {
			hitCount++
		}
	}
	
	if hitCount == 0 {
		t.Error("ConstantMedium should sometimes scatter rays passing through it")
	}
	
	if hitCount == 100 {
		t.Error("ConstantMedium should sometimes let rays pass through")
	}
}

func TestConstantMediumNoHit(t *testing.T) {
	mat := mockMaterial{}
	boundary := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	phaseFunction := Isotropic{albedo: ConstantTexture{color: Color{R: 1.0, G: 1.0, B: 1.0}}}
	
	cm := ConstantMedium{
		boundary:      boundary,
		density:       0.01,
		phaseFunction: phaseFunction,
	}
	
	ray := NewRay(NewVec3d(5.0, 5.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	hit, _ := cm.Hit(ray, 0.0, 10.0)
	
	if hit {
		t.Error("ConstantMedium should not scatter rays that miss boundary")
	}
}

func TestConstantMediumBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	boundary := NewSphere(NewVec3d(1.0, 2.0, 3.0), 1.0, mat)
	phaseFunction := Isotropic{albedo: ConstantTexture{color: Color{R: 1.0, G: 1.0, B: 1.0}}}
	
	cm := ConstantMedium{
		boundary:      boundary,
		density:       0.01,
		phaseFunction: phaseFunction,
	}
	
	var box AABB
	hasBox := cm.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("ConstantMedium should have bounding box if boundary has one")
	}
	
	var boundaryBox AABB
	boundary.BoundingBox(0.0, 1.0, &boundaryBox)
	
	if box.Min.X() != boundaryBox.Min.X() || box.Min.Y() != boundaryBox.Min.Y() || box.Min.Z() != boundaryBox.Min.Z() {
		t.Error("ConstantMedium bounding box should match boundary bounding box")
	}
}

func TestConstantMediumHitWithinBounds(t *testing.T) {
	mat := mockMaterial{}
	boundary := NewSphere(NewVec3d(0.0, 0.0, 0.0), 2.0, mat)
	phaseFunction := Isotropic{albedo: ConstantTexture{color: Color{R: 1.0, G: 1.0, B: 1.0}}}
	
	cm := ConstantMedium{
		boundary:      boundary,
		density:       10.0,
		phaseFunction: phaseFunction,
	}
	
	ray := NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	
	for i := 0; i < 10; i++ {
		hit, rec := cm.Hit(ray, 0.0, 10.0)
		if hit && rec.material == nil {
			t.Error("Hit record should have material set")
		}
	}
}

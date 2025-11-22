package rendim

import (
	"math"
	"testing"
)

func TestHitableListHit(t *testing.T) {
	mat := mockMaterial{}
	sphere1 := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	sphere2 := NewSphere(NewVec3d(5.0, 0.0, 0.0), 1.0, mat)
	
	hl := HitableList{sphere1, sphere2}
	
	ray := NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	hit, rec := hl.Hit(ray, 0.0, 10.0)
	
	if !hit {
		t.Error("HitableList should detect hit on first sphere")
	}
	
	if math.Abs(rec.P.X()-(-1.0)) > 0.1 {
		t.Errorf("Hit point X = %f, expected around -1.0", rec.P.X())
	}
}

func TestHitableListNoHit(t *testing.T) {
	mat := mockMaterial{}
	sphere1 := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	sphere2 := NewSphere(NewVec3d(5.0, 0.0, 0.0), 1.0, mat)
	
	hl := HitableList{sphere1, sphere2}
	
	ray := NewRay(NewVec3d(-5.0, 10.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	hit, _ := hl.Hit(ray, 0.0, 100.0)
	
	if hit {
		t.Error("HitableList should not detect hit when ray misses all objects")
	}
}

func TestHitableListLen(t *testing.T) {
	mat := mockMaterial{}
	sphere1 := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	sphere2 := NewSphere(NewVec3d(5.0, 0.0, 0.0), 1.0, mat)
	
	hl := HitableList{sphere1, sphere2}
	
	if hl.Len() != 2 {
		t.Errorf("Len() = %d, want 2", hl.Len())
	}
}

func TestFlipNormalsHit(t *testing.T) {
	mat := mockMaterial{}
	sphere := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	flipped := FlipNormals{hitable: sphere}
	
	ray := NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	hitOriginal, recOriginal := sphere.Hit(ray, 0.0, 10.0)
	hitFlipped, recFlipped := flipped.Hit(ray, 0.0, 10.0)
	
	if !hitOriginal || !hitFlipped {
		t.Fatal("Both should hit")
	}
	
	if math.Abs(recOriginal.Normal.X()+recFlipped.Normal.X()) > 1e-10 ||
		math.Abs(recOriginal.Normal.Y()+recFlipped.Normal.Y()) > 1e-10 ||
		math.Abs(recOriginal.Normal.Z()+recFlipped.Normal.Z()) > 1e-10 {
		t.Error("Flipped normal should be opposite of original")
	}
}

func TestFlipNormalsBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	sphere := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	flipped := FlipNormals{hitable: sphere}
	
	var box AABB
	hasBox := flipped.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("FlipNormals should have bounding box if underlying object has one")
	}
}

func TestTranslateHit(t *testing.T) {
	mat := mockMaterial{}
	sphere := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	offset := NewVec3d(3.0, 0.0, 0.0)
	translated := Translate{hitable: sphere, offset: offset}
	
	ray := NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	hit, rec := translated.Hit(ray, 0.0, 10.0)
	
	if !hit {
		t.Error("Translated sphere should be hit")
	}
	
	if math.Abs(rec.P.X()-2.0) > 0.1 {
		t.Errorf("Hit point X = %f, expected around 2.0", rec.P.X())
	}
}

func TestTranslateBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	sphere := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	offset := NewVec3d(5.0, 0.0, 0.0)
	translated := Translate{hitable: sphere, offset: offset}
	
	var box AABB
	hasBox := translated.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("Translate should have bounding box if underlying object has one")
	}
	
	expectedMinX := 4.0
	if math.Abs(box.Min.X()-expectedMinX) > 1e-10 {
		t.Errorf("Bounding box Min.X = %f, want %f", box.Min.X(), expectedMinX)
	}
}

func TestNewRotateY(t *testing.T) {
	mat := mockMaterial{}
	sphere := NewSphere(NewVec3d(2.0, 0.0, 0.0), 1.0, mat)
	
	rotated := NewRotateY(sphere, 90.0)
	
	if math.Abs(rotated.sinTheta-1.0) > 1e-10 {
		t.Errorf("sin(90°) = %f, want 1.0", rotated.sinTheta)
	}
	
	if math.Abs(rotated.cosTheta) > 1e-10 {
		t.Errorf("cos(90°) = %f, want 0.0", rotated.cosTheta)
	}
}

func TestRotateYHit(t *testing.T) {
	mat := mockMaterial{}
	sphere := NewSphere(NewVec3d(2.0, 0.0, 0.0), 1.0, mat)
	rotated := NewRotateY(sphere, 90.0)
	
	ray := NewRay(NewVec3d(0.0, 0.0, -5.0), NewVec3d(0.0, 0.0, 1.0), 0.0)
	hit, _ := rotated.Hit(ray, 0.0, 10.0)
	
	if !hit {
		t.Error("Rotated sphere should be hit from the new position")
	}
}

func TestRotateYBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	sphere := NewSphere(NewVec3d(2.0, 0.0, 0.0), 1.0, mat)
	rotated := NewRotateY(sphere, 90.0)
	
	var box AABB
	hasBox := rotated.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("RotateY should have bounding box if underlying object has one")
	}
}

package rendim

import (
	"testing"
)

func TestNewBVHNodeSingle(t *testing.T) {
	mat := mockMaterial{}
	sphere := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	
	list := HitableList{sphere}
	bvh := NewBVHNode(list, 0.0, 1.0, NewRNG(0))
	
	var box AABB
	hasBox := bvh.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("BVHNode should have bounding box")
	}
}

func TestNewBVHNodeDouble(t *testing.T) {
	mat := mockMaterial{}
	sphere1 := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	sphere2 := NewSphere(NewVec3d(5.0, 0.0, 0.0), 1.0, mat)
	
	list := HitableList{sphere1, sphere2}
	bvh := NewBVHNode(list, 0.0, 1.0, NewRNG(0))
	
	var box AABB
	hasBox := bvh.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("BVHNode should have bounding box")
	}
	
	if box.Min.X() > -1.0 {
		t.Error("BVH bounding box should contain first sphere")
	}
	
	if box.Max.X() < 6.0 {
		t.Error("BVH bounding box should contain second sphere")
	}
}

func TestNewBVHNodeMultiple(t *testing.T) {
	mat := mockMaterial{}
	var list HitableList
	
	for i := 0; i < 10; i++ {
		sphere := NewSphere(NewVec3d(float64(i)*2.0, 0.0, 0.0), 0.5, mat)
		list = append(list, sphere)
	}
	
	bvh := NewBVHNode(list, 0.0, 1.0, NewRNG(0))
	
	var box AABB
	hasBox := bvh.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("BVHNode should have bounding box")
	}
}

func TestBVHNodeHit(t *testing.T) {
	mat := mockMaterial{}
	sphere1 := NewSphere(NewVec3d(0.0, 0.0, 0.0), 1.0, mat)
	sphere2 := NewSphere(NewVec3d(5.0, 0.0, 0.0), 1.0, mat)
	
	list := HitableList{sphere1, sphere2}
	bvh := NewBVHNode(list, 0.0, 1.0, NewRNG(0))
	
	tests := []struct {
		name      string
		ray       Ray
		shouldHit bool
	}{
		{
			"Ray hits first sphere",
			NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			true,
		},
		{
			"Ray hits second sphere",
			NewRay(NewVec3d(10.0, 0.0, 0.0), NewVec3d(-1.0, 0.0, 0.0), 0.0),
			true,
		},
		{
			"Ray misses all",
			NewRay(NewVec3d(0.0, 10.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0),
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hit, _ := bvh.Hit(tt.ray, 0.0, 100.0)
			if hit != tt.shouldHit {
				t.Errorf("Hit() = %v, want %v", hit, tt.shouldHit)
			}
		})
	}
}

func TestBVHNodeHitClosest(t *testing.T) {
	mat := mockMaterial{}
	sphere1 := NewSphere(NewVec3d(5.0, 0.0, 0.0), 1.0, mat)
	sphere2 := NewSphere(NewVec3d(10.0, 0.0, 0.0), 1.0, mat)
	
	list := HitableList{sphere1, sphere2}
	bvh := NewBVHNode(list, 0.0, 1.0, NewRNG(0))
	
	ray := NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	hit, rec := bvh.Hit(ray, 0.0, 100.0)
	
	if !hit {
		t.Error("Should hit one of the spheres")
	}
	
	if rec.P.X() > 5.0 {
		t.Error("Should hit closer sphere first")
	}
}

func TestBVHNodeBoundingBox(t *testing.T) {
	mat := mockMaterial{}
	sphere1 := NewSphere(NewVec3d(-5.0, -5.0, -5.0), 1.0, mat)
	sphere2 := NewSphere(NewVec3d(5.0, 5.0, 5.0), 1.0, mat)
	
	list := HitableList{sphere1, sphere2}
	bvh := NewBVHNode(list, 0.0, 1.0, NewRNG(0))
	
	var box AABB
	hasBox := bvh.BoundingBox(0.0, 1.0, &box)
	
	if !hasBox {
		t.Error("BVHNode should always have bounding box")
	}
	
	if box.Min.X() > -6.0 || box.Min.Y() > -6.0 || box.Min.Z() > -6.0 {
		t.Error("Bounding box should encompass all objects")
	}
	
	if box.Max.X() < 6.0 || box.Max.Y() < 6.0 || box.Max.Z() < 6.0 {
		t.Error("Bounding box should encompass all objects")
	}
}

func TestBVHNodeHitBothChildren(t *testing.T) {
	mat := mockMaterial{}
	sphere1 := NewSphere(NewVec3d(0.0, 0.0, 0.0), 2.0, mat)
	sphere2 := NewSphere(NewVec3d(1.0, 0.0, 0.0), 2.0, mat)
	
	list := HitableList{sphere1, sphere2}
	bvh := NewBVHNode(list, 0.0, 1.0, NewRNG(0))
	
	ray := NewRay(NewVec3d(-5.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	hit, _ := bvh.Hit(ray, 0.0, 100.0)
	
	if !hit {
		t.Error("Should hit when ray passes through overlapping spheres")
	}
}

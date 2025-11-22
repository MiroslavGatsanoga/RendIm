package rendim

import (
	"math"
	"testing"
)

func TestNewCamera(t *testing.T) {
	lookFrom := NewVec3d(0.0, 0.0, 5.0)
	lookAt := NewVec3d(0.0, 0.0, 0.0)
	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 90.0
	aspect := 2.0
	aperture := 0.0
	focusDist := 5.0
	
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspect, aperture, focusDist, 0.0, 1.0)
	
	if cam.origin.X() != 0.0 || cam.origin.Y() != 0.0 || cam.origin.Z() != 5.0 {
		t.Errorf("Camera origin = (%f, %f, %f), want (0.0, 0.0, 5.0)",
			cam.origin.X(), cam.origin.Y(), cam.origin.Z())
	}
	
	if cam.time0 != 0.0 || cam.time1 != 1.0 {
		t.Errorf("Camera time0=%f, time1=%f, want 0.0, 1.0", cam.time0, cam.time1)
	}
	
	if cam.lensRadius != 0.0 {
		t.Errorf("Camera lensRadius = %f, want 0.0", cam.lensRadius)
	}
}

func TestCameraGetRay(t *testing.T) {
	lookFrom := NewVec3d(0.0, 0.0, 0.0)
	lookAt := NewVec3d(0.0, 0.0, -1.0)
	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 90.0
	aspect := 2.0
	aperture := 0.0
	focusDist := 1.0
	
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspect, aperture, focusDist, 0.0, 1.0)
	
	ray := cam.GetRay(0.5, 0.5)
	
	if ray.Origin().X() != 0.0 || ray.Origin().Y() != 0.0 || ray.Origin().Z() != 0.0 {
		t.Errorf("Ray origin should be at camera position")
	}
	
	if ray.Time() < 0.0 || ray.Time() > 1.0 {
		t.Errorf("Ray time = %f, should be between 0.0 and 1.0", ray.Time())
	}
	
	if ray.Direction().Z() >= 0.0 {
		t.Error("Ray should point in negative Z direction")
	}
}

func TestCameraGetRayWithAperture(t *testing.T) {
	lookFrom := NewVec3d(0.0, 0.0, 0.0)
	lookAt := NewVec3d(0.0, 0.0, -1.0)
	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 90.0
	aspect := 2.0
	aperture := 2.0
	focusDist := 1.0
	
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspect, aperture, focusDist, 0.0, 1.0)
	
	ray := cam.GetRay(0.5, 0.5)
	
	if ray.Time() < 0.0 || ray.Time() > 1.0 {
		t.Errorf("Ray time = %f, should be between 0.0 and 1.0", ray.Time())
	}
}

func TestRandomInUnitDisk(t *testing.T) {
	for i := 0; i < 100; i++ {
		p := randomInUnitDisk()
		
		if p.Z() != 0.0 {
			t.Errorf("randomInUnitDisk Z component = %f, want 0.0", p.Z())
		}
		
		length := math.Sqrt(p.X()*p.X() + p.Y()*p.Y())
		if length >= 1.0 {
			t.Errorf("randomInUnitDisk returned vector with length %f >= 1.0", length)
		}
	}
}

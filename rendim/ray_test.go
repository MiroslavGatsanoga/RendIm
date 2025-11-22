package rendim

import (
	"math"
	"testing"
)

func TestNewRay(t *testing.T) {
	origin := NewVec3d(1.0, 2.0, 3.0)
	direction := NewVec3d(0.0, 1.0, 0.0)
	time := 0.5
	
	ray := NewRay(origin, direction, time)
	
	if ray.Origin().X() != 1.0 || ray.Origin().Y() != 2.0 || ray.Origin().Z() != 3.0 {
		t.Errorf("Origin() = (%f, %f, %f), want (1.0, 2.0, 3.0)",
			ray.Origin().X(), ray.Origin().Y(), ray.Origin().Z())
	}
	
	if ray.Direction().X() != 0.0 || ray.Direction().Y() != 1.0 || ray.Direction().Z() != 0.0 {
		t.Errorf("Direction() = (%f, %f, %f), want (0.0, 1.0, 0.0)",
			ray.Direction().X(), ray.Direction().Y(), ray.Direction().Z())
	}
	
	if ray.Time() != 0.5 {
		t.Errorf("Time() = %f, want 0.5", ray.Time())
	}
}

func TestRayPointAt(t *testing.T) {
	tests := []struct {
		name      string
		origin    Vec3d
		direction Vec3d
		t         float64
		expected  Vec3d
	}{
		{
			"At origin",
			NewVec3d(0.0, 0.0, 0.0),
			NewVec3d(1.0, 0.0, 0.0),
			0.0,
			NewVec3d(0.0, 0.0, 0.0),
		},
		{
			"Unit distance",
			NewVec3d(0.0, 0.0, 0.0),
			NewVec3d(1.0, 0.0, 0.0),
			1.0,
			NewVec3d(1.0, 0.0, 0.0),
		},
		{
			"Multiple units",
			NewVec3d(1.0, 2.0, 3.0),
			NewVec3d(2.0, 0.0, 0.0),
			3.0,
			NewVec3d(7.0, 2.0, 3.0),
		},
		{
			"Diagonal",
			NewVec3d(0.0, 0.0, 0.0),
			NewVec3d(1.0, 1.0, 1.0),
			2.0,
			NewVec3d(2.0, 2.0, 2.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ray := NewRay(tt.origin, tt.direction, 0.0)
			got := ray.PointAt(tt.t)
			
			if math.Abs(got.X()-tt.expected.X()) > 1e-10 ||
				math.Abs(got.Y()-tt.expected.Y()) > 1e-10 ||
				math.Abs(got.Z()-tt.expected.Z()) > 1e-10 {
				t.Errorf("PointAt(%f) = (%f, %f, %f), want (%f, %f, %f)",
					tt.t,
					got.X(), got.Y(), got.Z(),
					tt.expected.X(), tt.expected.Y(), tt.expected.Z())
			}
		})
	}
}

func TestRayAccessors(t *testing.T) {
	origin := NewVec3d(5.0, 6.0, 7.0)
	direction := NewVec3d(1.0, 2.0, 3.0)
	time := 1.5
	
	ray := NewRay(origin, direction, time)
	
	o := ray.Origin()
	if o.X() != 5.0 || o.Y() != 6.0 || o.Z() != 7.0 {
		t.Errorf("Origin() failed")
	}
	
	d := ray.Direction()
	if d.X() != 1.0 || d.Y() != 2.0 || d.Z() != 3.0 {
		t.Errorf("Direction() failed")
	}
	
	if ray.Time() != 1.5 {
		t.Errorf("Time() = %f, want 1.5", ray.Time())
	}
}

package rendim

import (
	"math"
	"testing"
)

func TestNewVec3d(t *testing.T) {
	v := NewVec3d(1.0, 2.0, 3.0)
	if v.X() != 1.0 || v.Y() != 2.0 || v.Z() != 3.0 {
		t.Errorf("NewVec3d failed: got (%f, %f, %f), want (1.0, 2.0, 3.0)", v.X(), v.Y(), v.Z())
	}
}

func TestVec3dAccessors(t *testing.T) {
	v := NewVec3d(4.5, 7.2, -3.1)
	if v.X() != 4.5 {
		t.Errorf("X() = %f, want 4.5", v.X())
	}
	if v.Y() != 7.2 {
		t.Errorf("Y() = %f, want 7.2", v.Y())
	}
	if v.Z() != -3.1 {
		t.Errorf("Z() = %f, want -3.1", v.Z())
	}
}

func TestVec3dLength(t *testing.T) {
	tests := []struct {
		name     string
		v        Vec3d
		expected float64
	}{
		{"Unit X", NewVec3d(1.0, 0.0, 0.0), 1.0},
		{"Unit Y", NewVec3d(0.0, 1.0, 0.0), 1.0},
		{"Unit Z", NewVec3d(0.0, 0.0, 1.0), 1.0},
		{"3-4-5 Triangle", NewVec3d(3.0, 4.0, 0.0), 5.0},
		{"Zero Vector", NewVec3d(0.0, 0.0, 0.0), 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.Length()
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("Length() = %f, want %f", got, tt.expected)
			}
		})
	}
}

func TestVec3dMultiplyScalar(t *testing.T) {
	v := NewVec3d(1.0, 2.0, 3.0)
	result := v.MultiplyScalar(2.0)
	if result.X() != 2.0 || result.Y() != 4.0 || result.Z() != 6.0 {
		t.Errorf("MultiplyScalar(2.0) = (%f, %f, %f), want (2.0, 4.0, 6.0)", result.X(), result.Y(), result.Z())
	}
}

func TestVec3dDivideScalar(t *testing.T) {
	v := NewVec3d(4.0, 8.0, 12.0)
	result := v.DivideScalar(4.0)
	if result.X() != 1.0 || result.Y() != 2.0 || result.Z() != 3.0 {
		t.Errorf("DivideScalar(4.0) = (%f, %f, %f), want (1.0, 2.0, 3.0)", result.X(), result.Y(), result.Z())
	}
}

func TestVec3dAdd(t *testing.T) {
	v1 := NewVec3d(1.0, 2.0, 3.0)
	v2 := NewVec3d(4.0, 5.0, 6.0)
	result := v1.Add(v2)
	if result.X() != 5.0 || result.Y() != 7.0 || result.Z() != 9.0 {
		t.Errorf("Add() = (%f, %f, %f), want (5.0, 7.0, 9.0)", result.X(), result.Y(), result.Z())
	}
}

func TestVec3dSubtract(t *testing.T) {
	v1 := NewVec3d(5.0, 7.0, 9.0)
	v2 := NewVec3d(1.0, 2.0, 3.0)
	result := v1.Subtract(v2)
	if result.X() != 4.0 || result.Y() != 5.0 || result.Z() != 6.0 {
		t.Errorf("Subtract() = (%f, %f, %f), want (4.0, 5.0, 6.0)", result.X(), result.Y(), result.Z())
	}
}

func TestVec3dDot(t *testing.T) {
	tests := []struct {
		name     string
		v1       Vec3d
		v2       Vec3d
		expected float64
	}{
		{"Perpendicular", NewVec3d(1.0, 0.0, 0.0), NewVec3d(0.0, 1.0, 0.0), 0.0},
		{"Parallel", NewVec3d(2.0, 0.0, 0.0), NewVec3d(3.0, 0.0, 0.0), 6.0},
		{"General", NewVec3d(1.0, 2.0, 3.0), NewVec3d(4.0, 5.0, 6.0), 32.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v1.Dot(tt.v2)
			if math.Abs(got-tt.expected) > 1e-10 {
				t.Errorf("Dot() = %f, want %f", got, tt.expected)
			}
		})
	}
}

func TestVec3dCross(t *testing.T) {
	tests := []struct {
		name     string
		v1       Vec3d
		v2       Vec3d
		expected Vec3d
	}{
		{"X cross Y = Z", NewVec3d(1.0, 0.0, 0.0), NewVec3d(0.0, 1.0, 0.0), NewVec3d(0.0, 0.0, 1.0)},
		{"Y cross Z = X", NewVec3d(0.0, 1.0, 0.0), NewVec3d(0.0, 0.0, 1.0), NewVec3d(1.0, 0.0, 0.0)},
		{"Z cross X = Y", NewVec3d(0.0, 0.0, 1.0), NewVec3d(1.0, 0.0, 0.0), NewVec3d(0.0, 1.0, 0.0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v1.Cross(tt.v2)
			if math.Abs(got.X()-tt.expected.X()) > 1e-10 ||
				math.Abs(got.Y()-tt.expected.Y()) > 1e-10 ||
				math.Abs(got.Z()-tt.expected.Z()) > 1e-10 {
				t.Errorf("Cross() = (%f, %f, %f), want (%f, %f, %f)",
					got.X(), got.Y(), got.Z(),
					tt.expected.X(), tt.expected.Y(), tt.expected.Z())
			}
		})
	}
}

func TestVec3dUnitVector(t *testing.T) {
	v := NewVec3d(3.0, 4.0, 0.0)
	unit := v.UnitVector()
	expectedLength := 1.0
	gotLength := unit.Length()
	
	if math.Abs(gotLength-expectedLength) > 1e-10 {
		t.Errorf("UnitVector().Length() = %f, want %f", gotLength, expectedLength)
	}
	
	if math.Abs(unit.X()-0.6) > 1e-10 || math.Abs(unit.Y()-0.8) > 1e-10 || math.Abs(unit.Z()-0.0) > 1e-10 {
		t.Errorf("UnitVector() = (%f, %f, %f), want (0.6, 0.8, 0.0)", unit.X(), unit.Y(), unit.Z())
	}
}

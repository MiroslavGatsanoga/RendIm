package rendim

import (
	"math"
	"testing"
)

type constantTexture struct {
	color Color
}

func (ct constantTexture) Value(u, v float64, p Vec3d) Color {
	return ct.color
}

func TestLambertianScatter(t *testing.T) {
	albedo := constantTexture{color: Color{R: 0.5, G: 0.5, B: 0.5}}
	mat := Lambertian{albedo: albedo}
	
	rayIn := NewRay(NewVec3d(-1.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	rec := HitRecord{
		P:      NewVec3d(0.0, 0.0, 0.0),
		Normal: NewVec3d(-1.0, 0.0, 0.0),
		t:      1.0,
	}
	
	var attenuation Color
	isScattered, _ := mat.Scatter(rayIn, rec, &attenuation)
	
	if !isScattered {
		t.Error("Lambertian should always scatter")
	}
	
	if attenuation.R != 0.5 || attenuation.G != 0.5 || attenuation.B != 0.5 {
		t.Errorf("Attenuation = (%f, %f, %f), want (0.5, 0.5, 0.5)",
			attenuation.R, attenuation.G, attenuation.B)
	}
}

func TestLambertianEmitted(t *testing.T) {
	albedo := constantTexture{color: Color{R: 0.5, G: 0.5, B: 0.5}}
	mat := Lambertian{albedo: albedo}
	
	emitted := mat.Emitted(0, 0, NewVec3d(0, 0, 0))
	
	if emitted.R != 0.0 || emitted.G != 0.0 || emitted.B != 0.0 {
		t.Error("Lambertian should not emit light")
	}
}

func TestIsotropicScatter(t *testing.T) {
	albedo := constantTexture{color: Color{R: 0.8, G: 0.7, B: 0.6}}
	mat := Isotropic{albedo: albedo}
	
	rayIn := NewRay(NewVec3d(-1.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	rec := HitRecord{
		P:      NewVec3d(0.0, 0.0, 0.0),
		Normal: NewVec3d(-1.0, 0.0, 0.0),
		t:      1.0,
	}
	
	var attenuation Color
	isScattered, _ := mat.Scatter(rayIn, rec, &attenuation)
	
	if !isScattered {
		t.Error("Isotropic should always scatter")
	}
	
	if attenuation.R != 0.8 || attenuation.G != 0.7 || attenuation.B != 0.6 {
		t.Errorf("Attenuation = (%f, %f, %f), want (0.8, 0.7, 0.6)",
			attenuation.R, attenuation.G, attenuation.B)
	}
}

func TestMetalScatter(t *testing.T) {
	albedo := constantTexture{color: Color{R: 0.9, G: 0.9, B: 0.9}}
	mat := Metal{albedo: albedo, fuzz: 0.0}
	
	rayIn := NewRay(NewVec3d(0.0, 0.0, 0.0), NewVec3d(1.0, -1.0, 0.0).UnitVector(), 0.0)
	rec := HitRecord{
		P:      NewVec3d(1.0, -1.0, 0.0),
		Normal: NewVec3d(0.0, 1.0, 0.0),
		t:      1.0,
	}
	
	var attenuation Color
	isScattered, scattered := mat.Scatter(rayIn, rec, &attenuation)
	
	if attenuation.R != 0.9 || attenuation.G != 0.9 || attenuation.B != 0.9 {
		t.Errorf("Attenuation = (%f, %f, %f), want (0.9, 0.9, 0.9)",
			attenuation.R, attenuation.G, attenuation.B)
	}
	
	if !isScattered {
		t.Error("Metal should scatter rays that reflect upward")
		return
	}
	
	if scattered.Direction().Y() <= 0.0 {
		t.Error("Metal should reflect upward when hitting surface from above")
	}
}

func TestMetalEmitted(t *testing.T) {
	albedo := constantTexture{color: Color{R: 0.9, G: 0.9, B: 0.9}}
	mat := Metal{albedo: albedo, fuzz: 0.1}
	
	emitted := mat.Emitted(0, 0, NewVec3d(0, 0, 0))
	
	if emitted.R != 0.0 || emitted.G != 0.0 || emitted.B != 0.0 {
		t.Error("Metal should not emit light")
	}
}

func TestDielectricScatter(t *testing.T) {
	mat := Dielectric{refIdx: 1.5}
	
	rayIn := NewRay(NewVec3d(-1.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	rec := HitRecord{
		P:      NewVec3d(0.0, 0.0, 0.0),
		Normal: NewVec3d(-1.0, 0.0, 0.0),
		t:      1.0,
	}
	
	var attenuation Color
	isScattered, _ := mat.Scatter(rayIn, rec, &attenuation)
	
	if !isScattered {
		t.Error("Dielectric should always scatter")
	}
	
	if attenuation.R != 1.0 || attenuation.G != 1.0 || attenuation.B != 1.0 {
		t.Errorf("Attenuation = (%f, %f, %f), want (1.0, 1.0, 1.0)",
			attenuation.R, attenuation.G, attenuation.B)
	}
}

func TestDiffuseLightScatter(t *testing.T) {
	emit := constantTexture{color: Color{R: 1.0, G: 1.0, B: 1.0}}
	mat := DiffuseLight{emit: emit}
	
	rayIn := NewRay(NewVec3d(-1.0, 0.0, 0.0), NewVec3d(1.0, 0.0, 0.0), 0.0)
	rec := HitRecord{
		P:      NewVec3d(0.0, 0.0, 0.0),
		Normal: NewVec3d(-1.0, 0.0, 0.0),
		t:      1.0,
	}
	
	var attenuation Color
	isScattered, _ := mat.Scatter(rayIn, rec, &attenuation)
	
	if isScattered {
		t.Error("DiffuseLight should not scatter")
	}
}

func TestDiffuseLightEmitted(t *testing.T) {
	emit := constantTexture{color: Color{R: 2.0, G: 3.0, B: 4.0}}
	mat := DiffuseLight{emit: emit}
	
	emitted := mat.Emitted(0.5, 0.5, NewVec3d(0, 0, 0))
	
	if emitted.R != 2.0 || emitted.G != 3.0 || emitted.B != 4.0 {
		t.Errorf("Emitted = (%f, %f, %f), want (2.0, 3.0, 4.0)",
			emitted.R, emitted.G, emitted.B)
	}
}

func TestReflect(t *testing.T) {
	v := NewVec3d(1.0, -1.0, 0.0)
	n := NewVec3d(0.0, 1.0, 0.0)
	
	reflected := reflect(v, n)
	
	expected := NewVec3d(1.0, 1.0, 0.0)
	
	if math.Abs(reflected.X()-expected.X()) > 1e-10 ||
		math.Abs(reflected.Y()-expected.Y()) > 1e-10 ||
		math.Abs(reflected.Z()-expected.Z()) > 1e-10 {
		t.Errorf("reflect() = (%f, %f, %f), want (%f, %f, %f)",
			reflected.X(), reflected.Y(), reflected.Z(),
			expected.X(), expected.Y(), expected.Z())
	}
}

func TestRefract(t *testing.T) {
	v := NewVec3d(1.0, -1.0, 0.0)
	n := NewVec3d(0.0, 1.0, 0.0)
	niOverNt := 1.0
	
	canRefract, refracted := refract(v, n, niOverNt)
	
	if !canRefract {
		t.Error("Should be able to refract with niOverNt = 1.0")
	}
	
	if refracted.Y() >= 0.0 {
		t.Error("Refracted ray should continue downward")
	}
}

func TestSchlick(t *testing.T) {
	tests := []struct {
		name    string
		cosine  float64
		refIdx  float64
		minProb float64
		maxProb float64
	}{
		{"Normal incidence", 1.0, 1.5, 0.03, 0.05},
		{"Grazing angle", 0.0, 1.5, 0.9, 1.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prob := schlick(tt.cosine, tt.refIdx)
			if prob < tt.minProb || prob > tt.maxProb {
				t.Errorf("schlick(%f, %f) = %f, want between %f and %f",
					tt.cosine, tt.refIdx, prob, tt.minProb, tt.maxProb)
			}
		})
	}
}

func TestRandomInUnitSphere(t *testing.T) {
	for i := 0; i < 100; i++ {
		p := randomInUnitSphere()
		length := p.Length()
		if length >= 1.0 {
			t.Errorf("randomInUnitSphere() returned vector with length %f >= 1.0", length)
		}
	}
}

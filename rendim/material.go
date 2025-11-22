package rendim

import (
	"math"
	"math/rand"
)

type Material interface {
	Scatter(rayIn Ray, rec HitRecord, attenuation *Color) (bool, Ray)
	Emitted(u, v float64, p Vec3d) Color
}

type Lambertian struct {
	albedo Texture
}

func (l Lambertian) Scatter(rayIn Ray, rec HitRecord, attenuation *Color) (isScattered bool, scattered Ray) {
	target := rec.P.Add(rec.Normal).Add(randomInUnitSphere())
	scattered = NewRay(rec.P, target.Subtract(rec.P), 0.0)
	*attenuation = l.albedo.Value(rec.u, rec.v, rec.P)
	return true, scattered
}

func (l Lambertian) Emitted(u, v float64, p Vec3d) Color {
	return Color{0, 0, 0}
}

type Isotropic struct {
	albedo Texture
}

func (i Isotropic) Scatter(rayIn Ray, rec HitRecord, attenuation *Color) (isScattered bool, scattered Ray) {
	scattered = NewRay(rec.P, randomInUnitSphere(), 0.0)
	*attenuation = i.albedo.Value(rec.u, rec.v, rec.P)
	return true, scattered
}

func (i Isotropic) Emitted(u, v float64, p Vec3d) Color {
	return Color{0, 0, 0}
}

type Metal struct {
	albedo Texture
	fuzz   float64
}

func (m Metal) Scatter(rayIn Ray, rec HitRecord, attenuation *Color) (isScattered bool, scattered Ray) {
	reflected := reflect(rayIn.Direction().UnitVector(), rec.Normal)
	scattered = NewRay(rec.P, reflected.Add(randomInUnitSphere().MultiplyScalar(m.fuzz)), 0.0)
	*attenuation = m.albedo.Value(0, 0, rec.P)
	return scattered.Direction().Dot(rec.Normal) > 0, scattered
}

func (m Metal) Emitted(u, v float64, p Vec3d) Color {
	return Color{0, 0, 0}
}

type Dielectric struct {
	refIdx float64
}

func (d Dielectric) Scatter(rayIn Ray, rec HitRecord, attenuation *Color) (isScattered bool, scattered Ray) {
	*attenuation = Color{R: 1.0, G: 1.0, B: 1.0}
	var (
		outwardNormal Vec3d
		NiOverNt      float64
		cosine        float64
	)

	rayInDotNormal := rayIn.Direction().Dot(rec.Normal)

	if rayInDotNormal > 0.0 {
		outwardNormal = rec.Normal.MultiplyScalar(-1.0)
		NiOverNt = d.refIdx
		cosine = d.refIdx * rayInDotNormal / rayIn.Direction().Length()
	} else {
		outwardNormal = rec.Normal
		NiOverNt = 1.0 / d.refIdx
		cosine = -rayInDotNormal / rayIn.Direction().Length()
	}

	var (
		reflectProb float64
		refracted   Vec3d
	)

	isRefracted, refracted := refract(rayIn.Direction(), outwardNormal, NiOverNt)

	if isRefracted {
		reflectProb = schlick(cosine, d.refIdx)
	} else {
		reflectProb = 1.0
	}

	if rand.Float64() < reflectProb { //nolint:gosec // G404: math/rand for reflection sampling
		reflected := reflect(rayIn.Direction(), rec.Normal)
		scattered = NewRay(rec.P, reflected, 0.0)
	} else {
		scattered = NewRay(rec.P, refracted, 0.0)
	}

	return true, scattered
}

func (d Dielectric) Emitted(u, v float64, p Vec3d) Color {
	return Color{0, 0, 0}
}

type DiffuseLight struct {
	emit Texture
}

func (dl DiffuseLight) Scatter(rayIn Ray, rec HitRecord, attenuation *Color) (isScattered bool, scattered Ray) {
	return false, Ray{}
}

func (dl DiffuseLight) Emitted(u, v float64, p Vec3d) Color {
	return dl.emit.Value(u, v, p)
}

func randomInUnitSphere() Vec3d {
	p := NewVec3d(rand.Float64(), rand.Float64(), rand.Float64()).MultiplyScalar(2.0).Subtract(NewVec3d(1.0, 1.0, 1.0)) //nolint:gosec // G404: math/rand for sampling
	for p.Dot(p) >= 1.0 {
		p = NewVec3d(rand.Float64(), rand.Float64(), rand.Float64()).MultiplyScalar(2.0).Subtract(NewVec3d(1.0, 1.0, 1.0)) //nolint:gosec // G404: math/rand for sampling
	}
	return p
}

func reflect(v, n Vec3d) Vec3d {
	tmp := n.MultiplyScalar(2.0 * v.Dot(n))
	return v.Subtract(tmp)
}

func refract(v, n Vec3d, NiOverNt float64) (isRefracted bool, refracted Vec3d) { //nolint:gocritic // captLocal: follows physics notation η_i/η_t
	uv := v.UnitVector()
	dt := uv.Dot(n)
	discriminant := 1.0 - NiOverNt*NiOverNt*(1.0-dt*dt)
	if discriminant > 0.0 {
		refracted = uv.
			Subtract(n.MultiplyScalar(dt)).
			MultiplyScalar(NiOverNt).
			Subtract(n.MultiplyScalar(math.Sqrt(discriminant)))

		return true, refracted
	}

	return false, Vec3d{}
}

func schlick(cosine, refIdx float64) float64 {
	r0 := (1.0 - refIdx) / (1.0 + refIdx)
	r0 *= r0
	return r0 + (1.0-r0)*math.Pow((1.0-cosine), 5.0)
}

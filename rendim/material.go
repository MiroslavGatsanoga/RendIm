package rendim

import (
	"math"
	"math/rand"
)

type Material interface {
	Scatter(rayIn Ray, rec HitRecord, attenuation *Color) (bool, Ray)
}

type Lambertian struct {
	albedo Color
}

func (l Lambertian) Scatter(rayIn Ray, rec HitRecord, attenuation *Color) (isScattered bool, scattered Ray) {
	target := rec.P.Add(rec.Normal).Add(randomInUnitSphere())
	scattered = NewRay(rec.P, target.Subtract(rec.P), 0.0)
	*attenuation = l.albedo
	return true, scattered
}

type Metal struct {
	albedo Color
	fuzz   float64
}

func (m Metal) Scatter(rayIn Ray, rec HitRecord, attenuation *Color) (isScattered bool, scattered Ray) {
	reflected := reflect(rayIn.Direction().UnitVector(), rec.Normal)
	scattered = NewRay(rec.P, reflected.Add(randomInUnitSphere().MultiplyScalar(m.fuzz)), 0.0)
	*attenuation = m.albedo
	return scattered.Direction().Dot(rec.Normal) > 0, scattered
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

	if rand.Float64() < reflectProb {
		reflected := reflect(rayIn.Direction(), rec.Normal)
		scattered = NewRay(rec.P, reflected, 0.0)
	} else {
		scattered = NewRay(rec.P, refracted, 0.0)
	}

	return true, scattered
}

func randomInUnitSphere() Vec3d {
	p := NewVec3d(rand.Float64(), rand.Float64(), rand.Float64()).MultiplyScalar(2.0).Subtract(NewVec3d(1.0, 1.0, 1.0))
	for p.Dot(p) >= 1.0 {
		p = NewVec3d(rand.Float64(), rand.Float64(), rand.Float64()).MultiplyScalar(2.0).Subtract(NewVec3d(1.0, 1.0, 1.0))
	}
	return p
}

func reflect(v, n Vec3d) Vec3d {
	tmp := n.MultiplyScalar(2.0 * v.Dot(n))
	return v.Subtract(tmp)
}

func refract(v, n Vec3d, NiOverNt float64) (isRefracted bool, refracted Vec3d) {
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
	r0 = r0 * r0
	return r0 + (1.0-r0)*math.Pow((1.0-cosine), 5.0)
}

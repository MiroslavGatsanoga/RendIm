package rendim

import "math"

type Material interface {
	Scatter(rayIn Ray, rec HitRecord, attenuation *Vec3d) (bool, Ray)
}

type Lambertian struct {
	albedo Vec3d
}

func (l Lambertian) Scatter(rayIn Ray, rec HitRecord, attenuation *Vec3d) (isScattered bool, scattered Ray) {
	target := rec.P.Add(rec.Normal).Add(randomInUnitSphere())
	scattered = NewRay(rec.P, target.Subtract(rec.P))
	*attenuation = l.albedo
	return true, scattered
}

type Metal struct {
	albedo Vec3d
	fuzz   float64
}

func (m Metal) Scatter(rayIn Ray, rec HitRecord, attenuation *Vec3d) (isScattered bool, scattered Ray) {
	reflected := reflect(rayIn.Direction().UnitVector(), rec.Normal)
	scattered = NewRay(rec.P, reflected.Add(randomInUnitSphere().MultiplyScalar(m.fuzz)))
	*attenuation = m.albedo
	return scattered.Direction().Dot(rec.Normal) > 0, scattered
}

type Dielectric struct {
	refIdx float64
}

func (d Dielectric) Scatter(rayIn Ray, rec HitRecord, attenuation *Vec3d) (isScattered bool, scattered Ray) {
	*attenuation = NewVec3d(1.0, 1.0, 1.0)
	var (
		outwardNormal Vec3d
		ni_over_nt    float64
	)

	rInNormalDot := rayIn.Direction().Dot(rec.Normal)

	if rInNormalDot > 0.0 {
		outwardNormal = rec.Normal.MultiplyScalar(-1.0)
		ni_over_nt = d.refIdx
	} else {
		outwardNormal = rec.Normal
		ni_over_nt = 1.0 / d.refIdx
	}

	if isRefracted, refracted := refract(rayIn.Direction(), outwardNormal, ni_over_nt); isRefracted {
		scattered = NewRay(rec.P, refracted)

	} else {
		reflected := reflect(rayIn.Direction(), rec.Normal)
		scattered = NewRay(rec.P, reflected)
	}

	return true, scattered
}

func randomInUnitSphere() Vec3d {
	p := NewVec3d(rnd.Float64(), rnd.Float64(), rnd.Float64()).MultiplyScalar(2.0).Subtract(NewVec3d(1.0, 1.0, 1.0))
	for p.Dot(p) >= 1.0 {
		p = NewVec3d(rnd.Float64(), rnd.Float64(), rnd.Float64()).MultiplyScalar(2.0).Subtract(NewVec3d(1.0, 1.0, 1.0))
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

package rendim

type Material interface {
	Scatter(rayIn Ray, rec HitRecord, attenuation *Vec3d) (bool, Ray)
}

type Lambertian struct {
	albedo Vec3d
}

func (l Lambertian) Scatter(rayIn Ray, rec HitRecord, attenuation *Vec3d) (isScattered bool, scattered Ray) {
	target := rec.P.Add(rec.Normal).Add(RandomInUnitSphere())
	scattered = NewRay(rec.P, target.Subtract(rec.P))
	*attenuation = l.albedo
	return true, scattered
}

type Metal struct {
	albedo Vec3d
	fuzz   float64
}

func (m Metal) Scatter(rayIn Ray, rec HitRecord, attenuation *Vec3d) (isScattered bool, scattered Ray) {
	reflected := Reflect(rayIn.Direction().UnitVector(), rec.Normal)
	scattered = NewRay(rec.P, reflected.Add(RandomInUnitSphere().MultiplyScalar(m.fuzz)))
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

	if isRefracted, refracted := Refract(rayIn.Direction(), outwardNormal, ni_over_nt); isRefracted {
		scattered = NewRay(rec.P, refracted)

	} else {
		reflected := Reflect(rayIn.Direction(), rec.Normal) //todo: direction unit vector ?
		scattered = NewRay(rec.P, reflected)
	}

	return true, scattered
}

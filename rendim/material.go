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

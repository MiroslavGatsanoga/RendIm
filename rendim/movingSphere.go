package rendim

import "math"

type MovingSphere struct {
	center0, center1 Vec3d
	time0, time1     float64
	Radius           float64
	material         Material
}

func NewMovingSphere(cen0, cen1 Vec3d, t0, t1, radius float64, material Material) MovingSphere {
	return MovingSphere{center0: cen0, center1: cen1, time0: t0, time1: t1, Radius: radius, material: material}
}

func (s MovingSphere) Center(time float64) Vec3d {
	return s.center0.Add(s.center1.Subtract(s.center0).MultiplyScalar((time - s.time0) / (s.time1 - s.time0)))
}

func (s MovingSphere) Hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	oc := r.Origin().Subtract(s.Center(r.Time()))
	a := r.Direction().Dot(r.Direction())
	b := oc.Dot(r.Direction())
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - a*c

	if discriminant > 0.0 {
		rec.material = s.material
		temp := (-b - math.Sqrt(discriminant)) / a
		if temp > tMin && temp < tMax {
			rec.t = temp
			rec.P = r.PointAt(rec.t)
			rec.Normal = rec.P.Subtract(s.Center(r.Time())).DivideScalar(s.Radius)
			return true
		}
		temp = (-b + math.Sqrt(discriminant)) / a
		if temp > tMin && temp < tMax {
			rec.t = temp
			rec.P = r.PointAt(rec.t)
			rec.Normal = rec.P.Subtract(s.Center(r.Time())).DivideScalar(s.Radius)
			return true
		}
	}

	return false
}

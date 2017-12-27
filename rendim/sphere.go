package rendim

import "math"

type Sphere struct {
	Center Vec3d
	Radius float64
}

func NewSphere(center Vec3d, radius float64) Sphere {
	return Sphere{Center: center, Radius: radius}
}

func (s Sphere) Hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	oc := r.Origin().Subtract(s.Center)
	a := r.Direction().Dot(r.Direction())
	b := oc.Dot(r.Direction())
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - a*c

	if discriminant > 0.0 {
		temp := (-b - math.Sqrt(discriminant)) / a
		if temp > tMin && temp < tMax {
			rec.t = temp
			rec.P = r.PointAt(rec.t)
			rec.Normal = rec.P.Subtract(s.Center).DivideScalar(s.Radius)
			return true
		}
		temp = (-b + math.Sqrt(discriminant)) / a
		if temp > tMin && temp < tMax {
			rec.t = temp
			rec.P = r.PointAt(rec.t)
			rec.Normal = rec.P.Subtract(s.Center).DivideScalar(s.Radius)
			return true
		}
	}

	return false
}

package rendim

import (
	"math"
)

type Sphere struct {
	Center   Vec3d
	Radius   float64
	material Material
}

func NewSphere(center Vec3d, radius float64, material Material) Sphere {
	return Sphere{Center: center, Radius: radius, material: material}
}

func (s Sphere) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	oc := r.Origin().Subtract(s.Center)
	a := r.Direction().Dot(r.Direction())
	b := oc.Dot(r.Direction())
	c := oc.Dot(oc) - s.Radius*s.Radius
	discriminant := b*b - a*c

	rec := HitRecord{}
	if discriminant > 0.0 {
		rec.material = s.material
		temp := (-b - math.Sqrt(discriminant)) / a
		if temp > tMin && temp < tMax {
			rec.t = temp
			rec.P = r.PointAt(rec.t)
			u, v := getSphereUV((rec.P.Subtract(s.Center)).DivideScalar(s.Radius))
			rec.u = u
			rec.v = v
			rec.Normal = rec.P.Subtract(s.Center).DivideScalar(s.Radius)
			return true, rec
		}
		temp = (-b + math.Sqrt(discriminant)) / a
		if temp > tMin && temp < tMax {
			rec.t = temp
			rec.P = r.PointAt(rec.t)
			u, v := getSphereUV((rec.P.Subtract(s.Center)).DivideScalar(s.Radius))
			rec.u = u
			rec.v = v
			rec.Normal = rec.P.Subtract(s.Center).DivideScalar(s.Radius)
			return true, rec
		}
	}

	return false, rec
}

func (s Sphere) BoundingBox(t0, t1 float64, box *AABB) bool {
	boxMin := s.Center.Subtract(NewVec3d(s.Radius, s.Radius, s.Radius))
	boxMax := s.Center.Add(NewVec3d(s.Radius, s.Radius, s.Radius))
	*box = AABB{Min: boxMin, Max: boxMax}
	return true
}

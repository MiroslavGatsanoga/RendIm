package rendim

type Ray struct {
	a, b Vec3d
}

func NewRay(origin, direction Vec3d) Ray {
	return Ray{a: origin, b: direction}
}

func (r Ray) Origin() Vec3d {
	return r.a
}

func (r Ray) Direction() Vec3d {
	return r.b
}

func (r Ray) PointAt(t float64) Vec3d {
	return r.a.Add(r.b.MultiplyScalar(t))
}

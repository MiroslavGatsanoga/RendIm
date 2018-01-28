package rendim

type Ray struct {
	a, b Vec3d
	time float64
}

func NewRay(origin, direction Vec3d, ti float64) Ray {
	return Ray{a: origin, b: direction, time: ti}
}

func (r Ray) Origin() Vec3d {
	return r.a
}

func (r Ray) Direction() Vec3d {
	return r.b
}

func (r Ray) Time() float64 {
	return r.time
}

func (r Ray) PointAt(t float64) Vec3d {
	return r.a.Add(r.b.MultiplyScalar(t))
}

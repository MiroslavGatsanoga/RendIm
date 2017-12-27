package rendim

type Camera struct {
	origin          Vec3d
	lowerLeftCorner Vec3d
	horizontal      Vec3d
	vertical        Vec3d
}

func NewCamera() Camera {
	c := Camera{}
	c.lowerLeftCorner = NewVec3d(-2.0, -1.0, -1.0)
	c.horizontal = NewVec3d(4.0, 0.0, 0.0)
	c.vertical = NewVec3d(0.0, 2.0, 0.0)
	c.origin = NewVec3d(0.0, 0.0, 0.0)
	return c
}

func (c Camera) GetRay(u, v float64) Ray {
	rayDirection := c.lowerLeftCorner.Add(c.horizontal.MultiplyScalar(u).Add(c.vertical.MultiplyScalar(v)))
	return NewRay(c.origin, rayDirection)
}

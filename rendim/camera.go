package rendim

import (
	"math"
)

type Camera struct {
	origin          Vec3d
	lowerLeftCorner Vec3d
	horizontal      Vec3d
	vertical        Vec3d
	u, v, w         Vec3d
}

func NewCamera(lookFrom, lookAt, vUp Vec3d, vFov, aspect float64) Camera {
	theta := vFov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight

	c := Camera{}
	c.origin = lookFrom
	c.w = lookFrom.Subtract(lookAt).UnitVector()
	c.u = vUp.Cross(c.w).UnitVector()
	c.v = c.w.Cross(c.u)
	c.lowerLeftCorner = c.origin.Subtract(c.u.MultiplyScalar(halfWidth)).Subtract(c.v.MultiplyScalar(halfHeight)).Subtract(c.w)
	c.horizontal = c.u.MultiplyScalar(2.0 * halfWidth)
	c.vertical = c.v.MultiplyScalar(2.0 * halfHeight)
	return c
}

func (c Camera) GetRay(s, t float64) Ray {
	rayDirection := c.lowerLeftCorner.Add(c.horizontal.MultiplyScalar(s).Add(c.vertical.MultiplyScalar(t)).Subtract(c.origin))
	return NewRay(c.origin, rayDirection)
}

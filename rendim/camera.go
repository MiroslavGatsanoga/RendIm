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
	lensRadius      float64
	time0, time1    float64
}

func NewCamera(lookFrom, lookAt, vUp Vec3d, vFov, aspect, aperture, focusDist, t0, t1 float64) Camera {
	theta := vFov * math.Pi / 180.0
	halfHeight := math.Tan(theta / 2.0)
	halfWidth := aspect * halfHeight

	c := Camera{}
	c.time0 = t0
	c.time1 = t1
	c.lensRadius = aperture / 2.0
	c.origin = lookFrom
	c.w = lookFrom.Subtract(lookAt).UnitVector()
	c.u = vUp.Cross(c.w).UnitVector()
	c.v = c.w.Cross(c.u)
	c.lowerLeftCorner = c.origin.Subtract(c.u.MultiplyScalar(halfWidth * focusDist)).Subtract(c.v.MultiplyScalar(halfHeight * focusDist)).Subtract(c.w.MultiplyScalar(focusDist))
	c.horizontal = c.u.MultiplyScalar(2.0 * halfWidth * focusDist)
	c.vertical = c.v.MultiplyScalar(2.0 * halfHeight * focusDist)
	return c
}

func (c Camera) GetRay(s, t float64, rng *RNG) Ray {
	rd := randomInUnitDisk(rng).MultiplyScalar(c.lensRadius)
	offset := c.u.MultiplyScalar(rd.X()).Add(c.v.MultiplyScalar(rd.Y()))
	time := c.time0 + rng.Float64()*(c.time1-c.time0)
	rayDirection := c.lowerLeftCorner.Add(c.horizontal.MultiplyScalar(s)).Add(c.vertical.MultiplyScalar(t)).Subtract(c.origin)
	return NewRay(c.origin.Add(offset), rayDirection.Subtract(offset), time)
}

func randomInUnitDisk(rng *RNG) Vec3d {
	p := NewVec3d(rng.Float64(), rng.Float64(), 0.0).MultiplyScalar(2.0).Subtract(NewVec3d(1.0, 1.0, 0.0))
	for p.Dot(p) >= 1.0 {
		p = NewVec3d(rng.Float64(), rng.Float64(), 0.0).MultiplyScalar(2.0).Subtract(NewVec3d(1.0, 1.0, 0.0))
	}
	return p
}

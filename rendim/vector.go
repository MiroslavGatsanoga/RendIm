package rendim

import (
	"math"
)

type Vec3d struct {
	e [3]float64
}

func NewVec3d(x, y, z float64) Vec3d {
	return Vec3d{e: [3]float64{x, y, z}}
}

func (v Vec3d) X() float64 {
	return v.e[0]
}

func (v Vec3d) Y() float64 {
	return v.e[1]
}

func (v Vec3d) Z() float64 {
	return v.e[2]
}

func (v Vec3d) Length() float64 {
	return math.Sqrt(v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2])
}

func (v Vec3d) LengthSquared() float64 {
	return v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2]
}

func (v Vec3d) MultiplyScalar(s float64) Vec3d {
	return Vec3d{e: [3]float64{
		v.e[0] * s,
		v.e[1] * s,
		v.e[2] * s,
	}}
}

func (v Vec3d) DivideScalar(s float64) Vec3d {
	return Vec3d{e: [3]float64{
		v.e[0] / s,
		v.e[1] / s,
		v.e[2] / s,
	}}
}

func (v Vec3d) Add(v2 Vec3d) Vec3d {
	return Vec3d{e: [3]float64{
		v.e[0] + v2.e[0],
		v.e[1] + v2.e[1],
		v.e[2] + v2.e[2],
	}}
}

func (v Vec3d) Subtract(v2 Vec3d) Vec3d {
	return Vec3d{e: [3]float64{
		v.e[0] - v2.e[0],
		v.e[1] - v2.e[1],
		v.e[2] - v2.e[2],
	}}
}

func (v Vec3d) Dot(v2 Vec3d) float64 {
	return v.e[0]*v2.e[0] + v.e[1]*v2.e[1] + v.e[2]*v2.e[2]
}

func (v Vec3d) Cross(v2 Vec3d) Vec3d {
	return Vec3d{e: [3]float64{
		v.e[1]*v2.e[2] - v.e[2]*v2.e[1],
		-(v.e[0]*v2.e[2] - v.e[2]*v2.e[0]),
		v.e[0]*v2.e[1] - v.e[1]*v2.e[0],
	}}
}

func (v Vec3d) UnitVector() Vec3d {
	return v.DivideScalar(v.Length())
}

func RandomInUnitSphere() Vec3d {
	p := NewVec3d(rnd.Float64(), rnd.Float64(), rnd.Float64()).MultiplyScalar(2.0).Subtract(NewVec3d(1.0, 1.0, 1.0))
	for p.Dot(p) >= 1.0 {
		p = NewVec3d(rnd.Float64(), rnd.Float64(), rnd.Float64()).MultiplyScalar(2.0).Subtract(NewVec3d(1.0, 1.0, 1.0))
	}
	return p
}

func Reflect(v, n Vec3d) Vec3d {
	tmp := n.MultiplyScalar(2.0 * v.Dot(n))
	return v.Subtract(tmp)
}

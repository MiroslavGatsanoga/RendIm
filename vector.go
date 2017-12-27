package main

import (
	"math"
)

type Vec3d struct {
	e [3]float64
}

func NewVec3d(x, y, z float64) Vec3d {
	return Vec3d{e: [3]float64{x, y, z}}
}

func (v Vec3d) Length() float64 {
	return math.Sqrt(v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2])
}

func (v Vec3d) LengthSquared() float64 {
	return v.e[0]*v.e[0] + v.e[1]*v.e[1] + v.e[2]*v.e[2]
}

func (v Vec3d) Multiply(s float64) Vec3d {
	return Vec3d{e: [3]float64{
		v.e[0] * s,
		v.e[1] * s,
		v.e[2] * s,
	}}
}

func (v Vec3d) Divide(s float64) Vec3d {
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
	return v.Divide(v.Length())
}

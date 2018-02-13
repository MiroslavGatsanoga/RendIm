package rendim

import "math"

type Hitable interface {
	Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord)
	BoundingBox(t0, t1 float64, box *AABB) bool
}

type HitRecord struct {
	u, v     float64
	t        float64
	P        Vec3d
	Normal   Vec3d
	material Material
}

type HitableList []Hitable

func (hl HitableList) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	hitAnything := false
	closestSoFar := tMax
	rec := HitRecord{}
	for _, h := range hl {
		if isHit, hr := h.Hit(r, tMin, closestSoFar); isHit {
			hitAnything = true
			closestSoFar = hr.t

			rec.t = hr.t
			rec.u = hr.u
			rec.v = hr.v
			rec.P = hr.P
			rec.Normal = hr.Normal
			rec.material = hr.material
		}
	}

	return hitAnything, rec
}

func (hl HitableList) Len() int {
	return len(hl)
}

type FlipNormals struct {
	hitable Hitable
}

func (f FlipNormals) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	if isHit, rec := f.hitable.Hit(r, tMin, tMax); isHit {
		rec.Normal = rec.Normal.MultiplyScalar(-1.0)
		return true, rec
	}
	return false, HitRecord{}
}

func (f FlipNormals) BoundingBox(t0, t1 float64, box *AABB) bool {
	return f.hitable.BoundingBox(t0, t1, box)
}

type Translate struct {
	hitable Hitable
	offset  Vec3d
}

func (t Translate) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	movedRay := NewRay(r.Origin().Subtract(t.offset), r.Direction(), r.Time())
	if isHit, rec := t.hitable.Hit(movedRay, tMin, tMax); isHit {
		rec.P = rec.P.Add(t.offset)
		return true, rec
	}
	return false, HitRecord{}
}

func (t Translate) BoundingBox(t0, t1 float64, box *AABB) bool {
	if t.hitable.BoundingBox(t0, t1, box) {
		*box = AABB{Min: (*box).Min.Add(t.offset), Max: (*box).Max.Add(t.offset)}
		return true
	}
	return false
}

type RotateY struct {
	hitable            Hitable
	sinTheta, cosTheta float64
	hasBox             bool
	bbox               AABB
}

func NewRotateY(h Hitable, angle float64) RotateY {
	ry := RotateY{hitable: h}
	radians := (math.Pi / 180.0) * angle
	ry.sinTheta = math.Sin(radians)
	ry.cosTheta = math.Cos(radians)
	ry.hasBox = h.BoundingBox(0.0, 1.0, &ry.bbox)

	min := NewVec3d(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64)
	max := NewVec3d(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64)
	ijk := [2]float64{0.0, 1.0}
	for _, i := range ijk {
		for _, j := range ijk {
			for _, k := range ijk {
				x := i*ry.bbox.Max.X() + (1.0-i)*ry.bbox.Min.X()
				y := j*ry.bbox.Max.Y() + (1.0-j)*ry.bbox.Min.Y()
				z := k*ry.bbox.Max.Z() + (1.0-k)*ry.bbox.Min.Z()
				newX := ry.cosTheta*x + ry.sinTheta*z
				newZ := -ry.sinTheta*x + ry.cosTheta*z
				tester := NewVec3d(newX, y, newZ)
				for c := 0; c < 3; c++ {
					if tester.e[c] > max.e[c] {
						max.e[c] = tester.e[c]
					}
					if tester.e[c] < min.e[c] {
						min.e[c] = tester.e[c]
					}
				}
			}
		}
	}
	ry.bbox = AABB{Min: min, Max: max}
	return ry
}

func (ry RotateY) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	origin := r.Origin()
	direction := r.Direction()

	origin.e[0] = ry.cosTheta*r.Origin().e[0] - ry.sinTheta*r.Origin().e[2]
	origin.e[2] = ry.sinTheta*r.Origin().e[0] + ry.cosTheta*r.Origin().e[2]
	direction.e[0] = ry.cosTheta*r.Direction().e[0] - ry.sinTheta*r.Direction().e[2]
	direction.e[2] = ry.sinTheta*r.Direction().e[0] + ry.cosTheta*r.Direction().e[2]

	rotatedRay := NewRay(origin, direction, r.Time())
	if isHit, rec := ry.hitable.Hit(rotatedRay, tMin, tMax); isHit {
		p := rec.P
		normal := rec.Normal
		p.e[0] = ry.cosTheta*rec.P.e[0] + ry.sinTheta*rec.P.e[2]
		p.e[2] = -ry.sinTheta*rec.P.e[0] + ry.cosTheta*rec.P.e[2]
		normal.e[0] = ry.cosTheta*rec.Normal.e[0] + ry.sinTheta*rec.Normal.e[2]
		normal.e[2] = -ry.sinTheta*rec.Normal.e[0] + ry.cosTheta*rec.Normal.e[2]
		rec.P = p
		rec.Normal = normal
		return true, rec
	}
	return false, HitRecord{}
}

func (ry RotateY) BoundingBox(t0, t1 float64, box *AABB) bool {
	*box = ry.bbox
	return ry.hasBox
}

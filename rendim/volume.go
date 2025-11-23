package rendim

import (
	"math"
)

type ConstantMedium struct {
	boundary      Hitable
	density       float64
	phaseFunction Material
	rng           *RNG
}

func (cm ConstantMedium) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	if isHit1, rec1 := cm.boundary.Hit(r, -math.MaxFloat64, math.MaxFloat64); isHit1 {
		if isHit2, rec2 := cm.boundary.Hit(r, rec1.t+0.0001, math.MaxFloat64); isHit2 {
			if rec1.t < tMin {
				rec1.t = tMin
			}
			if rec2.t > tMax {
				rec2.t = tMax
			}
			if rec1.t >= rec2.t {
				return false, HitRecord{}
			}
			if rec1.t < 0 {
				rec1.t = 0
			}

			rec := HitRecord{}
			distanceInsideBoundary := (rec2.t - rec1.t) * r.Direction().Length()
			hitDistance := -(1.0 / cm.density) * math.Log(cm.rng.Float64())
			if hitDistance < distanceInsideBoundary {
				rec.t = rec1.t + hitDistance/r.Direction().Length()
				rec.P = r.PointAt(rec.t)
				rec.Normal = NewVec3d(1.0, 0.0, 0.0)
				rec.material = cm.phaseFunction
				return true, rec
			}
		}
	}
	return false, HitRecord{}
}

func (cm ConstantMedium) BoundingBox(t0, t1 float64, box *AABB) bool {
	return cm.boundary.BoundingBox(t0, t1, box)
}

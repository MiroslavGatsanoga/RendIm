package rendim

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

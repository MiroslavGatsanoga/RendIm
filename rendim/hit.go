package rendim

type Hitable interface {
	Hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool
	BoundingBox(t0, t1 float64, box *AABB) bool
}

type HitRecord struct {
	t        float64
	P        Vec3d
	Normal   Vec3d
	material Material
}

type HitableList []Hitable

func (hl HitableList) Hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	tempRec := &HitRecord{}
	hitAnything := false
	closestSoFar := tMax
	for _, h := range hl {
		if h.Hit(r, tMin, closestSoFar, tempRec) {
			hitAnything = true
			closestSoFar = tempRec.t

			//rec = tempRec //todo:NOTE!!! copy tempRec to rec
			rec.t = tempRec.t
			rec.P = tempRec.P
			rec.Normal = tempRec.Normal
			rec.material = tempRec.material
		}
	}
	return hitAnything
}

func (hl HitableList) Len() int {
	return len(hl)
}

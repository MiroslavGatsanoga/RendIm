package rendim

type Hitable interface {
	Hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool
}

type HitRecord struct {
	t      float64
	p      Vec3d
	Normal Vec3d
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
			rec.p = tempRec.p
			rec.Normal = tempRec.Normal
		}
	}
	return hitAnything
}

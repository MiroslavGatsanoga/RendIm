package rendim

type AABB struct {
	Min, Max Vec3d
}

func ffMin(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func ffMax(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func (box AABB) hit(r Ray, tMin float64, tMax float64) bool {
	for a := 0; a < 3; a++ {
		invD := 1.0 / r.Direction().e[a]
		t0 := (box.Min.e[a] - r.Origin().e[a]) * invD
		t1 := (box.Max.e[a] - r.Origin().e[a]) * invD

		if invD < 0.0 {
			t0, t1 = t1, t0
		}

		tMin = ffMax(t0, tMin)
		tMax = ffMin(t1, tMax)
		if tMax <= tMin {
			return false
		}
	}
	return true
}

func surroundingBox(box0, box1 AABB) AABB {
	small := NewVec3d(
		ffMin(box0.Min.X(), box1.Min.X()),
		ffMin(box0.Min.Y(), box1.Min.Y()),
		ffMin(box0.Min.Z(), box1.Min.Z()),
	)
	big := NewVec3d(
		ffMax(box0.Max.X(), box1.Max.X()),
		ffMax(box0.Max.Y(), box1.Max.Y()),
		ffMax(box0.Max.Z(), box1.Max.Z()),
	)
	return AABB{small, big}
}

package rendim

type XYRect struct {
	x0, x1, y0, y1, k float64
	material          Material
}

func (rect XYRect) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	t := (rect.k - r.Origin().Z()) / r.Direction().Z()
	if t < tMin || t > tMax {
		return false, HitRecord{}
	}

	x := r.Origin().X() + t*r.Direction().X()
	y := r.Origin().Y() + t*r.Direction().Y()
	if x < rect.x0 || x > rect.x1 || y < rect.y0 || y > rect.y1 {
		return false, HitRecord{}
	}

	rec := HitRecord{}
	rec.u = (x - rect.x0) / (rect.x1 - rect.x0)
	rec.v = (y - rect.y0) / (rect.y1 - rect.y0)
	rec.t = t
	rec.material = rect.material
	rec.P = r.PointAt(t)
	rec.Normal = NewVec3d(0.0, 0.0, 1.0)
	return true, rec
}

func (rect XYRect) BoundingBox(t0, t1 float64, box *AABB) bool {
	boxMin := NewVec3d(rect.x0, rect.y0, rect.k-0.0001)
	boxMax := NewVec3d(rect.x1, rect.y1, rect.k+0.0001)
	*box = AABB{Min: boxMin, Max: boxMax}
	return true
}

type XZRect struct {
	x0, x1, z0, z1, k float64
	material          Material
}

func (rect XZRect) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	t := (rect.k - r.Origin().Y()) / r.Direction().Y()
	if t < tMin || t > tMax {
		return false, HitRecord{}
	}

	x := r.Origin().X() + t*r.Direction().X()
	z := r.Origin().Z() + t*r.Direction().Z()
	if x < rect.x0 || x > rect.x1 || z < rect.z0 || z > rect.z1 {
		return false, HitRecord{}
	}

	rec := HitRecord{}
	rec.u = (x - rect.x0) / (rect.x1 - rect.x0)
	rec.v = (z - rect.z0) / (rect.z1 - rect.z0)
	rec.t = t
	rec.material = rect.material
	rec.P = r.PointAt(t)
	rec.Normal = NewVec3d(0.0, 1.0, 0.0)
	return true, rec
}

func (rect XZRect) BoundingBox(t0, t1 float64, box *AABB) bool {
	boxMin := NewVec3d(rect.x0, rect.k-0.0001, rect.z0)
	boxMax := NewVec3d(rect.x1, rect.k+0.0001, rect.z1)
	*box = AABB{Min: boxMin, Max: boxMax}
	return true
}

type YZRect struct {
	y0, y1, z0, z1, k float64
	material          Material
}

func (rect YZRect) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	t := (rect.k - r.Origin().X()) / r.Direction().X()
	if t < tMin || t > tMax {
		return false, HitRecord{}
	}

	y := r.Origin().Y() + t*r.Direction().Y()
	z := r.Origin().Z() + t*r.Direction().Z()
	if y < rect.y0 || y > rect.y1 || z < rect.z0 || z > rect.z1 {
		return false, HitRecord{}
	}

	rec := HitRecord{}
	rec.u = (y - rect.y0) / (rect.y1 - rect.y0)
	rec.v = (z - rect.z0) / (rect.z1 - rect.z0)
	rec.t = t
	rec.material = rect.material
	rec.P = r.PointAt(t)
	rec.Normal = NewVec3d(1.0, 0.0, 0.0)
	return true, rec
}

func (rect YZRect) BoundingBox(t0, t1 float64, box *AABB) bool {
	boxMin := NewVec3d(rect.k-0.0001, rect.y0, rect.z0)
	boxMax := NewVec3d(rect.k+0.0001, rect.y1, rect.z1)
	*box = AABB{Min: boxMin, Max: boxMax}
	return true
}

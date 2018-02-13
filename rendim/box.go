package rendim

type Box struct {
	pMin, pMax Vec3d
	faces      HitableList
}

func NewBox(p0, p1 Vec3d, mat Material) Box {
	box := Box{}
	box.pMin = p0
	box.pMax = p1

	box.faces = append(box.faces, XYRect{x0: p0.X(), x1: p1.X(), y0: p0.Y(), y1: p1.Y(), k: p1.Z(), material: mat})
	box.faces = append(box.faces, FlipNormals{hitable: XYRect{x0: p0.X(), x1: p1.X(), y0: p0.Y(), y1: p1.Y(), k: p0.Z(), material: mat}})
	box.faces = append(box.faces, XZRect{x0: p0.X(), x1: p1.X(), z0: p0.Z(), z1: p1.Z(), k: p1.Y(), material: mat})
	box.faces = append(box.faces, FlipNormals{hitable: XZRect{x0: p0.X(), x1: p1.X(), z0: p0.Z(), z1: p1.Z(), k: p0.Y(), material: mat}})
	box.faces = append(box.faces, YZRect{y0: p0.Y(), y1: p1.Y(), z0: p0.Z(), z1: p1.Z(), k: p1.X(), material: mat})
	box.faces = append(box.faces, FlipNormals{hitable: YZRect{y0: p0.Y(), y1: p1.Y(), z0: p0.Z(), z1: p1.Z(), k: p0.X(), material: mat}})

	return box
}

func (b Box) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	return b.faces.Hit(r, tMin, tMax)
}

func (b Box) BoundingBox(t0, t1 float64, box *AABB) bool {
	*box = AABB{Min: b.pMin, Max: b.pMax}
	return true
}

package rendim

import (
	"math/rand"
	"sort"
)

type BVHNode struct {
	left, right *Hitable
	box         AABB
}

func NewBVHNode(l HitableList, time0, time1 float64) Hitable {
	bvh := BVHNode{}

	axis := rand.Intn(3) //nolint:gosec // G404: math/rand appropriate for Monte Carlo sampling in graphics

	if axis == 0 { //nolint:gocritic,staticcheck // ifElseChain/QF1003: if-else more readable than switch for sort logic
		// sort by x
		sort.Slice(l,
			func(i, j int) bool {
				boxLeft, boxRight := &AABB{}, &AABB{}
				if !l[i].BoundingBox(time0, time1, boxLeft) ||
					!l[j].BoundingBox(time0, time1, boxRight) {
					panic("No bounding box in BVHNode constructor.\n")
				}

				return boxLeft.Min.X() < boxRight.Min.X()
			})
	} else if axis == 1 { // sort by y
		sort.Slice(l,
			func(i, j int) bool {
				boxLeft, boxRight := &AABB{}, &AABB{}
				if !l[i].BoundingBox(time0, time1, boxLeft) ||
					!l[j].BoundingBox(time0, time1, boxRight) {
					panic("No bounding box in BVHNode constructor.\n")
				}

				return boxLeft.Min.Y() < boxRight.Min.Y()
			})
	} else { // sort by z
		sort.Slice(l,
			func(i, j int) bool {
				boxLeft, boxRight := &AABB{}, &AABB{}
				if !l[i].BoundingBox(time0, time1, boxLeft) ||
					!l[j].BoundingBox(time0, time1, boxRight) {
					panic("No bounding box in BVHNode constructor.\n")
				}

				return boxLeft.Min.Z() < boxRight.Min.Z()
			})
	}

	n := l.Len()
	if n == 1 { //nolint:gocritic,staticcheck // ifElseChain/QF1003: if-else clearer than switch for recursive partitioning
		bvh.left = &l[0]
		bvh.right = &l[0]
	} else if n == 2 {
		bvh.left = &l[0]
		bvh.right = &l[1]
	} else {
		newLeft := NewBVHNode(l[:n/2], time0, time1)
		bvh.left = &newLeft
		newRight := NewBVHNode(l[n/2:], time0, time1)
		bvh.right = &newRight
	}

	boxLeft, boxRight := &AABB{}, &AABB{}
	if !(*bvh.left).BoundingBox(time0, time1, boxLeft) ||
		!(*bvh.right).BoundingBox(time0, time1, boxRight) {
		panic("No bounding box in BVHNode constructor.\n")
	}

	bvh.box = surroundingBox(*boxLeft, *boxRight)

	return bvh
}

func (n BVHNode) BoundingBox(t0, t1 float64, box *AABB) bool {
	*box = n.box
	return true
}

func (n BVHNode) Hit(r Ray, tMin float64, tMax float64) (bool, HitRecord) {
	if n.box.hit(r, tMin, tMax) {
		hitLeft, leftRec := (*n.left).Hit(r, tMin, tMax)
		hitRight, rightRec := (*n.right).Hit(r, tMin, tMax)

		rec := HitRecord{}
		if hitLeft && hitRight { //nolint:gocritic // ifElseChain: boolean conditions better as if-else than switch
			if leftRec.t < rightRec.t {
				rec.t = leftRec.t
				rec.u = leftRec.u
				rec.v = leftRec.v
				rec.P = leftRec.P
				rec.Normal = leftRec.Normal
				rec.material = leftRec.material
			} else {
				rec.t = rightRec.t
				rec.u = rightRec.u
				rec.v = rightRec.v
				rec.P = rightRec.P
				rec.Normal = rightRec.Normal
				rec.material = rightRec.material
			}
			return true, rec
		} else if hitLeft {
			rec.t = leftRec.t
			rec.u = leftRec.u
			rec.v = leftRec.v
			rec.P = leftRec.P
			rec.Normal = leftRec.Normal
			rec.material = leftRec.material
			return true, rec
		} else if hitRight {
			rec.t = rightRec.t
			rec.u = rightRec.u
			rec.v = rightRec.v
			rec.P = rightRec.P
			rec.Normal = rightRec.Normal
			rec.material = rightRec.material
			return true, rec
		}
		return false, rec
	}
	return false, HitRecord{}
}

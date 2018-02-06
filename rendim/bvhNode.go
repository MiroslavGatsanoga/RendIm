package rendim

import (
	"fmt"
	"math/rand"
	"sort"
)

type BVHNode struct {
	left, right *Hitable
	box         AABB
}

func NewBVHNode(l HitableList, time0, time1 float64) Hitable {
	bvh := BVHNode{}

	axis := rand.Intn(3)

	if axis == 0 { // sort by x
		sort.Slice(l,
			func(i, j int) bool {
				boxLeft, boxRight := &AABB{}, &AABB{}
				if !l[i].BoundingBox(time0, time1, boxLeft) ||
					!l[j].BoundingBox(time0, time1, boxRight) {
					panic(fmt.Sprintf("No bounding box in BVHNode constructor.\n"))
				}

				return (*boxLeft).Min.X() < (*boxRight).Min.X()
			})
	} else if axis == 1 { // sort by y
		sort.Slice(l,
			func(i, j int) bool {
				boxLeft, boxRight := &AABB{}, &AABB{}
				if !l[i].BoundingBox(time0, time1, boxLeft) ||
					!l[j].BoundingBox(time0, time1, boxRight) {
					panic(fmt.Sprintf("No bounding box in BVHNode constructor.\n"))
				}

				return (*boxLeft).Min.Y() < (*boxRight).Min.Y()
			})
	} else { // sort by z
		sort.Slice(l,
			func(i, j int) bool {
				boxLeft, boxRight := &AABB{}, &AABB{}
				if !l[i].BoundingBox(time0, time1, boxLeft) ||
					!l[j].BoundingBox(time0, time1, boxRight) {
					panic(fmt.Sprintf("No bounding box in BVHNode constructor.\n"))
				}

				return (*boxLeft).Min.Z() < (*boxRight).Min.Z()
			})
	}

	n := l.Len()
	if n == 1 {
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
		panic(fmt.Sprintf("No bounding box in BVHNode constructor.\n"))
	}

	bvh.box = surroundingBox(*boxLeft, *boxRight)

	return bvh
}

func (n BVHNode) BoundingBox(t0, t1 float64, box *AABB) bool {
	*box = n.box
	return true
}

func (n BVHNode) Hit(r Ray, tMin float64, tMax float64, rec *HitRecord) bool {
	if n.box.hit(r, tMin, tMax) {
		leftRec, rightRec := &HitRecord{}, &HitRecord{}
		hitLeft := (*n.left).Hit(r, tMin, tMax, leftRec)
		hitRight := (*n.right).Hit(r, tMin, tMax, rightRec)

		if hitLeft && hitRight {
			if leftRec.t < rightRec.t {
				rec.t = leftRec.t
				rec.P = leftRec.P
				rec.Normal = leftRec.Normal
				rec.material = leftRec.material
			} else {
				rec.t = rightRec.t
				rec.P = rightRec.P
				rec.Normal = rightRec.Normal
				rec.material = rightRec.material
			}
			return true
		} else if hitLeft {
			rec.t = leftRec.t
			rec.P = leftRec.P
			rec.Normal = leftRec.Normal
			rec.material = leftRec.material
			return true
		} else if hitRight {
			rec.t = rightRec.t
			rec.P = rightRec.P
			rec.Normal = rightRec.Normal
			rec.material = rightRec.material
			return true
		}
		return false
	}
	return false
}

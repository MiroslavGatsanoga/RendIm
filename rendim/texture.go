package rendim

import (
	"math"
)

type Texture interface {
	Value(u, v float64, p Vec3d) Color
}

type ConstantTexture struct {
	color Color
}

func (t ConstantTexture) Value(u, v float64, p Vec3d) Color {
	return t.color
}

type CheckerTexture struct {
	even, odd Texture
}

func (t CheckerTexture) Value(u, v float64, p Vec3d) Color {
	sines := math.Sin(10*p.X()) * math.Sin(10*p.Y()) * math.Sin(10*p.Z())
	if sines < 0.0 {
		return t.odd.Value(u, v, p)
	}

	return t.even.Value(u, v, p)
}

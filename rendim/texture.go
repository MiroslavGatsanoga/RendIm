package rendim

import (
	"math"
	"math/rand"
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

type NoiseTexture struct {
	noise perlin
	scale float64
}

func (t NoiseTexture) Value(u, v float64, p Vec3d) Color {
	clr := Color{R: 1.0, G: 1.0, B: 1.0}
	t.noise = perlinNoise
	return clr.MultiplyScalar(0.5 * (1.0 + math.Sin(t.scale*p.Z()+10.0*t.noise.Turbulence(p))))
}

type perlin struct {
	ranVec              []Vec3d
	permX, permY, permZ []int
}

func (pn perlin) Noise(p Vec3d) float64 {
	u := p.X() - math.Floor(p.X())
	v := p.Y() - math.Floor(p.Y())
	w := p.Z() - math.Floor(p.Z())
	i := int(math.Floor(p.X()))
	j := int(math.Floor(p.Y()))
	k := int(math.Floor(p.Z()))

	var c [2][2][2]Vec3d
	for di := 0; di < 2; di++ {
		for dj := 0; dj < 2; dj++ {
			for dk := 0; dk < 2; dk++ {
				c[di][dj][dk] = pn.ranVec[pn.permX[(i+di)&255]^pn.permY[(j+dj)&255]^pn.permZ[(k+dk)&255]]
			}
		}
	}

	return perlinInterp(c, u, v, w)
}
func (pn perlin) Turbulence(p Vec3d) float64 {
	depth := 7
	var accum float64
	tempP := p
	weigth := 1.0
	for i := 0; i < depth; i++ {
		accum += weigth * pn.Noise(tempP)
		weigth *= 0.5
		tempP = tempP.MultiplyScalar(2.0)
	}
	return math.Abs(accum)
}

func perlinInterp(c [2][2][2]Vec3d, u, v, w float64) float64 {
	uu := u * u * (3.0 - 2.0*u)
	vv := v * v * (3.0 - 2.0*v)
	ww := w * w * (3.0 - 2.0*w)

	var accum float64
	ijk := [2]float64{0.0, 1.0}
	for _, i := range ijk {
		for _, j := range ijk {
			for _, k := range ijk {
				weightVec := NewVec3d(u-i, v-j, w-k)
				accum += (i*uu + (1.0-i)*(1.0-uu)) *
					(j*vv + (1.0-j)*(1.0-vv)) *
					(k*ww + (1.0-k)*(1.0-ww)) *
					c[int(i)][int(j)][int(k)].Dot(weightVec)
			}
		}
	}
	return accum
}

var perlinNoise = perlin{
	ranVec: perlinGenerate(),
	permX:  perlinGeneratePerm(),
	permY:  perlinGeneratePerm(),
	permZ:  perlinGeneratePerm(),
}

func perlinGenerate() []Vec3d {
	p := make([]Vec3d, 256)
	for i := range p {
		p[i] = NewVec3d(
			-1.0+2.0*rand.Float64(),
			-1.0+2.0*rand.Float64(),
			-1.0+2.0*rand.Float64(),
		)
	}
	return p
}

func permute(p []int) {
	n := len(p)
	for i := n - 1; i > 0; i-- {
		target := int(rand.Float64() * float64(i+1))
		p[i], p[target] = p[target], p[i]
	}
}

func perlinGeneratePerm() []int {
	p := make([]int, 256)
	for i := range p {
		p[i] = i
	}

	permute(p)
	return p
}

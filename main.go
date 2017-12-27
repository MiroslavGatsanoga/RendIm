package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"path-tracer/rendim"
	"time"
)

const (
	width   = 400
	height  = 200
	samples = 100
)

var rnd = rand.New(rand.NewSource(42))

func main() {
	start := time.Now()
	img := render()
	elapsed := time.Since(start)
	fmt.Println("Image rendering took:", elapsed)

	f, err := os.Create("out.png")
	defer f.Close()
	if err != nil {
		panic("cannot create out.png")
	}

	png.Encode(f, img)
}

func render() image.Image {

	//scene setup
	world := rendim.HitableList{}
	world = append(world, rendim.NewSphere(rendim.NewVec3d(0.0, 0.0, -1.0), 0.5))
	world = append(world, rendim.NewSphere(rendim.NewVec3d(0.0, -100.5, -1.0), 100.0))

	//camera setup
	cam := rendim.NewCamera()

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			var rayClr rendim.Vec3d
			for s := 0; s < samples; s++ {
				u := (float64(px) + rnd.Float64()) / float64(width)
				v := (float64(height-py) + rnd.Float64()) / float64(height)
				r := cam.GetRay(u, v)
				rayClr = rayClr.Add(rayColor(r, &world))
			}
			rayClr = rayClr.DivideScalar(samples)
			rayClrGamma := rendim.NewVec3d(
				math.Sqrt(rayClr.X()),
				math.Sqrt(rayClr.Y()),
				math.Sqrt(rayClr.Z()))

			ir := uint8(255.99 * rayClrGamma.X())
			ig := uint8(255.99 * rayClrGamma.Y())
			ib := uint8(255.99 * rayClrGamma.Z())

			clr := color.RGBA{R: ir, G: ig, B: ib, A: 255}

			img.Set(px, py, clr)
		}
	}

	return img
}

func rayColor(r rendim.Ray, world *rendim.HitableList) rendim.Vec3d {
	rec := &rendim.HitRecord{}
	if world.Hit(r, 0.0, math.MaxFloat64, rec) {
		target := rec.P.Add(rec.Normal).Add(randomInUnitSphere())
		newRay := rendim.NewRay(rec.P, target.Subtract(rec.P))
		return rayColor(newRay, world).MultiplyScalar(0.5)
	}

	unitDirection := r.Direction().UnitVector()
	t := 0.5 * (unitDirection.Y() + 1.0)
	white := rendim.NewVec3d(1.0, 1.0, 1.0)
	blue := rendim.NewVec3d(0.5, 0.7, 1.0)
	clr := white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
	return clr
}

func randomInUnitSphere() rendim.Vec3d {
	p := rendim.NewVec3d(rnd.Float64(), rnd.Float64(), rnd.Float64()).MultiplyScalar(2.0).Subtract(rendim.NewVec3d(1.0, 1.0, 1.0))
	for p.Dot(p) >= 1.0 {
		p = rendim.NewVec3d(rnd.Float64(), rnd.Float64(), rnd.Float64()).MultiplyScalar(2.0).Subtract(rendim.NewVec3d(1.0, 1.0, 1.0))
	}
	return p
}

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path-tracer/rendim"
	"time"
)

const (
	width, height = 800, 400
)

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
	//camera setup
	lowerLeftCorner := rendim.NewVec3d(-2.0, -1.0, -1.0)
	horizontal := rendim.NewVec3d(4.0, 0.0, 0.0)
	vertical := rendim.NewVec3d(0.0, 2.0, 0.0)
	origin := rendim.NewVec3d(0.0, 0.0, 0.0)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			u := float64(px) / width
			v := float64(height-py) / height

			rayDirection := lowerLeftCorner.Add(horizontal.MultiplyScalar(u).Add(vertical.MultiplyScalar(v)))
			r := rendim.NewRay(origin, rayDirection)

			rayClr := rayColor(r)

			ir := uint8(255.99 * rayClr.X())
			ig := uint8(255.99 * rayClr.Y())
			ib := uint8(255.99 * rayClr.Z())

			clr := color.RGBA{R: ir, G: ig, B: ib, A: 255}

			img.Set(px, py, clr)
		}
	}

	return img
}

func rayColor(r rendim.Ray) rendim.Vec3d {
	sphereCenter := rendim.NewVec3d(0.0, 0.0, -1.0)
	t := hitSphere(sphereCenter, 0.5, r)
	if t > 0.0 {
		N := r.PointAt(t).Subtract(sphereCenter).UnitVector()
		return rendim.NewVec3d(N.X()+1.0, N.Y()+1, N.Z()+1).MultiplyScalar(0.5)
	}

	unitDirection := r.Direction().UnitVector()
	t = 0.5 * (unitDirection.Y() + 1.0)
	white := rendim.NewVec3d(1.0, 1.0, 1.0)
	blue := rendim.NewVec3d(0.5, 0.7, 1.0)
	clr := white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
	return clr
}

func hitSphere(center rendim.Vec3d, radius float64, r rendim.Ray) float64 {
	oc := r.Origin().Subtract(center)
	a := r.Direction().Dot(r.Direction())
	b := 2.0 * oc.Dot(r.Direction())
	c := oc.Dot(oc) - radius*radius
	discriminant := b*b - 4*a*c

	if discriminant < 0.0 {
		return -1.0
	}

	return (-b - math.Sqrt(discriminant)) / (2.0 * a)
}

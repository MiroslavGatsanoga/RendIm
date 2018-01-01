package rendim

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"
)

const (
	width   = 400
	height  = 200
	samples = 100
)

var rnd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func Render() image.Image {

	//scene setup
	world := HitableList{}
	world = append(world, NewSphere(NewVec3d(0.0, 0.0, -1.0), 0.5, Lambertian{albedo: NewVec3d(0.1, 0.2, 0.5)}))
	world = append(world, NewSphere(NewVec3d(0.0, -100.5, -1.0), 100.0, Lambertian{albedo: NewVec3d(0.8, 0.8, 0.0)}))
	world = append(world, NewSphere(NewVec3d(1.0, 0.0, -1.0), 0.5, Metal{albedo: NewVec3d(0.8, 0.6, 0.2), fuzz: 0.0}))
	world = append(world, NewSphere(NewVec3d(-1.0, 0.0, -1.0), 0.5, Dielectric{refIdx: 1.5}))
	world = append(world, NewSphere(NewVec3d(-1.0, 0.0, -1.0), -0.45, Dielectric{refIdx: 1.5}))

	//camera setup
	lookFrom := NewVec3d(-2.0, 2.0, 1.0)
	lookAt := NewVec3d(0.0, 0.0, -1.0)
	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 20.0 //vertical field of view in degrees
	aspectRatio := float64(width) / float64(height)
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspectRatio)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py := 0; py < height; py++ {
		for px := 0; px < width; px++ {
			var rayClr Vec3d
			for s := 0; s < samples; s++ {
				u := (float64(px) + rnd.Float64()) / float64(width)
				v := (float64(height-py) + rnd.Float64()) / float64(height)
				r := cam.GetRay(u, v)
				rayClr = rayClr.Add(rayColor(r, &world, 0))
			}
			rayClr = rayClr.DivideScalar(samples)
			rayClrGamma := NewVec3d(
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

func rayColor(r Ray, world *HitableList, depth int) Vec3d {
	rec := &HitRecord{}
	if world.Hit(r, 0.001, math.MaxFloat64, rec) {
		attenuation := &Vec3d{}
		if depth < 50 {
			isScattered, scattered := rec.material.Scatter(r, *rec, attenuation)
			if isScattered {
				clr := rayColor(scattered, world, depth+1)
				x := attenuation.X() * clr.X()
				y := attenuation.Y() * clr.Y()
				z := attenuation.Z() * clr.Z()

				return NewVec3d(x, y, z)
			}
		}

		return NewVec3d(0.0, 0.0, 0.0)
	}

	unitDirection := r.Direction().UnitVector()
	t := 0.5 * (unitDirection.Y() + 1.0)
	white := NewVec3d(1.0, 1.0, 1.0)
	blue := NewVec3d(0.5, 0.7, 1.0)
	clr := white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
	return clr
}

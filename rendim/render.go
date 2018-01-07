package rendim

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"strings"
	"time"
)

const (
	width   = 1200
	height  = 800
	samples = 50
)

var rnd = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

func Render() image.Image {
	//scene setup
	world := randomScene()

	//camera setup
	lookFrom := NewVec3d(13.0, 2.0, 3.0)
	lookAt := NewVec3d(0.0, 0.0, 0.0)
	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 20.0 //vertical field of view in degrees
	aspectRatio := float64(width) / float64(height)
	distToFocus := 10.0
	aperture := 0.1
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspectRatio, aperture, distToFocus)

	var (
		py, px, s int
	)

	done := make(chan bool)
	go showProgress(&py, &px, &s, done)

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for py = 0; py < height; py++ {
		for px = 0; px < width; px++ {
			var rayClr Color
			for s = 0; s < samples; s++ {
				u := (float64(px) + rnd.Float64()) / float64(width)
				v := (float64(height-py) + rnd.Float64()) / float64(height)
				r := cam.GetRay(u, v)
				rayClr = rayClr.Add(rayColor(r, &world, 0))
			}
			rayClr = rayClr.DivideScalar(samples)
			rayClrGamma := Color{
				R: math.Sqrt(rayClr.R),
				G: math.Sqrt(rayClr.G),
				B: math.Sqrt(rayClr.B)}

			clr := rayClrGamma.ToRGBA()

			img.Set(px, py, clr)
		}
	}

	done <- true
	<-done

	return img
}

func rayColor(r Ray, world *HitableList, depth int) Color {
	rec := &HitRecord{}
	if world.Hit(r, 0.001, math.MaxFloat64, rec) {
		attenuation := &Color{}
		if depth < 50 {
			isScattered, scattered := rec.material.Scatter(r, *rec, attenuation)
			if isScattered {
				clr := rayColor(scattered, world, depth+1)
				return attenuation.Multiply(clr)
			}
		}

		return Color{}
	}

	unitDirection := r.Direction().UnitVector()
	t := 0.5 * (unitDirection.Y() + 1.0)
	white := Color{R: 1.0, G: 1.0, B: 1.0}
	blue := Color{R: 0.5, G: 0.7, B: 1.0}
	clr := white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
	return clr
}

func randomScene() HitableList {
	list := HitableList{}
	list = append(list, NewSphere(NewVec3d(0.0, -1000.0, 0), 1000, Lambertian{albedo: Color{R: 0.5, G: 0.5, B: 0.5}}))
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMaterial := rnd.Float64()
			center := NewVec3d(float64(a)+0.9*rnd.Float64(), 0.2, float64(b)+0.9*rnd.Float64())
			if center.Subtract(NewVec3d(4.0, 0.2, 0.0)).Length() > 0.9 {
				if chooseMaterial < 0.8 { //diffuse
					list = append(list, NewSphere(center, 0.2,
						Lambertian{albedo: Color{R: rnd.Float64() * rnd.Float64(), G: rnd.Float64() * rnd.Float64(), B: rnd.Float64() * rnd.Float64()}}))
				} else if chooseMaterial < 0.95 { //metal
					list = append(list, NewSphere(center, 0.2,
						Metal{albedo: Color{R: 0.5 * (1.0 + rnd.Float64()), G: 0.5 * (1.0 + rnd.Float64()), B: 0.5 * (1.0 + rnd.Float64())}, fuzz: 0.5 * rnd.Float64()}))
				} else { //glass
					list = append(list, NewSphere(center, 0.2, Dielectric{refIdx: 1.5}))
				}
			}
		}
	}

	list = append(list, NewSphere(NewVec3d(0.0, 1.0, 0.0), 1.0, Dielectric{refIdx: 1.5}))
	list = append(list, NewSphere(NewVec3d(-4.0, 1.0, 0.0), 1.0, Lambertian{albedo: Color{R: 0.4, G: 0.2, B: 0.1}}))
	list = append(list, NewSphere(NewVec3d(4.0, 1.0, 0.0), 1.0, Metal{albedo: Color{R: 0.7, G: 0.6, B: 0.5}, fuzz: 0.0}))

	return list
}

func showProgress(py, px, s *int, done chan bool) {
	const tickIntervalMs = 200
	ticker := time.NewTicker(time.Millisecond * tickIntervalMs)
	elapsed := 0
	for {
		select {
		case <-ticker.C:
			cy, cx, cs := (*py), (*px), (*s)
			opNum := (cy * width * samples) + (cx * samples) + cs
			progress := float64(opNum) / float64(width*height*samples)
			progressPercent := int(100.0 * progress)

			var progressBar bytes.Buffer
			for i := 0; i < progressPercent/2; i++ {
				progressBar.WriteString("=")
			}
			progressBar.WriteString(">")

			elapsed += tickIntervalMs
			elapsedDuration := time.Second * time.Duration(elapsed/1000)
			fmt.Printf("\r[%-50s] %d %% %v", progressBar.String(), progressPercent, elapsedDuration)
		case <-done:
			fmt.Printf("\r[%50s] Done                    \n", strings.Repeat("=", 50))
			done <- true
			return
		}
	}
}

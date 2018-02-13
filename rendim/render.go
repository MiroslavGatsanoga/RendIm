package rendim

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"math"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	samples      = 100
	bucketSize   = 32
	workersCount = 4
)

var ops uint64

type Pixel struct {
	image.Point
	R, G, B uint8
}

func Render(width, height int, pixels chan Pixel) image.Image {
	scene := cornellBox(width, height)

	return renderBuckets(width, height, scene, pixels)
}

func renderBuckets(width, height int, scene Scene, pixels chan Pixel) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	buckets := getBuckets(img.Bounds())
	bucketChan := make(chan image.Rectangle, len(buckets))

	done := make(chan bool)
	go showProgress(width*height, done)

	var wg sync.WaitGroup
	wg.Add(workersCount)

	for w := 0; w < workersCount; w++ {
		go renderBucket(bucketChan, &scene, img, &wg, pixels)
	}

	for _, b := range buckets {
		bucketChan <- b
	}

	close(bucketChan)
	wg.Wait()

	done <- true
	<-done

	return img
}

func getBuckets(r image.Rectangle) []image.Rectangle {
	w := r.Bounds().Dx()
	h := r.Bounds().Dy()

	bw := w/bucketSize + 1
	bh := h/bucketSize + 1

	buckets := make([]image.Rectangle, 0)

	for y := 0; y < bh; y++ {
		if y%2 == 0 {
			for x := 0; x < bw; x++ {
				b := image.Rectangle{
					Min: image.Point{X: x * bucketSize, Y: y * bucketSize},
					Max: image.Point{X: (x+1)*bucketSize - 1, Y: (y+1)*bucketSize - 1},
				}

				buckets = append(buckets, b)
			}
		} else {
			for x := bw - 1; x >= 0; x-- {
				b := image.Rectangle{
					Min: image.Point{X: x * bucketSize, Y: y * bucketSize},
					Max: image.Point{X: (x+1)*bucketSize - 1, Y: (y+1)*bucketSize - 1},
				}

				buckets = append(buckets, b)
			}
		}
	}

	for i := range buckets {
		clip(&buckets[i], w-1, h-1)
	}

	return buckets
}

func clip(r *image.Rectangle, maxX, maxY int) {
	r.Max.X = int(math.Min(float64(r.Max.X), float64(maxX)))
	r.Max.Y = int(math.Min(float64(r.Max.Y), float64(maxY)))
}

func renderBucket(buckets chan image.Rectangle, scene *Scene, img *image.RGBA, wg *sync.WaitGroup, pixels chan Pixel) {
	defer wg.Done()

	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	for b := range buckets {
		for py := b.Min.Y; py <= b.Max.Y; py++ {
			for px := b.Min.X; px <= b.Max.X; px++ {
				clr := pixelColor(px, py, width, height, scene)
				img.Set(px, py, clr)
				pixels <- Pixel{
					image.Point{X: px, Y: py},
					clr.R,
					clr.G,
					clr.B,
				}
			}
		}
	}
}

func rayColor(r Ray, world *HitableList, depth int) Color {
	if isHit, rec := world.Hit(r, 0.001, math.MaxFloat64); isHit {
		attenuation := &Color{}
		emitted := rec.material.Emitted(rec.u, rec.v, rec.P)
		if depth < 50 {
			isScattered, scattered := rec.material.Scatter(r, rec, attenuation)
			if isScattered {
				clr := rayColor(scattered, world, depth+1)
				return emitted.Add(attenuation.Multiply(clr))
			}

			return emitted
		}

		return emitted
	}

	return Color{}
}

func pixelColor(px, py, width, height int, scene *Scene) color.RGBA {
	var rayClr Color
	for s := 0; s < samples; s++ {
		u := (float64(px) + rand.Float64()) / float64(width)
		v := (float64(height-py) + rand.Float64()) / float64(height)
		r := scene.camera.GetRay(u, v)
		rayClr = rayClr.Add(rayColor(r, &scene.world, 0))
	}
	rayClr = rayClr.DivideScalar(samples)
	rayClrGamma := Color{
		R: math.Sqrt(rayClr.R),
		G: math.Sqrt(rayClr.G),
		B: math.Sqrt(rayClr.B)}

	atomic.AddUint64(&ops, samples)

	return rayClrGamma.ToRGBA()
}

func showProgress(pixCount int, done chan bool) {
	const tickIntervalMs = 1000
	ticker := time.NewTicker(time.Millisecond * tickIntervalMs)
	elapsed := 0
	for {
		select {
		case <-ticker.C:
			progress := float64(atomic.LoadUint64(&ops)) / float64(pixCount*samples)
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

func simpleLightScene(width, height int) Scene {
	perlinTexture := NoiseTexture{scale: 4.0}

	world := HitableList{}
	world = append(world, NewSphere(NewVec3d(0.0, -1000.0, 0.0), 1000, Lambertian{albedo: perlinTexture}))
	world = append(world, NewSphere(NewVec3d(0.0, 2.0, 0.0), 2, Lambertian{albedo: perlinTexture}))

	light := DiffuseLight{emit: ConstantTexture{color: Color{4.0, 4.0, 4.0}}}
	world = append(world, NewSphere(NewVec3d(0.0, 7.0, 0.0), 2.0, light))
	world = append(world, XYRect{x0: 3.0, x1: 5.0, y0: 1.0, y1: 3.0, k: -2.0, material: light})

	lookFrom := NewVec3d(17.0, 2.0, 3.0)
	lookAt := NewVec3d(0.0, 2.0, 0.0)

	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 30.0 //vertical field of view in degrees
	aspectRatio := float64(width) / float64(height)

	distToFocus := 10.0
	aperture := 0.0
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspectRatio, aperture, distToFocus, 0.0, 1.0)

	bvh := HitableList{}
	bvh = append(bvh, NewBVHNode(world, 0.0, 1.0))
	return Scene{camera: cam, world: bvh}
}

func cornellBox(width, height int) Scene {
	red := Lambertian{albedo: ConstantTexture{color: Color{R: 0.65, G: 0.05, B: 0.05}}}
	white := Lambertian{albedo: ConstantTexture{color: Color{R: 0.73, G: 0.73, B: 0.73}}}
	green := Lambertian{albedo: ConstantTexture{color: Color{R: 0.12, G: 0.45, B: 0.15}}}
	light := DiffuseLight{emit: ConstantTexture{color: Color{R: 15, G: 15, B: 15}}}

	world := HitableList{}
	world = append(world, FlipNormals{hitable: YZRect{y0: 0.0, y1: 555.0, z0: 0.0, z1: 555.0, k: 555.0, material: green}})
	world = append(world, YZRect{y0: 0.0, y1: 555.0, z0: 0.0, z1: 555.0, k: 0.0, material: red})
	world = append(world, XZRect{x0: 213.0, x1: 343.0, z0: 227.0, z1: 332.0, k: 554.0, material: light})
	world = append(world, FlipNormals{hitable: XZRect{x0: 0.0, x1: 555.0, z0: 0.0, z1: 555.0, k: 555.0, material: white}})
	world = append(world, XZRect{x0: 0.0, x1: 555.0, z0: 0.0, z1: 555.0, k: 0.0, material: white})
	world = append(world, FlipNormals{hitable: XYRect{x0: 0.0, x1: 555.0, y0: 0.0, y1: 555.0, k: 555.0, material: white}})

	world = append(world,
		Translate{
			hitable: NewRotateY(
				NewBox(
					NewVec3d(0.0, 0.0, 0.0),
					NewVec3d(165.0, 165.0, 165.0),
					white),
				-18.0),
			offset: NewVec3d(130.0, 0.0, 65.0)})

	world = append(world,
		Translate{
			hitable: NewRotateY(
				NewBox(
					NewVec3d(0.0, 0.0, 0.0),
					NewVec3d(165.0, 330.0, 165.0),
					white),
				15.0),
			offset: NewVec3d(265.0, 0.0, 295.0)})

	lookFrom := NewVec3d(278.0, 278.0, -800.0)
	lookAt := NewVec3d(278.0, 278.0, 0.0)

	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 40.0 //vertical field of view in degrees
	aspectRatio := float64(width) / float64(height)

	distToFocus := 10.0
	aperture := 0.0
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspectRatio, aperture, distToFocus, 0.0, 1.0)

	bvh := HitableList{}
	bvh = append(bvh, NewBVHNode(world, 0.0, 1.0))
	return Scene{camera: cam, world: bvh}
}

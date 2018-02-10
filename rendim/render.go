package rendim

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"math"
	"math/rand"
	"os"
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
	scene := earthScene(width, height)

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
		if depth < 50 {
			// isScattered, scattered := rec.material.Scatter(r, rec, attenuation)
			isScattered, _ := rec.material.Scatter(r, rec, attenuation)
			if isScattered {
				// clr := rayColor(scattered, world, depth+1)
				// return attenuation.Multiply(clr)
				return *attenuation
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

func randomScene(width, height int) Scene {
	rnd := rand.New(rand.NewSource(42))

	world := HitableList{}
	checker := CheckerTexture{
		even: ConstantTexture{color: Color{0.2, 0.3, 0.1}},
		odd:  ConstantTexture{color: Color{0.9, 0.9, 0.9}},
	}
	world = append(world, NewSphere(NewVec3d(0.0, -1000.0, 0), 1000, Lambertian{albedo: checker}))
	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMaterial := rnd.Float64()
			center := NewVec3d(float64(a)+0.9*rnd.Float64(), 0.2, float64(b)+0.9*rnd.Float64())
			if center.Subtract(NewVec3d(4.0, 0.2, 0.0)).Length() > 0.9 {
				if chooseMaterial < 0.8 { //diffuse
					world = append(world, NewMovingSphere(center, center.Add(NewVec3d(0, 0.5*rnd.Float64(), 0)), 0.0, 1.0, 0.2,
						Lambertian{albedo: ConstantTexture{color: Color{R: rnd.Float64() * rnd.Float64(), G: rnd.Float64() * rnd.Float64(), B: rnd.Float64() * rnd.Float64()}}}))
				} else if chooseMaterial < 0.95 { //metal
					world = append(world, NewSphere(center, 0.2,
						Metal{albedo: ConstantTexture{color: Color{R: 0.5 * (1.0 + rnd.Float64()), G: 0.5 * (1.0 + rnd.Float64()), B: 0.5 * (1.0 + rnd.Float64())}}, fuzz: 0.5 * rnd.Float64()}))
				} else { //glass
					world = append(world, NewSphere(center, 0.2, Dielectric{refIdx: 1.5}))
				}
			}
		}
	}

	world = append(world, NewSphere(NewVec3d(0.0, 1.0, 0.0), 1.0, Dielectric{refIdx: 1.5}))
	world = append(world, NewSphere(NewVec3d(-4.0, 1.0, 0.0), 1.0, Lambertian{albedo: ConstantTexture{color: Color{R: 0.4, G: 0.2, B: 0.1}}}))
	world = append(world, NewSphere(NewVec3d(4.0, 1.0, 0.0), 1.0, Metal{albedo: ConstantTexture{color: Color{R: 0.7, G: 0.6, B: 0.5}}, fuzz: 0.0}))

	lookFrom := NewVec3d(13.0, 2.0, 3.0)
	lookAt := NewVec3d(0.0, 0.0, 0.0)

	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 20.0 //vertical field of view in degrees
	aspectRatio := float64(width) / float64(height)

	distToFocus := 10.0
	aperture := 0.0
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspectRatio, aperture, distToFocus, 0.0, 1.0)

	bvh := HitableList{}
	bvh = append(bvh, NewBVHNode(world, 0.0, 1.0))
	return Scene{camera: cam, world: bvh}
}

func testScene(width, height int) Scene {
	world := HitableList{}
	world = append(world, NewSphere(NewVec3d(0.0, 0.0, -1.0), 0.5, Lambertian{albedo: ConstantTexture{color: Color{R: 0.1, G: 0.2, B: 0.5}}}))
	world = append(world, NewSphere(NewVec3d(0.0, -100.5, -1.0), 100.0, Lambertian{albedo: ConstantTexture{color: Color{R: 0.8, G: 0.8, B: 0.0}}}))
	world = append(world, NewSphere(NewVec3d(1.0, 0.0, -1.0), 0.5, Metal{albedo: ConstantTexture{color: Color{R: 0.8, G: 0.6, B: 0.2}}, fuzz: 0.0}))
	world = append(world, NewSphere(NewVec3d(-1.0, 0.0, -1.0), 0.5, Dielectric{refIdx: 1.5}))
	world = append(world, NewSphere(NewVec3d(-1.0, 0.0, -1.0), -0.45, Dielectric{refIdx: 1.5}))

	lookFrom := NewVec3d(3.0, 3.0, 2.0)
	lookAt := NewVec3d(0.0, 0.0, -1.0)

	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 20.0 //vertical field of view in degrees
	aspectRatio := float64(width) / float64(height)

	distToFocus := lookFrom.Subtract(lookAt).Length()
	aperture := 0.5
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspectRatio, aperture, distToFocus, 0.0, 1.0)

	return Scene{camera: cam, world: world}
}

func twoPerlinSpheres(width, height int) Scene {
	perlinTexture := NoiseTexture{scale: 4.0}

	world := HitableList{}
	world = append(world, NewSphere(NewVec3d(0.0, -1000.0, 0.0), 1000, Lambertian{albedo: perlinTexture}))
	world = append(world, NewSphere(NewVec3d(0.0, 2.0, 0.0), 2, Lambertian{albedo: perlinTexture}))

	lookFrom := NewVec3d(13.0, 2.0, 3.0)
	lookAt := NewVec3d(0.0, 0.0, 0.0)

	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 20.0 //vertical field of view in degrees
	aspectRatio := float64(width) / float64(height)

	distToFocus := 10.0
	aperture := 0.0
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspectRatio, aperture, distToFocus, 0.0, 1.0)

	bvh := HitableList{}
	bvh = append(bvh, NewBVHNode(world, 0.0, 1.0))
	return Scene{camera: cam, world: bvh}
}

func earthScene(width, height int) Scene {
	f, err := os.Open("earthmap.jpg")
	defer f.Close()
	if err != nil {
		panic("cannot find earthmap.jpg")
	}

	img, _, err := image.Decode(f)
	if err != nil {
		panic("cannot decode earthmap.jpg")
	}

	earth := Lambertian{
		albedo: ImageTexture{
			image: img,
		}}

	world := HitableList{}
	world = append(world, NewSphere(NewVec3d(0.0, 0.0, 0.0), 2, earth))

	lookFrom := NewVec3d(17.0, 2.0, 3.0)
	lookAt := NewVec3d(0.0, 0.0, 0.0)

	vUp := NewVec3d(0.0, 1.0, 0.0)
	vFov := 20.0 //vertical field of view in degrees
	aspectRatio := float64(width) / float64(height)

	distToFocus := 10.0
	aperture := 0.0
	cam := NewCamera(lookFrom, lookAt, vUp, vFov, aspectRatio, aperture, distToFocus, 0.0, 1.0)

	bvh := HitableList{}
	bvh = append(bvh, NewBVHNode(world, 0.0, 1.0))
	return Scene{camera: cam, world: bvh}
}

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
	samples      = 10000
	bucketSize   = 32
	workersCount = 4
)

var ops uint64

type Pixel struct {
	image.Point
	R, G, B uint8
}

func Render(width, height int, pixels chan Pixel) image.Image {
	scene := finalScene(width, height)

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
			fmt.Println("Image rendered in", time.Second*time.Duration(elapsed/1000))
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
	light := DiffuseLight{emit: ConstantTexture{color: Color{R: 7, G: 7, B: 7}}}

	world := HitableList{}
	world = append(world, FlipNormals{hitable: YZRect{y0: 0.0, y1: 555.0, z0: 0.0, z1: 555.0, k: 555.0, material: green}})
	world = append(world, YZRect{y0: 0.0, y1: 555.0, z0: 0.0, z1: 555.0, k: 0.0, material: red})
	world = append(world, XZRect{x0: 113.0, x1: 443.0, z0: 127.0, z1: 432.0, k: 554.0, material: light})
	world = append(world, FlipNormals{hitable: XZRect{x0: 0.0, x1: 555.0, z0: 0.0, z1: 555.0, k: 555.0, material: white}})
	world = append(world, XZRect{x0: 0.0, x1: 555.0, z0: 0.0, z1: 555.0, k: 0.0, material: white})
	world = append(world, FlipNormals{hitable: XYRect{x0: 0.0, x1: 555.0, y0: 0.0, y1: 555.0, k: 555.0, material: white}})

	b1 := Translate{
		hitable: NewRotateY(
			NewBox(
				NewVec3d(0.0, 0.0, 0.0),
				NewVec3d(165.0, 165.0, 165.0),
				white),
			-18.0),
		offset: NewVec3d(130.0, 0.0, 65.0)}

	b2 := Translate{
		hitable: NewRotateY(
			NewBox(
				NewVec3d(0.0, 0.0, 0.0),
				NewVec3d(165.0, 330.0, 165.0),
				white),
			15.0),
		offset: NewVec3d(265.0, 0.0, 295.0)}

	world = append(world, ConstantMedium{boundary: b1, density: 0.01, phaseFunction: Isotropic{albedo: ConstantTexture{color: Color{R: 1.0, G: 1.0, B: 1.0}}}})
	world = append(world, ConstantMedium{boundary: b2, density: 0.01, phaseFunction: Isotropic{albedo: ConstantTexture{color: Color{R: 0.0, G: 0.0, B: 0.0}}}})

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

func finalScene(width, height int) Scene {
	rnd := rand.New(rand.NewSource(3))
	boxlist := HitableList{}
	ground := Lambertian{albedo: ConstantTexture{color: Color{R: 0.48, G: 0.83, B: 0.53}}}

	nb := 20
	for i := 0; i < nb; i++ {
		for j := 0; j < nb; j++ {
			w := 100.0
			x0 := -1000.0 + float64(i)*w
			z0 := -1000.0 + float64(j)*w
			y0 := 0.0
			x1 := x0 + w
			y1 := 100 * (rnd.Float64() + 0.01)
			z1 := z0 + w
			boxlist = append(boxlist, NewBox(
				NewVec3d(x0, y0, z0),
				NewVec3d(x1, y1, z1),
				ground))

		}
	}

	world := HitableList{}
	world = append(world, NewBVHNode(boxlist, 0.0, 1.0))

	light := DiffuseLight{emit: ConstantTexture{color: Color{R: 7, G: 7, B: 7}}}
	world = append(world, XZRect{x0: 123.0, x1: 423.0, z0: 147.0, z1: 412.0, k: 554.0, material: light})

	center := NewVec3d(400.0, 400.0, 200.0)
	world = append(world, NewMovingSphere(center, center.Add(NewVec3d(30.0, 0.0, 0.0)), 0.0, 1.0, 50.0, Lambertian{albedo: ConstantTexture{color: Color{R: 0.7, G: 0.3, B: 0.1}}}))
	world = append(world, NewSphere(NewVec3d(260.0, 150.0, 45.0), 50.0, Dielectric{refIdx: 1.5}))
	world = append(world, NewSphere(NewVec3d(0.0, 150.0, 145.0), 50.0, Metal{albedo: ConstantTexture{color: Color{R: 0.8, G: 0.8, B: 0.9}}, fuzz: 1.0}))

	bnd := NewSphere(NewVec3d(360.0, 150.0, 145.0), 70.0, Dielectric{refIdx: 1.5})
	world = append(world, bnd)
	world = append(world, ConstantMedium{boundary: bnd, density: 0.2, phaseFunction: Isotropic{albedo: ConstantTexture{color: Color{R: 0.2, G: 0.4, B: 0.9}}}})
	bnd2 := NewSphere(NewVec3d(0.0, 0.0, 0.0), 5000.0, Dielectric{refIdx: 1.5})
	world = append(world, ConstantMedium{boundary: bnd2, density: 0.0001, phaseFunction: Isotropic{albedo: ConstantTexture{color: Color{R: 1.0, G: 1.0, B: 1.0}}}})

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

	world = append(world, NewSphere(NewVec3d(400.0, 200.0, 400.0), 100.0, earth))

	perlinTexture := NoiseTexture{scale: 0.1}
	world = append(world, NewSphere(NewVec3d(220.0, 280.0, 300.0), 80.0, Lambertian{albedo: perlinTexture}))

	ns := 1000
	sphereList := HitableList{}
	white := Lambertian{albedo: ConstantTexture{color: Color{R: 0.73, G: 0.73, B: 0.73}}}
	for j := 0; j < ns; j++ {
		sphereList = append(sphereList, NewSphere(NewVec3d(165.0*rnd.Float64(), 165.0*rnd.Float64(), 165.0*rnd.Float64()), 10.0, white))
	}

	sphereBox := Translate{
		hitable: NewRotateY(
			NewBVHNode(sphereList, 0.0, 1.0),
			15.0),
		offset: NewVec3d(-100.0, 270.0, 395.0)}

	world = append(world, sphereBox)

	lookFrom := NewVec3d(478.0, 278.0, -600.0)
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

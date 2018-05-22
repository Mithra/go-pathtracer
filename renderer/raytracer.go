package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"
)

const computeReflectionAndRefractions bool = true

type Pixel struct {
	x     uint
	y     uint
	color Vector3
}

func mix(a float64, b float64, mix float64) float64 {
	return b*mix + a*(1-mix)
}

func trace(ray Ray, scene Scene, depthLeft uint) Vector3 {
	tnear := math.MaxFloat64
	collisionIndex := -1

	var nearestHit = Hit{Valid: false}

	// Compute nearest intersection
	for i := 0; i < len(scene.Objects); i++ {
		hit := scene.Objects[i].Intersects(ray)
		if hit.Valid {
			if hit.Distance < tnear {
				tnear = hit.Distance
				collisionIndex = i
				nearestHit = hit
			}
		}
	}

	if collisionIndex == -1 {
		return Vector3{X: 1, Y: 1, Z: 1}
	}

	collidingObject := scene.Objects[collisionIndex]
	material := collidingObject.Material()
	surfaceColor := Vector3{X: 0, Y: 0, Z: 0}

	// Intersection point
	phit := nearestHit.Position

	// Normal at intersection
	nhit := nearestHit.Normal

	bias := 0.0001
	inside := false

	// Inside the sphere
	if ray.Direction.Dot(nhit) > 0 {
		nhit = nhit.MulScalar(-1)
		inside = true
	}

	if (material.Transparency > 0.0 || material.Reflectivity > 0.0) && depthLeft > 0 {
		rayDotNormal := ray.Direction.Dot(nhit)

		facingRatio := -rayDotNormal
		fresnelEffect := mix(math.Pow(1-facingRatio, 3), 1, 0.1)

		reflectionDir := ray.Direction.Sub(nhit.MulScalar(2 * rayDotNormal)).Normalize()

		reflectionRay := Ray{
			Origin:    phit.Sub(nhit.MulScalar(bias)),
			Direction: reflectionDir,
		}

		reflection := trace(reflectionRay, scene, depthLeft-1)
		refraction := Vector3{X: 0, Y: 0, Z: 0}

		if material.Transparency > 0 {
			ior := 1.1
			eta := ior
			if !inside {
				eta = 1 / ior
			}
			cosi := -nhit.Dot(ray.Direction)
			k := 1 - eta*eta*(1-cosi*cosi)

			refractionDir := ray.Direction.MulScalar(eta).Add(nhit.MulScalar(eta*cosi - math.Sqrt(k))).Normalize()

			refractionRay := Ray{
				Origin:    phit.Sub(nhit.MulScalar(bias)),
				Direction: refractionDir,
			}

			refraction = trace(refractionRay, scene, depthLeft-1)
		}

		reflectionColor := reflection.MulScalar(fresnelEffect)
		refractionColor := refraction.MulScalar(1 - fresnelEffect).MulScalar(material.Transparency)

		if computeReflectionAndRefractions {
			surfaceColor = material.Color.Mul(reflectionColor.Add(refractionColor))
		}
	} else {
		for i := 0; i < len(scene.Lights); i++ {
			light := scene.Lights[i]
			lightDistance := phit.DistanceTo(light.Position)
			lightDirection := light.Position.Sub(phit).Normalize()

			lightRay := Ray{
				Origin:    phit.Add(nhit.MulScalar(bias)),
				Direction: lightDirection,
			}

			transmission := true
			for j := 0; j < len(scene.Objects); j++ {
				// Ignore self collisions
				if j == collisionIndex {
					continue
				}
				lightHit := scene.Objects[j].Intersects(lightRay)

				// Object is blocking the light
				if lightHit.Valid && lightHit.Distance < lightDistance {
					transmission = false
					break
				}
			}

			if transmission {
				intensity := math.Max(0.0, nhit.Dot(lightDirection))
				color := material.Color.Mul(light.EmissionColor).MulScalar(intensity)

				surfaceColor = surfaceColor.Add(color)
			}
		}
	}

	return surfaceColor
}

func worker(options RenderingOptions, scene Scene, inputQueue chan Pixel, outputQueue chan Pixel, wg *sync.WaitGroup) {
	defer wg.Done()

	var aspectRatio = float64(options.Width) / float64(options.Height)
	var angle = math.Tan(0.5 * options.Fov * math.Pi / 180.0)

	for p := range inputQueue {
		x := p.x
		y := p.y

		// Translate from raster space to screen space
		// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-generating-camera-rays/generating-camera-rays

		// Normalized Device Coordinates ([0,1])
		// We add 0.5 because we want to pass through the center of the pixel, not the top left corner
		pixelNdcX := (float64(x) + 0.5) / float64(options.Width)
		pixelNdcY := (float64(y) + 0.5) / float64(options.Height)

		// Screen space ([-1,1])
		pixelScreenX := 2*pixelNdcX - 1
		pixelScreenY := 1 - 2*pixelNdcY

		// Camera space
		pixelCameraX := pixelScreenX * aspectRatio * angle
		pixelCameraY := pixelScreenY * angle

		// if (x == 0 && y == 0) || (x == options.Width-1 && y == 0) || (x == 0 && y == options.Height-1) || (x == options.Width-1 && y == options.Height-1) {
		// 	fmt.Printf("Raster: (X=%v; Y=%v) | Ndc: (X=%v; Y=%v) | Screen: (X=%v; Y=%v) | Camera: (X=%v; Y=%v)\n", x, y, pixelNdcX, pixelNdcY, pixelScreenX, pixelScreenY, pixelCameraX, pixelCameraY)
		// }

		rayDirection := Vector3{X: pixelCameraX, Y: pixelCameraY, Z: -1}.Normalize()

		ray := Ray{
			Origin:    Vector3{X: 0, Y: 0, Z: 0},
			Direction: rayDirection,
		}

		p.color = trace(ray, scene, options.MaxDepth)
		outputQueue <- p
	}
}

// Render computes and write the result as a PNG file
func Render(options RenderingOptions, scene Scene) {
	const maxThreads = 8

	inputQueue := make(chan Pixel)
	outputQueue := make(chan Pixel)

	var wgImageBuilder sync.WaitGroup
	wgImageBuilder.Add(1)

	// Start image building thread
	m := image.NewRGBA(image.Rect(0, 0, int(options.Width), int(options.Height)))
	go func() {
		defer wgImageBuilder.Done()

		for p := range outputQueue {
			c := color.RGBA{
				uint8(math.Min(1, p.color.X) * 255),
				uint8(math.Min(1, p.color.Y) * 255),
				uint8(math.Min(1, p.color.Z) * 255),
				255,
			}
			m.Set(int(p.x), int(p.y), c)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(maxThreads)

	// Spawn workers
	for i := 0; i < maxThreads; i++ {
		go worker(options, scene, inputQueue, outputQueue, &wg)
	}

	// Enqueue all pixels
	for y := uint(0); y < options.Height; y++ {
		for x := uint(0); x < options.Width; x++ {
			inputQueue <- Pixel{x: x, y: y}
		}
	}
	close(inputQueue)

	// Wait for workers to finish
	wg.Wait()
	close(outputQueue)

	// Wait for image to be composed
	wgImageBuilder.Wait()

	f, err := os.OpenFile("rgb.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	fmt.Println("Result saved to rgb.png")
	png.Encode(f, m)
}

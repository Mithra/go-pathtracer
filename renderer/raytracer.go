package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func mix(a float64, b float64, mix float64) float64 {
	return b*mix + a*(1-mix)
}

func trace(ray Ray, objects []Geometry, lights []Light, depth int) Vector3 {
	tnear := math.MaxFloat64
	collisionIndex := -1

	var nearestHit = Hit{Valid: false}

	// Compute nearest intersection
	for i := 0; i < len(objects); i++ {
		hit := objects[i].Intersects(ray)
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

	collidingObject := objects[collisionIndex]
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

	if (collidingObject.Transparency() > 0.0 || collidingObject.Reflectivity() > 0.0) && depth < 5 {
		facingRatio := -ray.Direction.Dot(nhit)
		fresnelEffect := mix(math.Pow(1-facingRatio, 3), 1, 0.1)

		reflectionDir := ray.Direction.Sub(nhit.MulScalar(2 * ray.Direction.Dot(nhit))).Normalize()

		reflectionRay := Ray{
			Origin:    phit.Sub(nhit.MulScalar(bias)),
			Direction: reflectionDir,
		}

		reflection := trace(reflectionRay, objects, lights, depth+1)
		refraction := Vector3{X: 0, Y: 0, Z: 0}

		if collidingObject.Transparency() > 0 {
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

			refraction = trace(refractionRay, objects, lights, depth+1)
		}

		reflectionColor := reflection.MulScalar(fresnelEffect)
		refractionColor := refraction.MulScalar(1 - fresnelEffect).MulScalar(collidingObject.Transparency())

		surfaceColor = collidingObject.Color().Mul(reflectionColor.Add(refractionColor))
	} else {
		for i := 0; i < len(lights); i++ {
			light := lights[i]
			lightDirection := light.Position.Sub(phit).Normalize()

			lightRay := Ray{
				Origin:    phit.Add(nhit.MulScalar(bias)),
				Direction: lightDirection,
			}

			transmission := true
			for j := 0; j < len(objects); j++ {
				lightHit := objects[j].Intersects(lightRay)

				// Object is blocking the light
				if lightHit.Valid {
					transmission = false
					break
				}
			}

			if transmission {
				intensity := math.Max(0.0, nhit.Dot(lightDirection))
				color := collidingObject.Color().Mul(light.EmissionColor).MulScalar(intensity)

				surfaceColor = surfaceColor.Add(color)
			}
		}
	}

	//return collidingObject.Color()
	return surfaceColor
}

// Render computes and write the result as a PNG file
func Render(objects []Geometry, lights []Light) {
	width := 640
	height := 480
	fov := 30.0
	aspectRatio := float64(width) / float64(height)
	angle := math.Tan(0.5 * fov * math.Pi / 180.0)

	m := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {

			// Translate from raster space to screen space
			// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-generating-camera-rays/generating-camera-rays

			// Normalized Device Coordinates ([0,1])
			// We add 0.5 because we want to pass through the center of the pixel, not the top left corner
			pixelNdcX := (float64(x) + 0.5) / float64(width)
			pixelNdcY := (float64(y) + 0.5) / float64(height)

			// Screen space ([-1,1])
			pixelScreenX := 2*pixelNdcX - 1
			pixelScreenY := 1 - 2*pixelNdcY

			// Camera space
			pixelCameraX := pixelScreenX * aspectRatio * angle
			pixelCameraY := pixelScreenY * angle

			if (x == 0 && y == 0) || (x == width-1 && y == 0) || (x == 0 && y == height-1) || (x == width-1 && y == height-1) {
				fmt.Printf("Raster: (X=%v; Y=%v) | Ndc: (X=%v; Y=%v) | Screen: (X=%v; Y=%v) | Camera: (X=%v; Y=%v)\n", x, y, pixelNdcX, pixelNdcY, pixelScreenX, pixelScreenY, pixelCameraX, pixelCameraY)
			}

			rayDirection := Vector3{X: pixelCameraX, Y: pixelCameraY, Z: -1}.Normalize()

			ray := Ray{
				Origin:    Vector3{X: 0, Y: 0, Z: 0},
				Direction: rayDirection,
			}

			pixelColor := trace(ray, objects, lights, 0)

			c := color.RGBA{
				uint8(pixelColor.X * 255),
				uint8(pixelColor.Y * 255),
				uint8(pixelColor.Z * 255),
				255,
			}
			m.Set(x, y, c)
		}
	}

	f, err := os.OpenFile("rgb.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	fmt.Println("Result saved to rgb.png")
	png.Encode(f, m)
}

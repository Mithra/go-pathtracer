package renderer

import (
	"math"
	"math/rand"
)

const computeReflectionAndRefractions bool = true

type RayTracer struct {
}

func (r RayTracer) Sample(x, y uint, camera Camera, scene Scene, options RenderingOptions, rnd *rand.Rand) Vector3 {
	var aspectRatio = float64(options.Width) / float64(options.Height)
	var angle = math.Tan(0.5 * options.Fov * math.Pi / 180.0)

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

	return r.trace(ray, scene, options.MaxDepth)
}

func (r RayTracer) trace(ray Ray, scene Scene, depthLeft uint) Vector3 {
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

		reflection := r.trace(reflectionRay, scene, depthLeft-1)
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

			refraction = r.trace(refractionRay, scene, depthLeft-1)
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

func mix(a float64, b float64, mix float64) float64 {
	return b*mix + a*(1-mix)
}

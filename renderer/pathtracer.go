package renderer

import (
	"math"
	"math/rand"
)

type PathTracer struct {
}

var defaultColor Vector3 = NewVector3(0, 0, 0)

var M_PI float64 = math.Pi
var M_1_PI float64 = 1. / M_PI
var nbSamples int = 4

func (r PathTracer) Sample(x, y uint, camera Camera, scene Scene, options RenderingOptions, rnd *rand.Rand) Vector3 {
	// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-generating-camera-rays/generating-camera-rays
	w := float64(options.Width)
	h := float64(options.Height)

	cam := NewRay(camera.position, camera.direction)

	aspectRatio := w / h
	scale := math.Tan(0.5 * degToRad(options.Fov))

	pixelColor := NewVector3(0, 0, 0)
	for sy := 0; sy < 2; sy++ {
		for sx := 0; sx < 2; sx++ {
			acc := NewVector3(0, 0, 0)
			for subSample := 0; subSample < nbSamples; subSample++ {
				r1 := 2. * rand.Float64()
				dx := ternaryFloat64(r1 < 1, math.Sqrt(r1)-1., 1.-math.Sqrt(2.-r1))

				r2 := 2. * rand.Float64()
				dy := ternaryFloat64(r2 < 1, math.Sqrt(r2)-1., 1.-math.Sqrt(2.-r2))

				// Normalized Device Coordinates ([0,1])
				// We add 0.5 because we want to pass through the center of the pixel, not the top left corner
				pixelNdcX := ((dx+float64(sx))/2. + float64(x) + 0.5) / w
				pixelNdcY := ((dy+float64(sy))/2. + float64(y) + 0.5) / h

				// Screen space ([-1,1])
				pixelScreenX := 2*pixelNdcX - 1
				pixelScreenY := 1 - 2*pixelNdcY // We want the Y axis to go UP, not DOWN so we "inverse" it

				// Camera space (Applying aspect ratio, scale and camera transform)
				pixelCameraX := pixelScreenX * aspectRatio * scale
				pixelCameraY := pixelScreenY * scale
				pixelCameraSpace := camera.cameraToWorld.MultDirection(NewVector3(pixelCameraX, pixelCameraY, -1)).Normalize()

				// Compute color for that pixel
				radiance := r.radiance(NewRay(cam.Origin, pixelCameraSpace), scene, options, 0, rnd, 1)
				acc = acc.Add(radiance.MulScalar(1. / float64(nbSamples)))
			}

			// Sum up pixel color
			pixelColor = pixelColor.Add(NewVector3(clamp(acc.X), clamp(acc.Y), clamp(acc.Z)).MulScalar(.25))
		}
	}

	pixelColor.X = gammaCorrection(pixelColor.X)
	pixelColor.Y = gammaCorrection(pixelColor.Y)
	pixelColor.Z = gammaCorrection(pixelColor.Z)

	return pixelColor
}

func (r PathTracer) intersect(ray Ray, scene Scene, ignoreLights bool) (hit Hit, index int) {
	tnear := math.MaxFloat64
	collisionIndex := -1

	var nearestHit = Hit{Valid: false}

	// Compute nearest intersection
	for i := 0; i < len(scene.Objects); i++ {
		obj := scene.Objects[i]
		if ignoreLights && obj.Material().EmissionColor != Vector3Zero {
			continue
		}

		hit := scene.Objects[i].Intersects(ray)
		if hit.Valid {
			if hit.Distance < tnear {
				tnear = hit.Distance
				collisionIndex = i
				nearestHit = hit
			}
		}
	}

	return nearestHit, collisionIndex
}

func (r PathTracer) radiance(ray Ray, scene Scene, options RenderingOptions, depth uint, rnd *rand.Rand, E float64) Vector3 {
	depth = depth + 1

	// We don't want to draw the light spheres
	ignoreLights := false
	if depth == 1 {
		ignoreLights = true
	}

	nearestHit, collisionIndex := r.intersect(ray, scene, ignoreLights)

	if collisionIndex == -1 {
		return defaultColor
	}

	collidingObject := scene.Objects[collisionIndex]
	material := collidingObject.Material()
	objectColor := material.Color

	// Intersection point
	phit := nearestHit.Position

	// Normal at intersection
	nhit := nearestHit.Normal
	nhitCleaned := nhit

	// Inside the sphere
	if nhit.Dot(ray.Direction) >= 0 {
		nhitCleaned = nhit.MulScalar(-1)
	}

	// Russian Roulette
	p := objectColor.Z
	if objectColor.X > objectColor.Y && objectColor.X > objectColor.Z {
		p = objectColor.X
	} else if objectColor.Y > objectColor.Z {
		p = objectColor.Y
	}

	if depth > 5 || p == 0 {
		if rand.Float64() < p {
			objectColor = objectColor.MulScalar(1 / p)
		} else {
			return material.EmissionColor.MulScalar(E)
		}
	}

	// Pure diffuse material
	if material.Reflectivity == 0 && material.Transparency == 0 {
		r1 := 2. * M_PI * rand.Float64()
		r2 := rand.Float64()
		r2s := math.Sqrt(r2)

		// Create orthonormal coordinate frame (w,u,v)
		w := nhitCleaned
		u := NewVector3(1, 1, 1)
		if math.Abs(w.X) > .1 {
			u = NewVector3(0, 1, 0)
		}
		u = u.Cross(w).Normalize()
		v := w.Cross(u).Normalize()

		d1 := u.MulScalar(math.Cos(r1) * r2s)
		d2 := v.MulScalar(math.Sin(r1) * r2s)
		d3 := w.MulScalar(math.Sqrt(1 - r2))

		// Random reflection ray
		d := d1.Add(d2).Add(d3).Normalize()

		e := NewVector3(0, 0, 0)

		for i := 0; i < len(scene.Objects); i++ {
			if scene.Objects[i].Material().EmissionColor == Vector3Zero {
				continue
			}

			light := scene.Objects[i]
			lightMaterial := light.Material()
			//lightDistance := phit.DistanceTo(light.Position())

			sw := light.Position().Sub(phit)
			su := NewVector3(1, 1, 1)
			if math.Abs(sw.X) > .1 {
				su = NewVector3(0, 1, 0)
			}
			su = su.Cross(sw).Normalize()
			sv := sw.Cross(su).Normalize()

			p := phit.Sub(light.Position())
			rad := 1.5 // TODO: ?

			cos_a_max := math.Sqrt(1 - (rad*rad)/p.Dot(p))
			eps1 := rand.Float64()
			eps2 := rand.Float64()
			cos_a := 1 - eps1 + eps1*cos_a_max
			sin_a := math.Sqrt(1 - cos_a*cos_a)
			phi := 2 * M_PI * eps2

			l1 := su.MulScalar(math.Cos(phi) * sin_a)
			l2 := sv.MulScalar(math.Sin(phi) * sin_a)
			l3 := sw.MulScalar(cos_a)

			lightDirection := l1.Add(l2).Add(l3).Normalize()

			lightRay := NewRay(phit, lightDirection)

			collision, index := r.intersect(lightRay, scene, false)

			// The only collision was the light itself
			if collision.Valid && index == i {
				omega := 2 * M_PI * (1 - cos_a_max)
				e = e.Add(objectColor.Mul(lightMaterial.EmissionColor.MulScalar(lightDirection.Dot(nhitCleaned) * omega)).MulScalar(M_1_PI))
			}
		}

		return material.EmissionColor.MulScalar(E).
			Add(e).
			Add(objectColor.Mul(r.radiance(NewRay(phit, d), scene, options, depth, rnd, 0)))
	}

	reflectionDirection := ray.Direction.Sub(nhit.MulScalar(2 * nhit.Dot(ray.Direction)))
	reflectionRay := NewRay(phit, reflectionDirection)

	// Specular reflection
	if material.Transparency == 0 {
		return material.EmissionColor.
			Add(objectColor.Mul(r.radiance(reflectionRay, scene, options, depth, rnd, 1)))
	}

	// Reflection + Refraction (dielectric (glass))
	into := nhit.Dot(nhitCleaned) > 0
	nc := 1.
	nt := 1.5
	nnt := ternaryFloat64(into, nc/nt, nt/nc)
	ddn := ray.Direction.Dot(nhitCleaned)
	cost2t := 1 - nnt*nnt*(1-ddn*ddn)

	// Total internal reflection
	if cost2t < 0 {
		return material.EmissionColor.
			Add(objectColor.Mul(r.radiance(reflectionRay, scene, options, depth, rnd, 1)))
	}

	// Choose reflection or refraction
	coeff := ternaryFloat64(into, 1, -1)

	tdir1 := ray.Direction.MulScalar(nnt)
	tdir2 := nhit.MulScalar(coeff * (ddn * nnt * math.Sqrt(cost2t)))
	tdir := tdir1.Sub(tdir2).Normalize()

	a := nt - nc
	b := nt + nc
	c := 1 - ternaryFloat64(into, -ddn, tdir.Dot(nhit))
	R0 := a * a / (b * b)
	Re := R0 + (1-R0)*c*c*c*c*c
	Tr := 1 - Re
	P := .25 + .5*Re
	RP := Re / P
	TP := Tr / (1 - P)

	// TODO ?
	colorDelta := Vector3Zero
	if depth > 2 {
		// Russian Roulette
		if rand.Float64() < P {
			colorDelta = r.radiance(reflectionRay, scene, options, depth, rnd, 1).MulScalar(RP)
		} else {
			colorDelta = r.radiance(NewRay(phit, tdir), scene, options, depth, rnd, 1).MulScalar(TP)
		}
	} else {
		c1 := r.radiance(reflectionRay, scene, options, depth, rnd, 1).MulScalar(Re)
		c2 := r.radiance(NewRay(phit, tdir), scene, options, depth, rnd, 1).MulScalar(Tr)
		colorDelta = c1.Add(c2)
	}

	return material.EmissionColor.Add(material.Color.Mul(colorDelta))
}

func clamp(x float64) float64 {
	if x < 0 {
		return 0
	}

	if x > 1 {
		return 1
	}

	return x
}

func gammaCorrection(x float64) float64 {
	// Gamma correction of 2.2
	return math.Pow(clamp(x), 1./2.2)*255. + .5
}

func ternaryFloat64(condition bool, a, b float64) float64 {
	if condition {
		return a
	}
	return b
}

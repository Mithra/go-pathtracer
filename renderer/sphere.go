package renderer

import "math"

type Sphere struct {
	center       Vector3
	radius       float64
	radiusSquare float64
	material     Material
}

func CreateSphere(center Vector3, radius float64, material Material) Sphere {
	return Sphere{
		center:       center,
		radius:       radius,
		radiusSquare: radius * radius,
		material:     material,
	}
}

func (s Sphere) Position() Vector3 {
	return s.center
}

func (s Sphere) Material() Material {
	return s.material
}

func (s Sphere) Intersects(r Ray) Hit {
	// See: https://www.scratchapixel.com/lessons/3d-basic-rendering/minimal-ray-tracer-rendering-simple-shapes/ray-sphere-intersection
	l := s.center.Sub(r.Origin)
	tca := l.Dot(r.Direction)

	if tca < 0 {
		return NoHit
	}

	d2 := l.Dot(l) - tca*tca
	if d2 > s.radiusSquare {
		return NoHit
	}

	thc := math.Sqrt(s.radiusSquare - d2)
	t0 := tca - thc
	t1 := tca + thc

	t := t0
	if t0 < 0 {
		t = t1
	}

	pHit := r.Origin.Add(r.Direction.MulScalar(t))
	nHit := pHit.Sub(s.center).Normalize()

	return Hit{Valid: true, Distance: t, Position: pHit, Normal: nHit}
}

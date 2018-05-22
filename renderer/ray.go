package renderer

import "fmt"

type Ray struct {
	Origin    Vector3
	Direction Vector3
}

func (r Ray) String() string {
	return fmt.Sprintf("Origin = %v; Direction = %v", r.Origin, r.Direction)
}

func NewRay(origin, direction Vector3) Ray {
	return Ray{
		Origin:    origin,
		Direction: direction,
	}
}

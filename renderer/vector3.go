package renderer

import (
	"fmt"
	"math"
)

type Vector3 struct {
	X, Y, Z float64
}

func (v Vector3) String() string {
	return fmt.Sprintf("(%0.24f, %0.24f, %0.24f)", v.X, v.Y, v.Z)
}

func (v Vector3) Add(v2 Vector3) Vector3 {
	return Vector3{X: v.X + v2.X, Y: v.Y + v2.Y, Z: v.Z + v2.Z}
}

func (v Vector3) Sub(v2 Vector3) Vector3 {
	return Vector3{X: v.X - v2.X, Y: v.Y - v2.Y, Z: v.Z - v2.Z}
}

func (v Vector3) MulScalar(f float64) Vector3 {
	return Vector3{X: v.X * f, Y: v.Y * f, Z: v.Z * f}
}

func (v Vector3) Mul(v2 Vector3) Vector3 {
	return Vector3{X: v.X * v2.X, Y: v.Y * v2.Y, Z: v.Z * v2.Z}
}

func (v Vector3) Div(v2 Vector3) Vector3 {
	return Vector3{X: v.X / v2.X, Y: v.Y / v2.Y, Z: v.Z / v2.Z}
}

func (v Vector3) Dot(v2 Vector3) float64 {
	return v.X*v2.X + v.Y*v2.Y + v.Z*v2.Z
}

func (v Vector3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

func (v Vector3) IsZero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

func (v Vector3) Normalize() Vector3 {
	length := v.Length()
	if length > 0 {
		invNor := 1 / length
		return v.MulScalar(invNor)
	}
	return v
}
package renderer

type Geometry interface {
	Position() Vector3
	Material() Material
	Intersects(ray Ray) Hit
}

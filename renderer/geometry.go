package renderer

type Geometry interface {
	Material() Material
	Intersects(ray Ray) Hit
}

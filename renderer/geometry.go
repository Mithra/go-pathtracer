package renderer

type Geometry interface {
	Color() Vector3
	Transparency() float64
	Reflectivity() float64

	Intersects(ray Ray) Hit
}

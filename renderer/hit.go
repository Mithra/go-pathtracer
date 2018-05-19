package renderer

type Hit struct {
	Valid    bool
	Distance float64
	Position Vector3
	Normal   Vector3
}

// NoHit represents a lack of collision between a ray and a geometry
var NoHit = Hit{Valid: false}

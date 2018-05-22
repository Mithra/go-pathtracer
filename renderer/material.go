package renderer

type Material struct {
	Color         Vector3
	EmissionColor Vector3
	Reflectivity  float64
	Transparency  float64
}

func NewMaterial(color Vector3, reflectivity float64, transparency float64, emissionColor Vector3) Material {
	return Material{Color: color, EmissionColor: emissionColor, Reflectivity: reflectivity, Transparency: transparency}
}

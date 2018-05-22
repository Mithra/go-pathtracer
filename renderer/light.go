package renderer

type Light struct {
	Position      Vector3
	EmissionColor Vector3
}

func NewLight(position, emissionColor Vector3) Light {
	return Light{Position: position, EmissionColor: emissionColor}
}

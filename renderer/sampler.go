package renderer

import "math/rand"

type Sampler interface {
	Sample(x, y uint, camera Camera, scene Scene, options RenderingOptions, rnd *rand.Rand) Vector3
}

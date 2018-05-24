package renderer

type Sampler interface {
	Sample(x, y uint, camera Camera, scene Scene, options RenderingOptions) Vector3
}

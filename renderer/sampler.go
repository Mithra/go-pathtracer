package renderer

type Sampler interface {
	Sample(x, y uint, scene Scene, options RenderingOptions) Vector3
}

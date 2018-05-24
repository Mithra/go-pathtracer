package renderer

type Camera struct {
	position  Vector3
	direction Vector3
}

func NewCamera(position, direction Vector3) Camera {
	return Camera{position: position, direction: direction}
}

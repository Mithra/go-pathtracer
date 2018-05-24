package renderer

type Camera struct {
	position      Vector3
	direction     Vector3
	cameraToWorld *Matrix4
}

func NewCamera(position, direction Vector3) Camera {
	return Camera{
		position:      position,
		direction:     direction,
		cameraToWorld: LookAt(position, position.Add(direction.Normalize())),
	}
}

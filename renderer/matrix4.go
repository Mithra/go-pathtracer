package renderer

type Matrix4 struct {
	m [4][4]float64
}

func NewMatrix4() *Matrix4 {
	return &Matrix4{m: [4][4]float64{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}}
}

func NewIdentityMatrix4() *Matrix4 {
	m := NewMatrix4()
	m.m[0][0] = 1
	m.m[1][1] = 1
	m.m[2][2] = 1
	m.m[3][3] = 1

	return m
}

func LookAt(from, to Vector3) *Matrix4 {
	tmp := NewVector3(0, 1, 0)

	forward := from.Sub(to).Normalize()
	right := tmp.Cross(forward)
	up := forward.Cross(right)

	camToWorld := NewIdentityMatrix4()

	camToWorld.m[0][0] = right.X
	camToWorld.m[0][1] = right.Y
	camToWorld.m[0][2] = right.Z

	camToWorld.m[1][0] = up.X
	camToWorld.m[1][1] = up.Y
	camToWorld.m[1][2] = up.Z

	camToWorld.m[2][0] = forward.X
	camToWorld.m[2][1] = forward.Y
	camToWorld.m[2][2] = forward.Z

	camToWorld.m[3][0] = from.X
	camToWorld.m[3][1] = from.Y
	camToWorld.m[3][2] = from.Z

	return camToWorld
}

func (matrix *Matrix4) MultDirection(v Vector3) Vector3 {
	return NewVector3(
		v.X*matrix.m[0][0]+v.Y*matrix.m[1][0]+v.Z*matrix.m[2][0],
		v.X*matrix.m[0][1]+v.Y*matrix.m[1][1]+v.Z*matrix.m[2][1],
		v.X*matrix.m[0][2]+v.Y*matrix.m[1][2]+v.Z*matrix.m[2][2])
}

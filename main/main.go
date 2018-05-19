package main

import (
	"github.com/go-pathtracer/renderer"
)

func main() {
	var objects []renderer.Geometry
	var lights []renderer.Light

	// Ground
	objects = append(objects, renderer.CreateSphere(
		renderer.Vector3{X: 0, Y: -10004, Z: -20},
		10000,
		renderer.Vector3{X: 0.2, Y: 0.2, Z: 0.2},
		0,
		0))

	// Left Wall
	// objects = append(objects, renderer.CreatePlane(
	// 	renderer.Vector3{X: -50, Y: 0, Z: 0},
	// 	renderer.Vector3{X: 1, Y: 0, Z: 0},
	// 	renderer.Vector3{X: 1.0, Y: 0.0, Z: 0.0},
	// 	0,
	// 	0))

	// Far Wall
	// objects = append(objects, renderer.CreatePlane(
	// 	renderer.Vector3{X: 0, Y: 0, Z: -500},
	// 	renderer.Vector3{X: 0, Y: 0, Z: 1},
	// 	renderer.Vector3{X: 0.0, Y: 0.0, Z: 1.0},
	// 	0,
	// 	0))

	// Red sphere
	objects = append(objects, renderer.CreateSphere(
		renderer.Vector3{X: 0, Y: 0, Z: -20},
		4,
		renderer.Vector3{X: 1.0, Y: 0.32, Z: 0.36},
		1,
		0.5))

	// Green sphere
	// objects = append(objects, renderer.CreateSphere(
	// 	renderer.Vector3{X: -3, Y: 7, Z: -30},
	// 	10,
	// 	renderer.Vector3{X: 1.0, Y: 1.32, Z: 0.36},
	// 	1,
	// 	0))

	// Yellow sphere
	objects = append(objects, renderer.CreateSphere(
		renderer.Vector3{X: 5, Y: -1, Z: -15},
		2,
		renderer.Vector3{X: 0.9, Y: 0.76, Z: 0.46},
		1,
		0))

	// Blue sphere
	objects = append(objects, renderer.CreateSphere(
		renderer.Vector3{X: 5, Y: 0, Z: -25},
		3,
		renderer.Vector3{X: 0.65, Y: 0.77, Z: 0.97},
		1,
		0))

	// White sphere
	objects = append(objects, renderer.CreateSphere(
		renderer.Vector3{X: -5.5, Y: 0, Z: -15},
		3,
		renderer.Vector3{X: 0.9, Y: 0.9, Z: 0.9},
		1,
		0))

	// Light
	lights = append(lights, renderer.Light{
		Position:      renderer.Vector3{X: 0, Y: 20, Z: -30},
		EmissionColor: renderer.Vector3{X: 3, Y: 3, Z: 3},
	})

	renderer.Render(objects, lights)
}

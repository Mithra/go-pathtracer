package main

import (
	"github.com/go-pathtracer/renderer"
)

func main() {
	var objects []renderer.Geometry
	var lights []renderer.Light

	options := renderer.RenderingOptions{
		Width:    800,
		Height:   600,
		Fov:      60,
		MaxDepth: 5,
	}

	// BoundingSphere
	//boundingSphereMaterial := renderer.Material{Color: renderer.Vector3{X: 0., Y: 0.2, Z: 1.}, Reflectivity: 0, Transparency: 0}
	//objects = append(objects, renderer.CreateSphere(renderer.Vector3{X: 0, Y: 0, Z: 0}, 10000, boundingSphereMaterial))

	// Ground
	groundMaterial := renderer.Material{Color: renderer.Vector3{X: 1., Y: 1., Z: 1.}, Reflectivity: 1, Transparency: 0}
	objects = append(objects, renderer.CreateSphere(renderer.Vector3{X: 0, Y: -50, Z: -40}, 50, groundMaterial))

	// Red sphere
	redMaterial := renderer.Material{Color: renderer.Vector3{X: 1, Y: 0.32, Z: 0.36}, Reflectivity: 1, Transparency: 0.5}
	objects = append(objects, renderer.CreateSphere(renderer.Vector3{X: 0, Y: 0, Z: -20}, 4, redMaterial))

	// Green sphere
	greenMaterial := renderer.Material{Color: renderer.Vector3{X: 0.5, Y: 1.0, Z: 0.36}, Reflectivity: 1, Transparency: 0}
	objects = append(objects, renderer.CreateSphere(renderer.Vector3{X: -5, Y: 10, Z: -30}, 10, greenMaterial))

	// Yellow sphere
	yellowMaterial := renderer.Material{Color: renderer.Vector3{X: 0.9, Y: 0.76, Z: 0.46}, Reflectivity: 1, Transparency: 0}
	objects = append(objects, renderer.CreateSphere(renderer.Vector3{X: 5, Y: -1, Z: -15}, 2, yellowMaterial))

	// Blue sphere
	blueMaterial := renderer.Material{Color: renderer.Vector3{X: 0.65, Y: 0.77, Z: 0.97}, Reflectivity: 1, Transparency: 0}
	objects = append(objects, renderer.CreateSphere(renderer.Vector3{X: 5, Y: 0, Z: -25}, 3, blueMaterial))

	// White sphere
	whiteMaterial := renderer.Material{Color: renderer.Vector3{X: 0.9, Y: 0.9, Z: 0.9}, Reflectivity: 1, Transparency: 0}
	objects = append(objects, renderer.CreateSphere(renderer.Vector3{X: -5.5, Y: 0, Z: -15}, 3, whiteMaterial))

	// Light
	lights = append(lights, renderer.Light{
		Position:      renderer.Vector3{X: 20, Y: 20, Z: 40},
		EmissionColor: renderer.Vector3{X: 1, Y: 1, Z: 1},
	})

	renderer.Render(options, renderer.Scene{Objects: objects, Lights: lights})
}

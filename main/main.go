package main

import (
	r "github.com/go-pathtracer/renderer"
)

func main() {
	options := r.RenderingOptions{
		Width:  800,
		Height: 600,
		//Fov:    90,
		Fov:      54.36,
		MaxDepth: 5,
	}

	sampler := r.PathTracer{}
	//sampler := r.RayTracer{}

	camera, scene := createCornellBoxScene2()

	r.Render(sampler, options, camera, scene)
}

func createCornellBoxScene2() (r.Camera, r.Scene) {
	lightMaterial := r.NewMaterial(r.Vector3Zero, 0, 0, r.NewVector3(1, 1, 1).MulScalar(400))

	leftWallMaterial := r.NewMaterial(r.NewVector3(.75, .25, .25), 0, 0, r.Vector3Zero)
	rightWallMaterial := r.NewMaterial(r.NewVector3(.25, .25, .75), 0, 0, r.Vector3Zero)
	backWallMaterial := r.NewMaterial(r.NewVector3(.75, .75, .75), 0, 0, r.Vector3Zero)
	frontWallMaterial := r.NewMaterial(r.NewVector3(0, 0, 0), 0, 0, r.Vector3Zero)

	mirrorMaterial := r.NewMaterial(r.NewVector3(.9999, .9999, .9999), 1.0, 0, r.Vector3Zero)
	glassMaterial := r.NewMaterial(r.NewVector3(.9999, .9999, .9999), 1.0, 0, r.Vector3Zero)

	const radius = 1e5

	objects := []r.Geometry{
		r.CreateSphere(r.NewVector3(-radius, 40.8, 81.6), radius, leftWallMaterial),  // Left
		r.CreateSphere(r.NewVector3(1e5+100, 40.8, 81.6), radius, rightWallMaterial), // Right
		r.CreateSphere(r.NewVector3(50, 40.8, -radius), radius, backWallMaterial),    // Back
		r.CreateSphere(r.NewVector3(50, 40.8, 1e5+170), radius, frontWallMaterial),   // Front
		r.CreateSphere(r.NewVector3(50, -radius, -81.6), radius, backWallMaterial),   // Bottom
		r.CreateSphere(r.NewVector3(50, 1e5+81.6, -81.6), radius, backWallMaterial),  // Top

		r.CreateSphere(r.NewVector3(27, 16.5, 47), 16.5, mirrorMaterial), // Sphere1
		r.CreateSphere(r.NewVector3(73, 16.5, 78), 16.5, glassMaterial),  // Sphere2

		r.CreateSphere(r.NewVector3(50, 81.6-16.5, 81.6), 1.5, lightMaterial), // Light
	}

	lights := []r.Light{}

	//camera := r.NewCamera(r.NewVector3(50, 52, 295.6), r.NewVector3(0, -0.042612, -1).Normalize())
	camera := r.NewCamera(r.NewVector3(50, 40.8, 190), r.NewVector3(0, -0.042612, -1).Normalize())

	return camera, r.Scene{Objects: objects, Lights: lights}
}

func createCornellBoxScene() (r.Camera, r.Scene) {
	lightMaterial := r.NewMaterial(r.Vector3Zero, 0, 0, r.NewVector3(1, 1, 1).MulScalar(4))
	leftWallMaterial := r.NewMaterial(r.NewVector3(.75, .25, .25), 0, 0, r.Vector3Zero)
	rightWallMaterial := r.NewMaterial(r.NewVector3(.25, .25, .75), 0, 0, r.Vector3Zero)
	backWallMaterial := r.NewMaterial(r.NewVector3(.75, .75, .75), 0, 0, r.Vector3Zero)
	frontWallMaterial := r.NewMaterial(r.NewVector3(0, 1, 0), 0, 0, r.Vector3Zero)

	mirrorMaterial := r.NewMaterial(r.NewVector3(.9999, .9999, .9999), 1, 0, r.Vector3Zero)
	glassMaterial := r.NewMaterial(r.NewVector3(.9999, .9999, .9999), 1, 0, r.Vector3Zero)

	const radius = 1e4
	const offsetX = 100
	const offsetZ = 50

	objects := []r.Geometry{
		r.CreateSphere(r.NewVector3(-radius-offsetX, 0, -offsetZ), radius, leftWallMaterial), // Left
		r.CreateSphere(r.NewVector3(radius+offsetX, 0, -offsetZ), radius, rightWallMaterial), // Right
		r.CreateSphere(r.NewVector3(0, 0, -radius-offsetZ-300), radius, backWallMaterial),    // Back
		r.CreateSphere(r.NewVector3(0, 0, radius+offsetZ+50), radius, frontWallMaterial),     // Front
		r.CreateSphere(r.NewVector3(0, -radius-offsetX, -offsetZ), radius, backWallMaterial), // Bottom
		r.CreateSphere(r.NewVector3(0, radius+offsetX, -offsetZ), radius, backWallMaterial),  // Top

		r.CreateSphere(r.NewVector3(-50, -75, -offsetZ-250), 30, mirrorMaterial), // Sphere1
		r.CreateSphere(r.NewVector3(50, -75, -offsetZ-200), 30, mirrorMaterial),  // Sphere2
		r.CreateSphere(r.NewVector3(-15, -75, -offsetZ-190), 30, glassMaterial),  // Sphere3

		r.CreateSphere(r.NewVector3(0, offsetX, -offsetZ-150), 30, lightMaterial), // Light
	}

	lights := []r.Light{
		r.NewLight(r.NewVector3(0, offsetX-10, -offsetZ-150), r.NewVector3(4, 4, 4).MulScalar(0.5)),
	}

	camera := r.NewCamera(r.NewVector3(0, 0, 0), r.NewVector3(0, 0, -1))

	return camera, r.Scene{Objects: objects, Lights: lights}
}

func createTestScene1() (r.Camera, r.Scene) {
	var objects []r.Geometry
	var lights []r.Light

	// Ground
	groundMaterial := r.Material{Color: r.Vector3{X: 1., Y: 1., Z: 1.}, Reflectivity: 1, Transparency: 0}
	objects = append(objects, r.CreateSphere(r.Vector3{X: 0, Y: -50, Z: -40}, 50, groundMaterial))

	// Red sphere
	redMaterial := r.Material{Color: r.Vector3{X: 1, Y: 0.32, Z: 0.36}, Reflectivity: 1, Transparency: 0.5}
	objects = append(objects, r.CreateSphere(r.Vector3{X: 0, Y: 0, Z: -20}, 4, redMaterial))

	// Green sphere
	greenMaterial := r.Material{Color: r.Vector3{X: 0.5, Y: 1.0, Z: 0.36}, Reflectivity: 1, Transparency: 0}
	objects = append(objects, r.CreateSphere(r.Vector3{X: -5, Y: 10, Z: -30}, 10, greenMaterial))

	// Yellow sphere
	yellowMaterial := r.Material{Color: r.Vector3{X: 0.9, Y: 0.76, Z: 0.46}, Reflectivity: 1, Transparency: 0}
	objects = append(objects, r.CreateSphere(r.Vector3{X: 5, Y: -1, Z: -15}, 2, yellowMaterial))

	// Blue sphere
	blueMaterial := r.Material{Color: r.Vector3{X: 0.65, Y: 0.77, Z: 0.97}, Reflectivity: 1, Transparency: 0}
	objects = append(objects, r.CreateSphere(r.Vector3{X: 5, Y: 0, Z: -25}, 3, blueMaterial))

	// White sphere
	whiteMaterial := r.Material{Color: r.Vector3{X: 0.9, Y: 0.9, Z: 0.9}, Reflectivity: 1, Transparency: 0}
	objects = append(objects, r.CreateSphere(r.Vector3{X: -5.5, Y: 0, Z: -15}, 3, whiteMaterial))

	// Light
	lights = append(lights, r.Light{
		Position:      r.Vector3{X: 20, Y: 20, Z: 40},
		EmissionColor: r.Vector3{X: 1, Y: 1, Z: 1},
	})

	camera := r.NewCamera(r.NewVector3(0, 0, 0), r.NewVector3(0, 0, -1))

	return camera, r.Scene{Objects: objects, Lights: lights}
}

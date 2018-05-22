package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"
)

type Pixel struct {
	x     uint
	y     uint
	color Vector3
}

// Render computes and write the result as a PNG file
func Render(sampler Sampler, options RenderingOptions, scene Scene) {
	const maxThreads = 8

	inputQueue := make(chan Pixel)
	outputQueue := make(chan Pixel)

	var wgImageBuilder sync.WaitGroup
	wgImageBuilder.Add(1)

	// Start image building thread
	m := image.NewRGBA(image.Rect(0, 0, int(options.Width), int(options.Height)))
	go func() {
		defer wgImageBuilder.Done()

		for p := range outputQueue {
			c := color.RGBA{
				uint8(math.Min(1, p.color.X) * 255),
				uint8(math.Min(1, p.color.Y) * 255),
				uint8(math.Min(1, p.color.Z) * 255),
				255,
			}
			m.Set(int(p.x), int(p.y), c)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(maxThreads)

	// Spawn workers
	for i := 0; i < maxThreads; i++ {
		go worker(sampler, options, scene, inputQueue, outputQueue, &wg)
	}

	// Enqueue all pixels
	for y := uint(0); y < options.Height; y++ {
		for x := uint(0); x < options.Width; x++ {
			inputQueue <- Pixel{x: x, y: y}
		}
	}
	close(inputQueue)

	// Wait for workers to finish
	wg.Wait()
	close(outputQueue)

	// Wait for image to be composed
	wgImageBuilder.Wait()

	f, err := os.OpenFile("rgb.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	fmt.Println("Result saved to rgb.png")
	png.Encode(f, m)
}

func worker(sampler Sampler, options RenderingOptions, scene Scene, inputQueue chan Pixel, outputQueue chan Pixel, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range inputQueue {
		x := p.x
		y := p.y

		p.color = sampler.Sample(x, y, scene, options)
		outputQueue <- p
	}
}

package renderer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
)

type Pixel struct {
	x     uint
	y     uint
	color Vector3
}

// Render computes and write the result as a PNG file
func Render(sampler Sampler, options RenderingOptions, camera Camera, scene Scene) {
	const maxThreads = 8

	inputQueue := make(chan Pixel)
	outputQueue := make(chan Pixel)

	var wgImageBuilder sync.WaitGroup
	wgImageBuilder.Add(1)

	// Start image building thread
	m := image.NewRGBA(image.Rect(0, 0, int(options.Width), int(options.Height)))
	go func() {
		defer wgImageBuilder.Done()

		totalPixels := options.Width * options.Height
		processedPixels := uint(0)
		progress := uint(0)

		precision := uint(25)
		reportingStep := uint(float64(totalPixels) * (float64(precision) / 100.))

		for p := range outputQueue {
			c := color.RGBA{
				uint8(p.color.X),
				uint8(p.color.Y),
				uint8(p.color.Z),
				255,
			}
			m.Set(int(p.x), int(p.y), c)

			processedPixels++
			if processedPixels%reportingStep == 0 {
				progress++
				fmt.Printf("%v%% completed\n", progress*precision)
				writeImage(m)
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(maxThreads)

	// Spawn workers
	for i := 0; i < maxThreads; i++ {
		go worker(sampler, options, camera, scene, inputQueue, outputQueue, &wg)
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

	fmt.Println("Result saved to rgb.png")
	writeImage(m)
}

func worker(sampler Sampler, options RenderingOptions, camera Camera, scene Scene, inputQueue chan Pixel, outputQueue chan Pixel, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range inputQueue {
		x := p.x
		y := p.y

		p.color = sampler.Sample(x, y, camera, scene, options)
		outputQueue <- p
	}
}

func writeImage(m *image.RGBA) {
	f, err := os.OpenFile("rgb.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	png.Encode(f, m)
	defer f.Close()
}

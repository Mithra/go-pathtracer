package renderer

import (
	"math"
)

func degToRad(angle float64) float64 {
	return angle * math.Pi / 180.
}

func radToDeg(angle float64) float64 {
	return angle * 180. / math.Pi
}

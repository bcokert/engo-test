package mathutil

import (
	"github.com/engoengine/math"

	"engo.io/engo"
)

// SolveQuadratic computes the real solutions to the given quadratic.
// It will return 2, 1, or 0 results in the output array
func SolveQuadratic(a, b, c float32) []float32 {
	det := b*b - 4*a*c
	if det < 0 {
		return []float32{}
	}

	negB := -b
	twoA := 2 * a
	sqrt := math.Sqrt(b*b - 4*a*c)

	if engo.FloatEqual(det, 0) {
		return []float32{(negB + sqrt) / twoA}
	}

	return []float32{(negB + sqrt) / twoA, (negB - sqrt) / twoA}
}

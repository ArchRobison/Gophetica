package math32

import (
	"math"
	"testing"
)

func TestTrunc(t *testing.T) {
	ival := []float32{}
	x := float32(1)
	for !math.IsInf(float64(x), 1) {
		ival = append(ival, x)
		x *= 2
	}
	fval := []float32{}
	x = math.Nextafter32(1, -1)
	for x != 0 {
		fval = append(fval, x)
		x = 1 - (1-x)*2
	}
	for _, i := range ival {
		for _, f := range fval {
			x := i + f
			y := Trunc(x)
			assert(y == i || x-i == 1, "Trunc", t)
			z := Trunc(-x)
			assert(z == -y, "Trunc", t)
			assert(Signbit(Trunc(-x)) == !Signbit(y), "Signbit/Trunc", t)
		}
	}
}

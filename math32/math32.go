// Package math32 implements math functions on float32
package math32

import (
	"math"
)

const Pi = math.Pi

func Atan2(y float32, x float32) float32 {
	return float32(math.Atan2(float64(y), float64(x)))
}

func Abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0 // Deals with negative zero case
	}
	return x
}

func Cos(x float32) float32 {
	return float32(math.Cos(float64(x)))
}

func Exp(x float32) float32 {
	return float32(math.Exp(float64(x)))
}

func Hypot(x float32, y float32) float32 {
	// Using 64-bit arithmetic avoids overflow/underflow without
	// the complications of the approach in math.Hypot.
	x64 := float64(x)
	y64 := float64(y)
	return float32(math.Sqrt(x64*x64 + y64*y64))
}

func Max(a, b float32) float32 {
	if a < b {
		return b
	} else {
		return a
	}
}

func Min(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func Signbit(x float32) bool {
	return math.Signbit(float64(x))
}

func Sin(x float32) float32 {
	return float32(math.Sin(float64(x)))
}

func Sincos(θ float32) (sin, cos float32) {
	y, x := math.Sincos(float64(θ))
	sin = float32(y)
	cos = float32(x)
	return
}

func Sqrt(x float32) float32 {
	return float32(math.Sqrt(float64(x)))
}

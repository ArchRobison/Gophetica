// Package math32 implements math functions on float32
package math32

func Trunc(x float32) float32

func trunc(x float32) float32 {
	const m = 1 << 23
	if x != 0 && -m < x && x < m {
		x = float32(int32(x))
	}
	return x
}

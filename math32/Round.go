package math32

// Round 32-bit float to the nearest integer value, rounding to even in case of a tie.
// Sadly there is no analogue of this function in the Go math library.
func Round(x float32) float32

func round(x float32) float32 {
	const m = 1 << 23
	var d float32
	if x == 0 {
		return x
	} else if x > 0 {
		if x >= m {
			return x
		}
		d = 0.5
	} else if x < 0 {
		if x <= -m {
			return x
		}
		d = -0.5
	}
	i := int32(x + d)
	z := float32(i)
	if z-x == d && i&1 != 0 {
		// Rounded .5 to even
		return x - d
	}
	return z
}

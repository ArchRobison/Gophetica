package math32

// Round 32-bit float to the nearest integer value, rounding to even in case of a tie.
// Sadly there is no analogue of this function in the Go math library.
func Round(x float32) float32 {
	if -(1<<24)<x && x<1<<24 {
        i := int32(float64(x)+0.5)
		z := float32(i)
		if z-x==0.5 && i&1!=0 {
			// Rounded .5 up to odd
			return z-1
		}
		return z
	} else {
		// Only integer bits fit
	    return x
	}
}
 

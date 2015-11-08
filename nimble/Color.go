package nimble

// A pixel value.  Currently the format is ARGB, but this might change in the future.
// Use the Pixel constructing functions RGB and Gray to construct pixels.
type Pixel uint32

// Field positions for ARGB format.
const (
	alphaShift = 24
	redShift   = 16
	greenShift = 8
	blueShift  = 0
)

func component(c float32) Pixel {
	return Pixel(c*255 + (0.5 - 1.0/256))
}

// RGB constructs a Pixel from its red, green, and blue components.
// Each component should be in the interval [0,1]
func RGB(red float32, green float32, blue float32) Pixel {
	return component(red)<<redShift |
		component(green)<<greenShift |
		component(blue)<<blueShift |
		0xFF<<alphaShift
}

// Gray constructs a Pixel with equal red, green, and blue components.
// frac should be in the interval [0,1]
func Gray(frac float32) Pixel {
	g := component(frac)
	return g<<redShift | g<<greenShift | g<<blueShift | 0xFF<<alphaShift
}

// Predefined pixel constants.
const (
	Black = Pixel(0xFF000000)
	White = Pixel(0xFFFFFFFF)
)

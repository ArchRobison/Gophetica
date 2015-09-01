package nimble

type Pixel uint32

const (
	alphaShift = 24
	redShift   = 16
	greenShift = 8
	blueShift  = 0
)

func component(c float32) Pixel {
	return Pixel(c*255 + (0.5 - 1.0/256))
}

//  RGB constructs a Pixel from its red, green, and blue components.
func RGB(red float32, green float32, blue float32) Pixel {
	return component(red)<<redShift |
		component(green)<<greenShift |
		component(blue)<<blueShift |
		0xFF<<alphaShift
}

//  RGB constructs a Pixel with equal red, green, and blue components.
func Gray(frac float32) Pixel {
	g := component(frac)
	return g<<redShift | g<<greenShift | g<<blueShift | 0xFF<<alphaShift
}

package nimble

import (
	"bufio"
	"image"
	_ "image/png"
	"os"
)

// c is a value in 0..0xfff
func component16(c uint32) Pixel {
	return Pixel(float32(c)*(255./65535.) + (0.5 - 1.0/256))
}

// ReadPixMap reads a PixMap from a file with the given name.
func ReadPixMap(filename string) (pm PixMap, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	im, _, err := image.Decode(bufio.NewReader(f))
	if err != nil {
		return
	}
	r := im.Bounds()
	w := int32(r.Max.X - r.Min.X)
	h := int32(r.Max.Y - r.Min.Y)
	pm = PixMap{buf: make([]Pixel, h*w), vstride: w, width: w, height: h}
	for i := int32(0); i < h; i++ {
		row := pm.Row(i)
		for j := range row {
			r, g, b, a := im.At(r.Min.X+j, r.Min.Y+int(i)).RGBA()
			row[j] = component16(r)<<redShift |
				component16(g)<<greenShift |
				component16(b)<<blueShift |
				component16(a)<<alphaShift
		}
	}
	return
}

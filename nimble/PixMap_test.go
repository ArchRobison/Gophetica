package nimble

import (
	"testing"
)

func assert(expr bool, msg string, t *testing.T) {
	if !expr {
		t.Fail()
		panic(msg)
	}
}

func TestPixMap(t *testing.T) {
	for _, w := range [...]int32{0, 1, 49, 50} {
		for _, h := range [...]int32{0, 49, 50} {
			for _, v := range []int32{w, w + 1, w + 2} {
				pm := MakePixMap(w, h, make([]Pixel, v*h, v*h), v)
				// Size
				a, b := pm.Size()
				assert(a == w, "Size", t)
				assert(b == h, "Size", t)

				// Width
				assert(pm.Width() == w, "Width", t)

				// Height
				assert(pm.Height() == h, "Height", t)

				// Contains
				assert(!pm.Contains(-1, -1), "Contains", t)
				assert(!pm.Contains(0, -1), "Contains", t)
				assert(!pm.Contains(-1, 0), "Contains", t)
				assert(pm.Contains(0, 0) == !pm.Empty(), "Contains/Empty", t)
				assert(pm.Contains(w-1, h-1) == !pm.Empty(), "Contains/Empty", t)
				assert(!pm.Contains(w, 0), "Contains", t)
				assert(!pm.Contains(0, h), "Contains", t)
				assert(!pm.Contains(w, h), "Contains", t)

				// SetPixel
				c := [3]Pixel{RGB(0, 0.5, 1), RGB(0.5, 1, 0), RGB(1, 0, 0.5)}
				for x := int32(0); x < w; x++ {
					for y := int32(0); y < h; y++ {
						pm.SetPixel(x, y, c[(x+y)%int32(len(c))])
					}
				}
				// Pixel
				for x := int32(0); x < w; x++ {
					for y := int32(0); y < h; y++ {
						assert(pm.Pixel(x, y) == c[(x+y)%int32(len(c))], "SetPixel/Pixel", t)
					}
				}
				// Fill
				d := RGB(.25, .5, .75)
				pm.Fill(d)
				for x := int32(0); x < w; x++ {
					for y := int32(0); y < h; y++ {
						assert(pm.Pixel(x, y) == d, "Fill/Pixel", t)
					}
				}
				// DrawRect
				r := Rect{w / 4, h / 4, w * 3 / 4, h * 3 / 4}
				f := RGB(.75, .25, .5)
				pm.DrawRect(r, f)
				for x := int32(0); x < w; x++ {
					for y := int32(0); y < h; y++ {
						if r.Contains(x, y) {
							assert(pm.Pixel(x, y) == f, "DrawRect/Contains/Pixel", t)
						} else {
							assert(pm.Pixel(x, y) == d, "DrawRect/!Contains/Pixel", t)
						}
					}
				}
			}
		}
	}
}

var Src PixMap = MakePixMap(100, 200, make([]Pixel, 100*200), 100)
var Dst PixMap = MakePixMap(100, 200, make([]Pixel, 100*200), 100)

func BenchmarkCopy(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Dst.Copy(0, 0, &Src)
	}
}

func BenchmarkFill(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Dst.Fill(Gray(0))
	}
}

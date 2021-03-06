package nimble

import (
	"fmt"
)

// A PixMap is a reference to a 2D array or subarray of pixels.
type PixMap struct {
	buf     []Pixel // underlying array of pixels.  [0] is pixel at (0,0).
	vstride int32   // stride between vertically adjacent pixels.  [vstride] is pixel at (0,1).
	width   int32   // width of the array in pixels
	height  int32   // height of array in pixels
}

// MakePixMap makes a PixMap referring to pixels in the given slice.
// The following identities describe the mapping:
//     pixels[0] maps to (0,0).
//     pixels[1] maps to (1,0).
//     pixels[vstride] to (0,1).
func MakePixMap(width, height int32, pixels []Pixel, vstride int32) (pm PixMap) {
	if devConfig {
		// limit is a sanity check limit.  Though a PixMap could conceivably be bigger, it's more likely a sign of a programmer error.
		const limit = 16384
		if uint32(width) > limit || uint32(height) >= limit || vstride < width || int64(height)*int64(vstride) > int64(len(pixels)) {
			panic(fmt.Sprintf("MakePixMap: width=%v height=%v len(pixels)=%v vstride=%v\n", width, height, len(pixels), vstride))
		}
	}
	pm.buf = pixels
	pm.vstride = vstride
	pm.width = width
	pm.height = height
	return
}

// Width returns the width of the PixMap.
func (pm *PixMap) Width() int32 {
	return pm.width
}

// Height returns the height of the PixMap.
func (pm *PixMap) Height() int32 {
	return pm.height
}

// Size returns the width and height of the PixMap
func (pm *PixMap) Size() (w, h int32) {
	w = pm.width
	h = pm.height
	return
}

// Empty is true if the PixMap has zero pixels.
func (pm *PixMap) Empty() bool {
	return pm.width <= 0 || pm.height <= 0
}

// Contains returns true iff the PixMap contains pointer (x,y).
func (pm *PixMap) Contains(x, y int32) bool {
	return uint32(x) < uint32(pm.width) && uint32(y) < uint32(pm.height)
}

// Intersect returns a PixMap referencing the pixels in the intersection of a PixMap and a Rect.
func (pm *PixMap) Intersect(r Rect) (result PixMap) {
	x0, y0, x1, y1 := pm.clip(r)
	if x0 > x1 || y0 > y1 {
		// Empty intersection
		return
	}
	result.buf = pm.buf[x0+y0*pm.vstride : (x1-1)+(y1-1)*pm.vstride+1]
	result.vstride = pm.vstride
	result.width = x1 - x0
	result.height = y1 - y0
	return
}

// Row returns a slice referring to the pixels with the given y coordinate.
// For example:
//     pm.Row(y)[x]
// refers to the same pixel as:
//	   pm.Pixel(x,y)
func (pm *PixMap) Row(y int32) []Pixel {
	i := y * pm.vstride
	return pm.buf[i : i+pm.width]
}

// Pixel returns value of pixel at (x,y)
func (pm *PixMap) Pixel(x int32, y int32) Pixel {
	return pm.buf[pm.index(x, y)]
}

// Set pixel at (x,y) to given color.
func (pm *PixMap) SetPixel(x, y int32, color Pixel) {
	pm.buf[pm.index(x, y)] = color
}

// DrawRect draws a rectangle with the given color.
func (pm *PixMap) DrawRect(r Rect, color Pixel) {
	x0, y0, x1, y1 := pm.clip(r)
	if x1 <= x0 || y1 <= y0 {
		return
	}
	pm.rawDrawRect(x0, y0, x1, y1, color)
}

// Fill fills the PixMap with the given color.
func (pm *PixMap) Fill(color Pixel) {
	pm.rawDrawRect(0, 0, pm.width, pm.height, color)
}

// Copy copies PixMap src to dst, mapping (0,0) of src onto (x0,y0) of dst.
func (dst *PixMap) Copy(x0, y0 int32, src *PixMap) {
	// FIXME - add clipping support
	for y := int32(0); y < src.Height(); y++ {
		s := src.Row(y)
		d := dst.Row(y0 + y)
		copy(d[x0:], s)
	}
}

func (pm *PixMap) index(x, y int32) int32 {
	if devConfig {
		if uint32(x) >= uint32(pm.width) || uint32(y) >= uint32(pm.height) {
			panic(fmt.Sprintf("index: x=%v y=%v width=%v height=%v\n", x, y, pm.width, pm.height))
		}
	}
	return x + y*pm.vstride
}

func max(a, b int32) int32 {
	if a < b {
		return b
	} else {
		return a
	}
}

func min(a, b int32) int32 {
	if a < b {
		return a
	} else {
		return b
	}
}

// Draw rectangle [x0,y0) x [x1,y1)
func (dst *PixMap) rawDrawRect(x0, y0, x1, y1 int32, color Pixel) {
	for y := y0; y < y1; y++ {
		d := dst.Row(y)[x0:x1]
		for j := range d {
			d[j] = color
		}
	}
}

func (pm *PixMap) clip(r Rect) (x0, y0, x1, y1 int32) {
	x0 = max(r.Left, 0)
	x1 = min(r.Right, pm.width)
	y0 = max(r.Top, 0)
	y1 = min(r.Bottom, pm.height)
	return
}

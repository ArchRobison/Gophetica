package nimble

import (
	"testing"
)

func TestRect(t *testing.T) {
	intervals := [][2]int32{
		{4, 4}, // Empty interval
		{-1, 0},
		{0, 10},
		{-30, 10},
		{5, 15},
		{12, 40},
	}
	for _, i := range intervals {
		for _, j := range intervals {
			x0, x1 := i[0], i[1]
			y0, y1 := j[0], j[1]
			r := Rect{x0, y0, x1, y1}
			w, h := r.Size()
			assert(w == x1-x0, "Size", t)
			assert(h == y1-y0, "Size", t)
			assert(r.Width() == w, "Width", t)
			assert(r.Height() == h, "Height", t)
			assert(r.Empty() == (w == 0 || h == 0), "Empty", t)
			assert(!r.Contains(x0-1, y0), "Contains", t)
			assert(!r.Contains(x0, y0-1), "Contains", t)
			assert(r.Contains(x0, y0) == !r.Empty(), "Contains/Empty", t)
			assert(r.Contains(x1-1, y1-1) == !r.Empty(), "Contains/Empty", t)
			assert(!r.Contains(x1-1, y1), "Contains/Empty", t)
			assert(!r.Contains(x1, y1-1), "Contains/Empty", t)
			dx, dy := r.RelativeToLeftTop(6, 15)
			assert(x0+dx == 6, "RelativeToLeftTop/x", t)
			assert(y0+dy == 15, "RelativeToLeftTop/y", t)
		}
	}
}

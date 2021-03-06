package nimble

// A rectangle or bounding box.
// The bounds form the product of half-open intervals [Left,Right) x [Top,Bottom)
type Rect struct {
	Left, Top, Right, Bottom int32
}

// MakeRect makes a Rect given upper left corner (x,y) and its width and height
func MakeRect(x, y, width, height int32) (r Rect) {
	r.Left = x
	r.Top = y
	r.Right = x + width
	r.Bottom = y + height
	return
}

// Empty returns true iff rectangle is empty
func (r *Rect) Empty() bool {
	return r.Right <= r.Left || r.Bottom <= r.Top
}

// Size returns the width and height of the Rect.
func (r *Rect) Size() (w, h int32) {
	w = r.Right - r.Left
	h = r.Bottom - r.Top
	return
}

// Width returns the width the Rect.
func (r *Rect) Width() int32 {
	return r.Right - r.Left
}

// Height returns the height the Rect.
func (r *Rect) Height() int32 {
	return r.Bottom - r.Top
}

// Contains returns truee iff the Rect contains (x,y)
func (r *Rect) Contains(x, y int32) bool {
	return r.Left <= x && x < r.Right && r.Top <= y && y < r.Bottom
}

func (r *Rect) RelativeToLeftTop(x, y int32) (xLocal, yLocal int32) {
	xLocal = x - r.Left
	yLocal = y - r.Top
	return
}

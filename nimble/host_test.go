package nimble

import (
	"fmt"
	"testing"
)

type context struct{}

func (*context) Init( width, height int32 ) {
}

func (*context) Render( pm PixMap ) {
     pm.Fill(Black)
     pm.DrawRect(Rect{Left: observedMouseX, Top: 0, Right: observedMouseX+1, Bottom: pm.Height()},crossColor)
     pm.DrawRect(Rect{Left: 0, Top: observedMouseY, Right: pm.Width(), Bottom: observedMouseY+1},crossColor)
}

var observedMouseX, observedMouseY int32
var crossColor Pixel

func (*context) ObserveMouse(event MouseEvent, x, y int32) {
	observedMouseX, observedMouseY = x, y
	switch event {
	case MouseDown:
		crossColor = RGB(1,0,0)
	case MouseUp:
		crossColor = RGB(0,0,1)
	case MouseMove:
		crossColor = RGB(0.5,0.5,0.5)
	case MouseDrag:
		crossColor = RGB(1,1,1)
	default:
		panic(fmt.Sprintf("ERROR: event=%v x=%v y=%v\n", event, x, y))
	}
}

func (*context) KeyDown(k Key) {
	switch k {
	case KeyEscape:
		Quit()
	}
}

func TestMouse(t *testing.T) {
	AddRenderClient(&context{})
	AddKeyObserver(&context{})
	AddMouseObserver(&context{})
    fmt.Printf("Mouse mouse around and click left button. Press Esc to quit.\n")
	Run()
}
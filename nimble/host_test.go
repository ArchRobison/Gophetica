package nimble

import (
	"fmt"
	"os"
	"testing"
)

type context struct{}

func (*context) Init(width, height int32) {
}

var font *Font

func (*context) Render(pm PixMap) {
	pm.Fill(Black)
	pm.DrawRect(Rect{Left: observedMouseX, Top: 0, Right: observedMouseX + 1, Bottom: pm.Height()}, crossColor)
	pm.DrawRect(Rect{Left: 0, Top: observedMouseY, Right: pm.Width(), Bottom: observedMouseY + 1}, crossColor)
	if font != nil {
		const text = "Quick brown fox"
		w, h := font.Size(text)
		pm.DrawText(observedMouseX-w/2, observedMouseY-h/2, text, RGB(0, 0, 1), font)
	}
}

var observedMouseX, observedMouseY int32
var crossColor Pixel

func (*context) ObserveMouse(event MouseEvent, x, y int32) {
	observedMouseX, observedMouseY = x, y
	switch event {
	case MouseDown:
		crossColor = RGB(1, 0, 0)
	case MouseUp:
		crossColor = RGB(0, 0, 1)
	case MouseMove:
		crossColor = RGB(0.5, 0.5, 0.5)
	case MouseDrag:
		crossColor = RGB(1, 1, 1)
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

type winContext struct{}

func (*winContext) Size() (width, height int32) {
	return 1024, 768
}

func (*winContext) Title() string {
	return "Test"
}

func TestMouse(t *testing.T) {
	var err error
	const fontfile = "Roboto-Regular.ttf"
	font, err = OpenFont("Roboto-Regular.ttf", 20)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not open %s\n", fontfile)
	}
	AddRenderClient(&context{})
	AddKeyObserver(&context{})
	AddMouseObserver(&context{})
	fmt.Printf("Mouse mouse around and click left button. Press Esc to quit.\n")
	var ctx winContext
	Run(&ctx) // Partiel screen test
	Run(nil)  // Full screen test
}

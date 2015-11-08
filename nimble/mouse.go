package nimble

type MouseEvent uint8

// Kinds of mouse events.
const (
	MouseMove = MouseEvent(iota) // Mouse moved with button up
	MouseDown                    // Button pressed
	MouseUp                      // Button released
	MouseDrag                    // Mouse moved with button down
)

// A MouseObserver observes mousse events.
type MouseObserver interface {
	ObserveMouse(event MouseEvent, x, y int32)
}

// Current state of mouse.
var (
	mouseIsDown    bool
	mouseX, mouseY int32
)

// Get state of mouse.
func MouseState() (x, y int32, isDown bool) {
	x = int32(mouseX)
	y = int32(mouseY)
	isDown = mouseIsDown
	return
}

// AddMouseObserver causes m to be notified of mouse events.
func AddMouseObserver(m MouseObserver) {
	mouseObserverList = append(mouseObserverList, m)
}

var mouseObserverList []MouseObserver

func forwardMouseEvent(event MouseEvent, x, y int32) {
	mouseX = x
	mouseY = y
	switch event {
	case MouseDown:
		mouseIsDown = true
	case MouseUp:
		mouseIsDown = false
	}
	for _, m := range mouseObserverList {
		m.ObserveMouse(event, x, y)
	}
}

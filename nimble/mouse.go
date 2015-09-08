package nimble

type MouseEvent uint8

const (
	MouseMove = MouseEvent(iota) // Mouse moved with button up
	MouseDown                    // Button pressed
	MouseUp                      // Button released
	MouseDrag                    // Mouse moved with button down
)

type mouseObserver interface {
	ObserveMouse(event MouseEvent, x, y int32)
}

// Current state of mouse
var (
	mouseIsDown    bool
	mouseX, mouseY int32
)

// Get state mouse
func MouseState() (x, y int32, isDown bool) {
	x = int32(mouseX)
	y = int32(mouseY)
	isDown = mouseIsDown
	return
}

var mouseObserverList []mouseObserver

func AddMouseObserver(m mouseObserver) {
	mouseObserverList = append(mouseObserverList, m)
}

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

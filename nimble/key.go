package nimble

// A Key represents a key on the keyoard.  Printable ASCII values represent themselves.
type Key uint8

// Key values for non-printing keys.
const (
	KeyBackspace = Key(8)
	KeyTab       = Key(9)
	KeyReturn    = Key(0xD)
	KeyEscape    = Key(0x1B)
	KeyDelete    = Key(0x7F)
	KeyLeft      = Key(0x11 + iota) // Borrow "device control" codes
	KeyRight
	KeyUp
	KeyDown
)

// A KeyObserver is notified of key events.
type KeyObserver interface {
	KeyDown(Key)
}

// AddKeyObserver causes KeyObserver k to be notified of key events.
func AddKeyObserver(k KeyObserver) {
	keyObserverList = append(keyObserverList, k)
}

var keyObserverList []KeyObserver

func forwardKeyEvent(k Key) {
	for _, c := range keyObserverList {
		c.KeyDown(k)
	}
}

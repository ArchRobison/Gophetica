package nimble

// A Key represents a key on the keyoard.  Printable ASCII values represent themselves.
type Key uint8

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

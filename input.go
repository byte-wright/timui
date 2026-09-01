package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

// Input is everything that reached the backend since it was last read.
type Input struct {
	// Wheel is the accumulated wheel movement, Y positive when scrolling down.
	Wheel mathi.Vec2

	// Keys are the presses in the order they arrived.
	Keys []KeyEvent
}

// KeyEvent is a single key press. For printable characters Key is KeyRune and Rune
// holds the character, for everything else Key names the key and Rune is zero.
type KeyEvent struct {
	Key  Key
	Rune rune
}

const (
	KeyNone Key = iota
	KeyRune
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyPageUp
	KeyPageDown
	KeyHome
	KeyEnd
	KeyEnter
	KeyEscape
	KeyTab
	KeyBackspace
	KeyDelete
)

// Keys returns the keys pressed since the previous frame. There is no focus, so
// every widget sees all of them and decides for itself whether a key means
// anything to it.
func (t *Timui) Keys() []KeyEvent {
	return t.keys
}

// KeyPressed reports whether the given key was pressed this frame.
func (t *Timui) KeyPressed(key Key) bool {
	for _, k := range t.keys {
		if k.Key == key {
			return true
		}
	}

	return false
}

// RunePressed reports whether the given character was typed this frame.
func (t *Timui) RunePressed(r rune) bool {
	for _, k := range t.keys {
		if k.Key == KeyRune && k.Rune == r {
			return true
		}
	}

	return false
}

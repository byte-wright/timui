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

	// Mods are the modifiers currently held. It is a state and not an event, so
	// it survives being read.
	Mods ModMask
}

// ModMask is a set of modifier keys. Backends track the currently held state:
// the vaxis backend sees a modifier pressed or let go on its own through the
// kitty keyboard protocol. On a backend that can't report those (tcell), the
// modifiers only change when the mouse moves or a key is typed.
type ModMask int

const (
	ModShift ModMask = 1 << iota
	ModCtrl
	ModAlt
)

// Has reports whether every modifier in o is set.
func (m ModMask) Has(o ModMask) bool {
	return m&o == o
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

// Mods returns the modifiers currently held.
func (t *Timui) Mods() ModMask {
	return t.mods
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

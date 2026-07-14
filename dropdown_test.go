package timui

import (
	"strconv"
	"testing"

	"github.com/byte-wright/expect"
	"gitlab.com/bytewright/gmath/mathi"
)

type interactiveBackend struct {
	size    mathi.Vec2
	mouse   mathi.Vec2
	pressed bool
}

func (b *interactiveBackend) Size() mathi.Vec2                             { return b.size }
func (b *interactiveBackend) MousePosition() mathi.Vec2                    { return b.mouse }
func (b *interactiveBackend) MousePressed(key Key) bool                    { return b.pressed }
func (b *interactiveBackend) Set(pos mathi.Vec2, char rune, fg, bg uint32) {}
func (b *interactiveBackend) Render()                                      {}

func TestDropdownOverlaySelectsItem(t *testing.T) {
	be := &interactiveBackend{size: mathi.Vec2{X: 30, Y: 15}, mouse: mathi.Vec2{X: 5, Y: 0}}
	tui := New(be)

	selected := 0
	frame := func() {
		tui.Dropdown("dd", 3, &selected, func(i int, s bool) {
			tui.Label("item " + strconv.Itoa(i))
		})
		tui.Finish()
	}

	// hover, press, release, and let the widget observe the release
	click := func() {
		frame()
		be.pressed = true
		frame()
		be.pressed = false
		frame()
		frame()
	}

	click() // on the collapsed row: opens the overlay

	be.mouse = mathi.Vec2{X: 5, Y: 4} // third item row inside the deferred overlay
	click()

	expect.Value(t, "selected item", selected).ToBe(2)
}

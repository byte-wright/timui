package timui_test

import (
	"strconv"
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui/internal/test"
	"gitlab.com/bytewright/gmath/mathi"
)

func TestDropdownOverlaySelectsItem(t *testing.T) {
	tui, be := test.New(t, 30, 15)
	be.Mouse = mathi.Vec2{X: 5, Y: 0}

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
		be.Pressed = true
		frame()
		be.Pressed = false
		frame()
		frame()
	}

	click() // on the collapsed row: opens the overlay

	be.Mouse = mathi.Vec2{X: 5, Y: 4} // third item row inside the deferred overlay
	click()

	expect.Value(t, "selected item", selected).ToBe(2)
}

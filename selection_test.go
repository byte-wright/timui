package timui_test

import (
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/internal/test"
	"gitlab.com/bytewright/gmath/mathi"
)

// drag runs the press/move/release/settle frames of a mouse drag from start to
// end, calling frame() once per step the way the real app loop would.
func drag(be *test.TestBackend, frame func(), start, end mathi.Vec2) {
	be.Mouse = start
	frame()

	be.Pressed = true
	frame()

	be.Mouse = end
	frame()

	be.Pressed = false
	frame()

	// lets the selection manager observe the release
	frame()
}

func TestSelectionCopiesDraggedText(t *testing.T) {
	tui, be := test.New(t, 30, 5)

	frame := func() {
		tui.Label("Hello, World")
		tui.Label("Second row")
		tui.Finish()
	}

	drag(be, frame, mathi.Vec2{X: 0, Y: 0}, mathi.Vec2{X: 4, Y: 1})

	expect.Value(t, "clipboard", be.Clipboard).ToBe("Hello\nSecon")
}

func TestSelectionDoesNotStartOnAWidget(t *testing.T) {
	tui, be := test.New(t, 30, 5)

	frame := func() {
		tui.Button("ClickMe")
		tui.Finish()
	}

	drag(be, frame, mathi.Vec2{X: 1, Y: 0}, mathi.Vec2{X: 5, Y: 0})

	expect.Value(t, "clipboard", be.Clipboard).ToBe("")
}

func TestSelectionWorksInsideAScrollArea(t *testing.T) {
	tui, be := test.New(t, 30, 5)

	frame := func() {
		tui.ScrollAreaV("scroll", func() {
			tui.Label("Scrolled text here")
		})
		tui.Finish()
	}

	// the wheel-catcher MouseInput covers this same area; a drag over it must
	// still start a selection rather than being blocked like a real widget
	drag(be, frame, mathi.Vec2{X: 0, Y: 0}, mathi.Vec2{X: 7, Y: 0})

	expect.Value(t, "clipboard", be.Clipboard).ToBe("Scrolled")
}

func TestSelectionIsScopedToItsColumn(t *testing.T) {
	tui, be := test.New(t, 30, 5)

	frame := func() {
		tui.Columns(timui.Split().Factor(1, 1),
			func() { tui.Label("Left side text") },
			func() { tui.Label("Right side text") },
		)
		tui.Finish()
	}

	// dragging from inside the left column past its edge must not pull in
	// characters from the right column
	drag(be, frame, mathi.Vec2{X: 0, Y: 0}, mathi.Vec2{X: 29, Y: 0})

	expect.Value(t, "clipboard", be.Clipboard).ToBe("Left side text")
}

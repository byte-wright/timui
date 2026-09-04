package timui

import (
	"testing"

	"gitlab.com/bytewright/gmath/mathi"
)

// recordingBackend is a minimal Backend that remembers the last background
// color pushed to each cell, so the test can check what actually got painted
// without a real terminal.
type recordingBackend struct {
	size    mathi.Vec2
	mouse   mathi.Vec2
	pressed bool
	bg      map[mathi.Vec2]uint32
}

func newRecordingBackend(size mathi.Vec2) *recordingBackend {
	return &recordingBackend{size: size, bg: map[mathi.Vec2]uint32{}}
}

func (b *recordingBackend) Size() mathi.Vec2           { return b.size }
func (b *recordingBackend) MousePosition() mathi.Vec2  { return b.mouse }
func (b *recordingBackend) MousePressed(key Key) bool  { return b.pressed }
func (b *recordingBackend) TakeInput() Input           { return Input{} }
func (b *recordingBackend) Render()                    {}
func (b *recordingBackend) WriteClipboard(text string) {}

func (b *recordingBackend) Set(pos mathi.Vec2, char rune, fg, bg uint32) {
	b.bg[pos] = bg
}

// TestSelectionHighlightsOnlySelectedCells checks the actual pixels a backend
// would receive: the drag renders a highlight over exactly the selected cells,
// with no lag and without disturbing anything outside the selection.
func TestSelectionHighlightsOnlySelectedCells(t *testing.T) {
	be := newRecordingBackend(mathi.Vec2{X: 10, Y: 3})
	tui := New(be)

	frame := func() {
		tui.Label("ABCDEFGHI")
		tui.Finish()
	}

	be.mouse = mathi.Vec2{X: 0, Y: 0}
	frame()

	be.pressed = true
	frame()

	be.mouse = mathi.Vec2{X: 3, Y: 0}
	frame() // still mid-drag: the highlight must already be visible this frame

	want := uint32(tui.Theme.Selection.BG)

	for x := 0; x <= 3; x++ {
		if got := be.bg[mathi.Vec2{X: x, Y: 0}]; got != want {
			t.Errorf("cell %d: bg = %#08x, want %#08x (selection highlight)", x, got, want)
		}
	}

	if got := be.bg[mathi.Vec2{X: 4, Y: 0}]; got == want {
		t.Errorf("cell 4 is outside the drag but was painted with the selection highlight")
	}
}

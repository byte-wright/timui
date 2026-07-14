package timui

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/byte-wright/expect"
	"gitlab.com/bytewright/gmath/mathi"
)

// snapshotBackend accumulates the pushed cell diffs like a real terminal, so
// after any Finish its grid equals that frame's front buffer.
type snapshotBackend struct {
	size    mathi.Vec2
	chars   [][]rune
	mouse   mathi.Vec2
	pressed bool
}

func newSnapshotBackend(x, y int) *snapshotBackend {
	chars := make([][]rune, y)
	for i := range chars {
		chars[i] = make([]rune, x)
		for j := range chars[i] {
			chars[i][j] = ' '
		}
	}

	return &snapshotBackend{
		size:  mathi.Vec2{X: x, Y: y},
		chars: chars,
		mouse: mathi.Vec2{X: -1, Y: -1},
	}
}

func (b *snapshotBackend) Size() mathi.Vec2          { return b.size }
func (b *snapshotBackend) MousePosition() mathi.Vec2 { return b.mouse }
func (b *snapshotBackend) MousePressed(key Key) bool { return b.pressed }
func (b *snapshotBackend) Render()                   {}

func (b *snapshotBackend) Set(pos mathi.Vec2, char rune, fg, bg uint32) {
	if char != 0 {
		b.chars[pos.Y][pos.X] = char
	}
}

// String frames the screen so snapshot files keep their trailing spaces even
// through editors that trim whitespace.
func (b *snapshotBackend) String() string {
	border := "+" + strings.Repeat("-", b.size.X) + "+\n"

	sb := strings.Builder{}
	sb.WriteString(border)
	for _, row := range b.chars {
		sb.WriteString("|" + string(row) + "|\n")
	}
	sb.WriteString(border)

	return sb.String()
}

func TestSnapshotWidgets(t *testing.T) {
	be := newSnapshotBackend(30, 12)
	tui := New(be)

	checkedOn := true
	checkedOff := false
	option := "b"

	tui.Label("widgets:")
	tui.Button("Click me")
	tui.Checkbox("on", &checkedOn)
	tui.Checkbox("off", &checkedOff)
	OptionGroup(tui, "group", &option, func(og *OptionGroupElement[string]) {
		og.Option("Alpha", "a")
		og.Option("Beta", "b")
	})

	tui.Finish()

	expect.Value(t, "widgets", be.String()).ToBeSnapshot("testdata/widgets.txt")
}

func TestSnapshotPanel(t *testing.T) {
	be := newSnapshotBackend(24, 8)
	tui := New(be)

	tui.Panel(func(p *Panel) {
		p.Header(func() { tui.Label(" Panel ") })
		tui.Label("line one")
		p.HLine()
		tui.Label("line two")
	})

	tui.Finish()

	expect.Value(t, "panel", be.String()).ToBeSnapshot("testdata/panel.txt")
}

func TestSnapshotGrid(t *testing.T) {
	be := newSnapshotBackend(40, 12)
	tui := New(be)

	tui.Grid(func(grid *Grid) {
		grid.Rows(Split().Fixed(3).Factor(1),
			func(*GridCell) { tui.Label("top") },
			func(cell *GridCell) {
				cell.Columns(Split().Factor(1, 1),
					func(*GridCell) { tui.Label("left") },
					func(*GridCell) { tui.Label("right") },
				)
			},
		)
	})

	tui.Finish()

	expect.Value(t, "grid", be.String()).ToBeSnapshot("testdata/grid.txt")
}

func TestSnapshotDialog(t *testing.T) {
	be := newSnapshotBackend(40, 12)
	tui := New(be)

	visible := true

	tui.Label("background content")
	tui.Dialog("Settings", &visible, func() {
		tui.Label("dialog body")
	})

	tui.Finish()

	expect.Value(t, "dialog", be.String()).ToBeSnapshot("testdata/dialog.txt")
}

func TestSnapshotScrollArea(t *testing.T) {
	be := newSnapshotBackend(16, 8)
	tui := New(be)

	frame := func() {
		tui.ScrollAreaV("list", func() {
			for i := 0; i < 20; i++ {
				tui.Label(fmt.Sprintf("line %02d", i))
			}
		})
		tui.Finish()
	}

	// second frame sizes the knob from the content height measured in the first
	frame()
	frame()

	expect.Value(t, "scrollarea", be.String()).ToBeSnapshot("testdata/scrollarea.txt")
}

func TestSnapshotDropdownOpen(t *testing.T) {
	be := newSnapshotBackend(20, 10)
	tui := New(be)

	selected := 1
	frame := func() {
		tui.Dropdown("dd", 4, &selected, func(i int, s bool) {
			tui.Label("item " + strconv.Itoa(i))
		})
		tui.Finish()
	}

	be.mouse = mathi.Vec2{X: 5, Y: 0}
	frame()
	be.pressed = true
	frame()
	be.pressed = false
	frame()
	frame()

	expect.Value(t, "open dropdown", be.String()).ToBeSnapshot("testdata/dropdown-open.txt")
}

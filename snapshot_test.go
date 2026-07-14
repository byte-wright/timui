package timui_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/internal/test"
	"gitlab.com/bytewright/gmath/mathi"
)

func TestSnapshotWidgets(t *testing.T) {
	be := test.NewBackend(30, 12)
	tui := timui.New(be)

	checkedOn := true
	checkedOff := false
	option := "b"

	tui.Label("widgets:")
	tui.Button("Click me")
	tui.Checkbox("on", &checkedOn)
	tui.Checkbox("off", &checkedOff)
	timui.OptionGroup(tui, "group", &option, func(og *timui.OptionGroupElement[string]) {
		og.Option("Alpha", "a")
		og.Option("Beta", "b")
	})

	tui.Finish()

	expect.Value(t, "widgets", be.String()).ToBeSnapshot("testdata/widgets.txt")
}

func TestSnapshotPanel(t *testing.T) {
	be := test.NewBackend(24, 8)
	tui := timui.New(be)

	tui.Panel(func(p *timui.Panel) {
		p.Header(func() { tui.Label(" Panel ") })
		tui.Label("line one")
		p.HLine()
		tui.Label("line two")
	})

	tui.Finish()

	expect.Value(t, "panel", be.String()).ToBeSnapshot("testdata/panel.txt")
}

func TestSnapshotGrid(t *testing.T) {
	be := test.NewBackend(40, 12)
	tui := timui.New(be)

	tui.Grid(func(grid *timui.Grid) {
		grid.Rows(timui.Split().Fixed(3).Factor(1),
			func(*timui.GridCell) { tui.Label("top") },
			func(cell *timui.GridCell) {
				cell.Columns(timui.Split().Factor(1, 1),
					func(*timui.GridCell) { tui.Label("left") },
					func(*timui.GridCell) { tui.Label("right") },
				)
			},
		)
	})

	tui.Finish()

	expect.Value(t, "grid", be.String()).ToBeSnapshot("testdata/grid.txt")
}

func TestSnapshotDialog(t *testing.T) {
	be := test.NewBackend(40, 12)
	tui := timui.New(be)

	visible := true

	tui.Label("background content")
	tui.Dialog("Settings", &visible, func() {
		tui.Label("dialog body")
	})

	tui.Finish()

	expect.Value(t, "dialog", be.String()).ToBeSnapshot("testdata/dialog.txt")
}

func TestSnapshotScrollArea(t *testing.T) {
	be := test.NewBackend(16, 8)
	tui := timui.New(be)

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
	be := test.NewBackend(20, 10)
	tui := timui.New(be)

	selected := 1
	frame := func() {
		tui.Dropdown("dd", 4, &selected, func(i int, s bool) {
			tui.Label("item " + strconv.Itoa(i))
		})
		tui.Finish()
	}

	be.Mouse = mathi.Vec2{X: 5, Y: 0}
	frame()
	be.Pressed = true
	frame()
	be.Pressed = false
	frame()
	frame()

	expect.Value(t, "open dropdown", be.String()).ToBeSnapshot("testdata/dropdown-open.txt")
}

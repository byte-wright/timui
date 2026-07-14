package timui_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/internal/test"
	"gitlab.com/bytewright/gmath/mathi"
)

func TestSnapshotWidgets(t *testing.T) {
	tui, be := test.New(t, 30, 12)

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

	be.CheckSnapshot("testdata/widgets.txt")
}

func TestSnapshotPanel(t *testing.T) {
	tui, be := test.New(t, 24, 8)

	tui.Panel(func(p *timui.Panel) {
		p.Header(func() { tui.Label(" Panel ") })
		tui.Label("line one")
		p.HLine()
		tui.Label("line two")
	})

	tui.Finish()

	be.CheckSnapshot("testdata/panel.txt")
}

func TestSnapshotGrid(t *testing.T) {
	tui, be := test.New(t, 40, 12)

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

	be.CheckSnapshot("testdata/grid.txt")
}

func TestSnapshotDialog(t *testing.T) {
	tui, be := test.New(t, 40, 12)

	visible := true

	tui.Label("background content")
	tui.Dialog("Settings", &visible, func() {
		tui.Label("dialog body")
	})

	tui.Finish()

	be.CheckSnapshot("testdata/dialog.txt")
}

func TestSnapshotScrollArea(t *testing.T) {
	tui, be := test.New(t, 16, 8)

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

	be.CheckSnapshot("testdata/scrollarea.txt")
}

func TestSnapshotDropdownOpen(t *testing.T) {
	tui, be := test.New(t, 20, 10)

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

	be.CheckSnapshot("testdata/dropdown-open.txt")
}

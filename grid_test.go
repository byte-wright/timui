package timui_test

import (
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/internal/test"
	"gitlab.com/bytewright/gmath/mathi"
)

func TestGridCellAreas(t *testing.T) {
	tui := timui.New(test.NewBackend(20, 10))

	areas := []mathi.Box2{}
	cell := func(*timui.GridCell) { areas = append(areas, *tui.CurrentArea()) }

	tui.Grid(func(grid *timui.Grid) {
		grid.Rows(timui.Split().Fixed(2).Factor(1), cell, cell)
	})

	expect.Value(t, "cell areas", areas).ToBe([]mathi.Box2{
		{From: mathi.Vec2{X: 1, Y: 1}, To: mathi.Vec2{X: 19, Y: 3}},
		{From: mathi.Vec2{X: 1, Y: 4}, To: mathi.Vec2{X: 19, Y: 9}},
	})
}

func TestGridDividerJunctions(t *testing.T) {
	tui := timui.New(test.NewBackend(20, 10))

	tui.Grid(func(grid *timui.Grid) {
		grid.Rows(timui.Split().Fixed(2).Factor(1),
			func(*timui.GridCell) {},
			func(cell *timui.GridCell) {
				cell.Columns(timui.Split().Factor(1, 1),
					func(*timui.GridCell) {},
					func(*timui.GridCell) {},
				)
			},
		)
	})

	glyph := func(x, y int) string {
		return string(tui.FrontCharForTest(mathi.Vec2{X: x, Y: y}))
	}

	expect.Value(t, "row divider left junction", glyph(0, 3)).ToBe("╠")
	expect.Value(t, "row divider line", glyph(5, 3)).ToBe("═")
	expect.Value(t, "row divider right junction", glyph(19, 3)).ToBe("╣")

	expect.Value(t, "column divider top junction", glyph(10, 3)).ToBe("╦")
	expect.Value(t, "column divider line", glyph(10, 6)).ToBe("║")
	expect.Value(t, "column divider bottom junction", glyph(10, 9)).ToBe("╩")
}

func TestGridPanicsOnCellCountMismatch(t *testing.T) {
	tui := timui.New(test.NewBackend(20, 10))

	defer func() {
		expect.Value(t, "recovered panic", recover() != nil).ToBe(true)
	}()

	tui.Grid(func(grid *timui.Grid) {
		grid.Rows(timui.Split().Factor(1, 1), func(*timui.GridCell) {})
	})
}

func TestGridLeavesSplitOptionsUntouched(t *testing.T) {
	tui := timui.New(test.NewBackend(20, 10))

	split := timui.Split().Factor(1, 1)

	tui.Grid(func(grid *timui.Grid) {
		grid.Rows(split, func(*timui.GridCell) {}, func(*timui.GridCell) {})
	})

	expect.Value(t, "split entries after use", split.NumSplitsForTest()).ToBe(2)
}

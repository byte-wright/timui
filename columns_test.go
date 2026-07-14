package timui_test

import (
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/internal/test"
	"gitlab.com/bytewright/gmath/mathi"
)

func TestColumnsCellAreas(t *testing.T) {
	tui := timui.New(test.NewBackend(20, 10))

	areas := []mathi.Box2{}
	cell := func() { areas = append(areas, *tui.CurrentArea()) }

	tui.Columns(timui.Split().Fixed(5).Factor(1), cell, cell)

	expect.Value(t, "column areas", areas).ToBe([]mathi.Box2{
		{From: mathi.Vec2{X: 0, Y: 0}, To: mathi.Vec2{X: 5, Y: 10}},
		{From: mathi.Vec2{X: 5, Y: 0}, To: mathi.Vec2{X: 20, Y: 10}},
	})
}

func TestColumnsAdvanceCursorToTallestColumn(t *testing.T) {
	tui := timui.New(test.NewBackend(20, 10))

	tui.Columns(timui.Split().Factor(1, 1),
		func() { tui.MoveCursorForTest(mathi.Vec2{Y: 3}) },
		func() { tui.MoveCursorForTest(mathi.Vec2{Y: 5}) },
	)

	expect.Value(t, "parent cursor after columns", tui.CurrentArea().From.Y).ToBe(5)
}

func TestColumnsPanicsOnCellCountMismatch(t *testing.T) {
	tui := timui.New(test.NewBackend(20, 10))

	defer func() {
		expect.Value(t, "recovered panic", recover() != nil).ToBe(true)
	}()

	tui.Columns(timui.Split().Factor(1, 1), func() {})
}

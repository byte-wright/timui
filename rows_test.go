package timui_test

import (
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/internal/test"
	"gitlab.com/bytewright/gmath/mathi"
)

func TestRowsCellAreas(t *testing.T) {
	tui, _ := test.New(t, 20, 10)

	areas := []mathi.Box2{}
	cell := func() { areas = append(areas, *tui.CurrentArea()) }

	tui.Rows(timui.Split().Fixed(2).Factor(1), cell, cell)

	expect.Value(t, "row areas", areas).ToBe([]mathi.Box2{
		{From: mathi.Vec2{X: 0, Y: 0}, To: mathi.Vec2{X: 20, Y: 2}},
		{From: mathi.Vec2{X: 0, Y: 2}, To: mathi.Vec2{X: 20, Y: 10}},
	})
}

func TestRowsPanicsOnCellCountMismatch(t *testing.T) {
	tui, _ := test.New(t, 20, 10)

	defer func() {
		expect.Value(t, "recovered panic", recover() != nil).ToBe(true)
	}()

	tui.Rows(timui.Split().Factor(1, 1), func() {})
}

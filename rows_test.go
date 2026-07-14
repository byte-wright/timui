package timui

import (
	"testing"

	"github.com/byte-wright/expect"
	"gitlab.com/bytewright/gmath/mathi"
)

type testBackend struct {
	size mathi.Vec2
}

func (b *testBackend) Size() mathi.Vec2                             { return b.size }
func (b *testBackend) MousePosition() mathi.Vec2                    { return mathi.Vec2{X: -1, Y: -1} }
func (b *testBackend) MousePressed(key Key) bool                    { return false }
func (b *testBackend) Set(pos mathi.Vec2, char rune, fg, bg uint32) {}
func (b *testBackend) Render()                                      {}

func TestRowsCellAreas(t *testing.T) {
	tui := New(&testBackend{size: mathi.Vec2{X: 20, Y: 10}})

	areas := []mathi.Box2{}
	cell := func() { areas = append(areas, *tui.CurrentArea()) }

	tui.Rows(Split().Fixed(2).Factor(1), cell, cell)

	expect.Value(t, "row areas", areas).ToBe([]mathi.Box2{
		{From: mathi.Vec2{X: 0, Y: 0}, To: mathi.Vec2{X: 20, Y: 2}},
		{From: mathi.Vec2{X: 0, Y: 2}, To: mathi.Vec2{X: 20, Y: 10}},
	})
}

func TestRowsPanicsOnCellCountMismatch(t *testing.T) {
	tui := New(&testBackend{size: mathi.Vec2{X: 20, Y: 10}})

	defer func() {
		expect.Value(t, "recovered panic", recover() != nil).ToBe(true)
	}()

	tui.Rows(Split().Factor(1, 1), func() {})
}

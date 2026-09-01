package timui_test

import (
	"fmt"
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/internal/test"
	"gitlab.com/bytewright/gmath/mathi"
)

func listFrame(tui *timui.Timui, state *timui.ScrollState, count int, drawn *[]int) func() {
	return func() {
		*drawn = (*drawn)[:0]

		tui.ScrollListV("log", state, count, func(i int) {
			*drawn = append(*drawn, i)
			tui.Label(fmt.Sprintf("line %02d", i))
		})

		tui.Finish()
	}
}

func TestScrollListDrawsOnlyVisibleRows(t *testing.T) {
	tui, _ := test.New(t, 16, 8)

	state := &timui.ScrollState{}
	drawn := []int{}

	listFrame(tui, state, 10000, &drawn)()

	expect.Value(t, "drawn rows", len(drawn)).ToBe(8)
	expect.Value(t, "view height", state.ViewHeight).ToBe(8)
}

func TestScrollListFollowsTail(t *testing.T) {
	tui, _ := test.New(t, 16, 8)

	state := &timui.ScrollState{Follow: true}
	drawn := []int{}

	listFrame(tui, state, 20, &drawn)()

	expect.Value(t, "first visible row", drawn[0]).ToBe(12)

	state.PageUp()
	listFrame(tui, state, 20, &drawn)()

	expect.Value(t, "follow after paging up", state.Follow).ToBe(false)
	expect.Value(t, "first visible row after paging up", drawn[0]).ToBe(4)

	state.Bottom()
	listFrame(tui, state, 20, &drawn)()

	expect.Value(t, "first visible row after jumping back", drawn[0]).ToBe(12)
}

func TestScrollListWheelScrolls(t *testing.T) {
	tui, be := test.New(t, 16, 8)

	state := &timui.ScrollState{Follow: true}
	drawn := []int{}
	frame := listFrame(tui, state, 40, &drawn)

	be.Mouse = mathi.Vec2{X: 4, Y: 4}
	frame()

	expect.Value(t, "first visible row", drawn[0]).ToBe(32)

	// the wheel is dispatched to whatever was hovered during the frame, so it
	// arrives in the one after
	be.Input.Wheel = mathi.Vec2{Y: -2}
	frame()
	frame()

	expect.Value(t, "follow after scrolling up", state.Follow).ToBe(false)
	expect.Value(t, "first visible row after scrolling up", drawn[0]).ToBe(26)
}

func TestScrollListClampsToContent(t *testing.T) {
	tui, _ := test.New(t, 16, 8)

	state := &timui.ScrollState{Offset: 500}
	drawn := []int{}

	listFrame(tui, state, 3, &drawn)()

	expect.Value(t, "offset for content that fits", state.Offset).ToBe(0)
	expect.Value(t, "drawn rows", len(drawn)).ToBe(3)
}

func TestSnapshotScrollList(t *testing.T) {
	tui, be := test.New(t, 16, 8)

	state := &timui.ScrollState{Offset: 4}
	drawn := []int{}

	listFrame(tui, state, 20, &drawn)()

	be.CheckSnapshot("testdata/scrolllist.txt")
}

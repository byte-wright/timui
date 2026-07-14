package timui_test

import (
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/internal/test"
	"gitlab.com/bytewright/gmath/mathi"
)

func TestPanelAreas(t *testing.T) {
	tui, _ := test.New(t, 20, 10)

	bodyArea := mathi.Box2{}
	headerArea := mathi.Box2{}

	tui.Panel(func(p *timui.Panel) {
		bodyArea = *tui.CurrentArea()

		p.Header(func() {
			headerArea = *tui.CurrentArea()
		})
	})

	expect.Value(t, "body area", bodyArea).
		ToBe(mathi.Box2{From: mathi.Vec2{X: 2, Y: 1}, To: mathi.Vec2{X: 18, Y: 9}})
	expect.Value(t, "header area", headerArea).
		ToBe(mathi.Box2{From: mathi.Vec2{X: 2, Y: 0}, To: mathi.Vec2{X: 18, Y: 1}})
	expect.Value(t, "area stack restored", *tui.CurrentArea()).
		ToBe(mathi.Box2{To: mathi.Vec2{X: 20, Y: 10}})
}

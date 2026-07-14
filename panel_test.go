package timui

import (
	"testing"

	"github.com/byte-wright/expect"
	"gitlab.com/bytewright/gmath/mathi"
)

func TestPanelAreas(t *testing.T) {
	tui := New(&testBackend{size: mathi.Vec2{X: 20, Y: 10}})

	bodyArea := mathi.Box2{}
	headerArea := mathi.Box2{}

	tui.Panel(func(p *Panel) {
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

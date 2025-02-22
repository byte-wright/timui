package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

type Panel struct {
	t    *Timui
	area mathi.Box2
}

func (t *Timui) Panel() *Panel {
	t.Border(t.Theme.BorderStyle.Rect, t.Theme.BorderLine, t.Theme.BorderBG)

	area := *t.CurrentArea()
	originalArea := area

	pad := mathi.Vec2{X: 2, Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	t.PushArea(area)

	return &Panel{
		t:    t,
		area: originalArea,
	}
}

func (p *Panel) HLine() {
	pos := p.t.CurrentArea()
	y := pos.From.Y

	a := p.area
	a.From.Y = y

	p.t.PushArea(a)

	p.t.HLine(p.t.Theme.BorderStyle.Horizontal, p.t.Theme.BorderLine, p.t.Theme.BorderBG)

	p.t.PopArea()
}

func (p *Panel) Header() {
	p.t.PopArea()
	area := *p.t.CurrentArea()
	area.To.Y = area.From.Y + 1
	area.From.X += 2
	area.To.X -= 2
	p.t.PushArea(area)
}

func (p *Panel) Finish() {
	p.t.PopArea()
}

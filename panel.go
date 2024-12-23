package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

var (
	borderStyleDouble = [6]rune{'═', '║', '╔', '╗', '╚', '╝'}
	hLineStyleDouble  = [3]rune{'╠', '═', '╣'}
)

type Panel[B Backend] struct {
	t    *Timui[B]
	area mathi.Box2
}

func (t *Timui[B]) Panel() *Panel[B] {
	t.Border(borderStyleDouble)

	area := *t.CurrentArea()
	originalArea := area

	pad := mathi.Vec2{X: 2, Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	t.PushArea(area)

	t.Clear(' ')

	return &Panel[B]{
		t:    t,
		area: originalArea,
	}
}

func (p *Panel[B]) HLine() {
	pos := p.t.CurrentArea()
	y := pos.From.Y

	a := p.area
	a.From.Y = y

	p.t.PushArea(a)

	fg := RGB(0xff, 0xff, 0xff)
	bg := RGB(0x00, 0x00, 0x00)

	p.t.HLine(hLineStyleDouble, fg, bg)

	p.t.PopArea()
}

func (p *Panel[B]) Header() {
	p.t.PopArea()
	area := *p.t.CurrentArea()
	area.To.Y = area.From.Y + 1
	area.From.X += 2
	area.To.X -= 2
	p.t.PushArea(area)
}

func (p *Panel[B]) Finish() {
	p.t.PopArea()
}

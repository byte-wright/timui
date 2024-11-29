package timui

import "gitlab.com/bytewright/gmath/mathi"

var borderStyleDouble = [6]rune{'═', '║', '╔', '╗', '╚', '╝'}

type Panel[B Backend] struct {
	t *Timui[B]
}

func (t *Timui[B]) Panel() *Panel[B] {
	t.Border(borderStyleDouble)

	area := *t.CurrentArea()

	pad := mathi.Vec2{X: 2, Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	t.PushArea(area)

	t.Clear(' ')

	return &Panel[B]{
		t: t,
	}
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

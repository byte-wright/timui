package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

type Panel struct {
	t    *Timui
	area mathi.Box2
}

// Panel draws a bordered container and runs body inside its padded area.
func (t *Timui) Panel(body func(p *Panel)) {
	t.Border(t.Theme.BorderStyle.Rect, t.Theme.BorderLine, t.Theme.BorderBG)

	area := *t.CurrentArea()
	p := &Panel{t: t, area: area}

	pad := mathi.Vec2{X: 2, Y: 1}
	padded := area
	padded.From = padded.From.Add(pad)
	padded.To = padded.To.Sub(pad)

	t.WithArea(padded, func() {
		body(p)
	})
}

func (p *Panel) HLine() {
	a := p.area
	a.From.Y = p.t.CurrentArea().From.Y

	p.t.WithArea(a, func() {
		p.t.HLine(p.t.Theme.BorderStyle.Horizontal, p.t.Theme.BorderLine, p.t.Theme.BorderBG)
	})

	p.t.moveCursor(mathi.Vec2{Y: 1})
}

// Header runs body inside a one-row strip on the panel's top border line.
func (p *Panel) Header(body func()) {
	area := p.area
	area.To.Y = area.From.Y + 1
	area.From.X += 2
	area.To.X -= 2

	p.t.WithArea(area, body)
}

package timui

type Pad struct {
	t *Timui
}

func (t *Timui) Pad(top, right, bottom, left int) *Pad {
	area := *t.CurrentArea()

	area.From.X += left
	area.From.Y += top
	area.To.X -= right
	area.To.Y -= bottom

	t.PushArea(area)

	return &Pad{
		t: t,
	}
}

func (p *Pad) Finish() {
	p.t.PopArea()
}

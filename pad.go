package timui

type Pad[B Backend] struct {
	t *Timui[B]
}

func (t *Timui[B]) Pad(top, right, bottom, left int) *Pad[B] {
	area := *t.CurrentArea()

	area.From.X += left
	area.From.Y += top
	area.To.X -= right
	area.To.Y -= bottom

	t.PushArea(area)

	return &Pad[B]{
		t: t,
	}
}

func (p *Pad[B]) Finish() {
	p.t.PopArea()
}

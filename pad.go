package timui

func (t *Timui) Pad(top, right, bottom, left int, body func()) {
	area := *t.CurrentArea()

	area.From.X += left
	area.From.Y += top
	area.To.X -= right
	area.To.Y -= bottom

	t.PushArea(area)

	body()

	t.PopArea()
}

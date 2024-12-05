package timui

import "gitlab.com/bytewright/gmath/mathi"

func (t *Timui[B]) Border(style [6]rune) {
	c := t.CurrentArea()

	for x := c.From.X + 1; x < c.To.X-1; x++ {
		t.Set(mathi.Vec2{X: x, Y: c.From.Y}, style[0])
		t.Set(mathi.Vec2{X: x, Y: c.To.Y - 1}, style[0])
	}

	for y := c.From.Y + 1; y < c.To.Y-1; y++ {
		t.Set(mathi.Vec2{X: c.From.X, Y: y}, style[1])
		t.Set(mathi.Vec2{X: c.To.X - 1, Y: y}, style[1])
	}

	t.Set(c.From, style[2])
	t.Set(mathi.Vec2{X: c.To.X - 1, Y: c.From.Y}, style[3])
	t.Set(mathi.Vec2{X: c.From.X, Y: c.To.Y - 1}, style[4])
	t.Set(mathi.Vec2{X: c.To.X - 1, Y: c.To.Y - 1}, style[5])
}

func (t *Timui[B]) Clear(char rune) {
	// panic("not implemented")
}

func (t *Timui[B]) Set(pos mathi.Vec2, char rune) {
	clip := t.PeekClip()
	if clip.Contains(pos) {
		t.front.set(pos, char)
	}
}

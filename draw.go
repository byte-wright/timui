package timui

import "gitlab.com/bytewright/gmath/mathi"

func (t *Timui[B]) Border(style [6]rune) {
	c := t.CurrentArea()

	for x := c.From.X + 1; x < c.To.X-1; x++ {
		t.front.set(mathi.Vec2{X: x, Y: c.From.Y}, style[0])
		t.front.set(mathi.Vec2{X: x, Y: c.To.Y - 1}, style[0])
	}

	for y := c.From.Y + 1; y < c.To.Y-1; y++ {
		t.front.set(mathi.Vec2{X: c.From.X, Y: y}, style[1])
		t.front.set(mathi.Vec2{X: c.To.X - 1, Y: y}, style[1])
	}

	t.front.set(c.From, style[2])
	t.front.set(mathi.Vec2{X: c.To.X - 1, Y: c.From.Y}, style[3])
	t.front.set(mathi.Vec2{X: c.From.X, Y: c.To.Y - 1}, style[4])
	t.front.set(mathi.Vec2{X: c.To.X - 1, Y: c.To.Y - 1}, style[5])
}

func (t *Timui[B]) Clear(char rune) {
}

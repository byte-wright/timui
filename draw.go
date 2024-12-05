package timui

import "gitlab.com/bytewright/gmath/mathi"

func (t *Timui[B]) Border(style [6]rune) {
	c := t.CurrentArea()

	fg := int32(0xffffff)
	bg := int32(0x000000)

	for x := c.From.X + 1; x < c.To.X-1; x++ {
		t.Set(mathi.Vec2{X: x, Y: c.From.Y}, style[0], fg, bg)
		t.Set(mathi.Vec2{X: x, Y: c.To.Y - 1}, style[0], fg, bg)
	}

	for y := c.From.Y + 1; y < c.To.Y-1; y++ {
		t.Set(mathi.Vec2{X: c.From.X, Y: y}, style[1], fg, bg)
		t.Set(mathi.Vec2{X: c.To.X - 1, Y: y}, style[1], fg, bg)
	}

	t.Set(c.From, style[2], fg, bg)
	t.Set(mathi.Vec2{X: c.To.X - 1, Y: c.From.Y}, style[3], fg, bg)
	t.Set(mathi.Vec2{X: c.From.X, Y: c.To.Y - 1}, style[4], fg, bg)
	t.Set(mathi.Vec2{X: c.To.X - 1, Y: c.To.Y - 1}, style[5], fg, bg)
}

func (t *Timui[B]) Clear(char rune) {
	// panic("not implemented")
}

func (t *Timui[B]) Set(pos mathi.Vec2, char rune, fg, bg int32) {
	clip := t.PeekClip()
	if clip.Contains(pos) {
		t.front.set(pos, char, fg, bg)
	}
}

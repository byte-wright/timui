package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

func (t *Timui) Border(style [6]rune, fg, bg RGBColor) {
	c := t.CurrentArea()

	for x := c.From.X + 1; x < c.To.X-1; x++ {
		t.SetBorder(mathi.Vec2{X: x, Y: c.From.Y}, style[0], fg, bg)
		t.SetBorder(mathi.Vec2{X: x, Y: c.To.Y - 1}, style[0], fg, bg)
	}

	for y := c.From.Y + 1; y < c.To.Y-1; y++ {
		t.SetBorder(mathi.Vec2{X: c.From.X, Y: y}, style[1], fg, bg)
		t.Set(mathi.Vec2{X: c.To.X - 1, Y: y}, style[1], fg, bg)
	}

	t.SetBorder(c.From, style[2], fg, bg)
	t.SetBorder(mathi.Vec2{X: c.To.X - 1, Y: c.From.Y}, style[3], fg, bg)
	t.SetBorder(mathi.Vec2{X: c.From.X, Y: c.To.Y - 1}, style[4], fg, bg)
	t.SetBorder(mathi.Vec2{X: c.To.X - 1, Y: c.To.Y - 1}, style[5], fg, bg)
}

func (t *Timui) HLine(style [3]rune, fg, bg RGBColor) {
	c := t.CurrentArea()

	for x := c.From.X + 1; x < c.To.X-1; x++ {
		t.SetBorder(mathi.Vec2{X: x, Y: c.From.Y}, style[1], fg, bg)
	}

	t.SetBorder(c.From, style[0], fg, bg)
	t.SetBorder(mathi.Vec2{X: c.To.X - 1, Y: c.From.Y}, style[2], fg, bg)
}

func (t *Timui) VLine(style [3]rune, fg, bg RGBColor) {
	c := t.CurrentArea()

	for y := c.From.Y + 1; y < c.To.Y-1; y++ {
		t.SetBorder(mathi.Vec2{X: c.From.X, Y: y}, style[1], fg, bg)
	}

	t.SetBorder(c.From, style[0], fg, bg)
	t.SetBorder(mathi.Vec2{X: c.From.X, Y: c.To.Y - 1}, style[2], fg, bg)
}

func (t *Timui) SetAlpha(pos mathi.Vec2, char rune, fg, bg RGBAColor) {
	clip := t.PeekClip()
	if clip.Contains(pos) {
		t.front.Set(pos, char, uint32(fg), uint32(bg))
	}
}

func (t *Timui) Set(pos mathi.Vec2, char rune, fg, bg RGBColor) {
	clip := t.PeekClip()
	if clip.Contains(pos) {
		t.front.Set(pos, char, uint32(fg.RGBA(0xff)), uint32(bg.RGBA(0xff)))
	}
}

func (t *Timui) SetArea(char rune, fg, bg RGBColor) {
	clip := t.PeekClip()
	area := t.CurrentArea()

	for y := area.From.Y; y < area.To.Y; y++ {
		for x := area.From.X; x < area.To.X; x++ {
			p := mathi.Vec2{X: x, Y: y}
			if clip.Contains(p) {
				t.front.Set(p, char, uint32(fg.RGBA(0xff)), uint32(bg.RGBA(0xff)))
			}
		}
	}
}

func (t *Timui) SetAreaAlpha(char rune, fg, bg RGBAColor) {
	clip := t.PeekClip()

	area := t.CurrentArea()

	for y := area.From.Y; y < area.To.Y; y++ {
		for x := area.From.X; x < area.To.X; x++ {
			p := mathi.Vec2{X: x, Y: y}
			if clip.Contains(p) {
				t.front.Set(p, char, uint32(fg), uint32(bg))
			}
		}
	}
}

func (t *Timui) SetBorder(pos mathi.Vec2, char rune, fg, bg RGBColor) {
	if t.ClipContains(pos) {
		current := t.front.Get(pos).Char

		char = mapBorderRune(char, current)

		t.front.Set(pos, char, uint32(fg.RGBA(0xff)), uint32(bg.RGBA(0xff)))
	}
}

func mapBorderRune(char rune, current rune) rune {
	switch char {
	case '═':
		switch current {
		case '╩':
			char = '╩'
		}
	case '║':
		switch current {
		case '╣':
			char = '╣'
		case '┤':
			char = '╢'
		}
	case '├':
		switch current {
		case '║':
			char = '╟'
		}
	case '╠':
		switch current {
		case '│':
			char = '╞'
		}
	case '╣':
		switch current {
		case '│':
			char = '╡'
		}
	case '╦':
		switch current {
		case '╩':
			char = '╬'
		case '─':
			char = '╥'
		}
	case '╩':
		switch current {
		case '─':
			char = '╨'
		}
	case '-':
		switch current {
		case '+':
			char = '+'
		}

	}
	return char
}

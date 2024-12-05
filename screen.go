package timui

import "gitlab.com/bytewright/gmath/mathi"

type screen struct {
	size     mathi.Vec2
	gridSize mathi.Vec2
	chars    []cell
}

type cell struct {
	char rune
	fg   int32
	bg   int32
}

func newScreen(size mathi.Vec2) *screen {
	return &screen{
		size:     size,
		gridSize: size,
		chars:    make([]cell, size.X*size.Y),
	}
}

func (s *screen) get(pos mathi.Vec2) cell {
	if pos.X < 0 || pos.Y < 0 || pos.X >= s.size.X || pos.Y >= s.size.Y {
		return cell{}
	}

	i := pos.Y*s.gridSize.X + pos.X
	return s.chars[i]
}

func (s *screen) set(pos mathi.Vec2, char rune, fg, bg int32) {
	if pos.X < 0 || pos.Y < 0 || pos.X >= s.size.X || pos.Y >= s.size.Y {
		return
	}

	i := pos.Y*s.gridSize.X + pos.X
	s.chars[i].char = char
	s.chars[i].fg = fg
	s.chars[i].bg = bg
}

func (s *screen) clear(char rune, fg, bg int32) {
	for y := 0; y < s.size.Y; y++ {
		for x := 0; x < s.size.X; x++ {
			s.set(mathi.Vec2{X: x, Y: y}, char, fg, bg)
		}
	}
}

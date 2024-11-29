package timui

import "gitlab.com/bytewright/gmath/mathi"

type screen struct {
	size     mathi.Vec2
	gridSize mathi.Vec2
	chars    []rune
}

func newScreen(size mathi.Vec2) *screen {
	return &screen{
		size:     size,
		gridSize: size,
		chars:    make([]rune, size.X*size.Y),
	}
}

func (s *screen) get(pos mathi.Vec2) rune {
	if pos.X < 0 || pos.Y < 0 || pos.X >= s.size.X || pos.Y >= s.size.Y {
		return 0
	}

	i := pos.Y*s.gridSize.X + pos.X
	return s.chars[i]
}

func (s *screen) set(pos mathi.Vec2, char rune) {
	if pos.X < 0 || pos.Y < 0 || pos.X >= s.size.X || pos.Y >= s.size.Y {
		return
	}

	i := pos.Y*s.gridSize.X + pos.X
	s.chars[i] = char
}

func (s *screen) clear(char rune) {
	for y := 0; y < s.size.Y; y++ {
		for x := 0; x < s.size.X; x++ {
			s.set(mathi.Vec2{X: x, Y: y}, char)
		}
	}
}

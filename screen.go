package timui

import (
	"fmt"

	"gitlab.com/bytewright/gmath/mathi"
)

type screen struct {
	size     mathi.Vec2
	gridSize mathi.Vec2
	chars    []cell
}

type cell struct {
	// char defines the character that should be painted. If it is 0 it means that the
	// character at the screen is currently undefined and must be updated.
	char rune
	fg   RGBColor
	bg   RGBColor
}

func newScreen(size mathi.Vec2) *screen {
	return &screen{
		size:     size,
		gridSize: size,
		chars:    make([]cell, size.X*size.Y),
	}
}

func (s *screen) resize(size mathi.Vec2) {
	if s.size == size {
		return
	}

	if s.size.X >= size.X && s.size.Y >= size.Y {
		// if we shrink the size, just set new size
		s.size = size
		fmt.Println("resize shrink")
		return
	}

	if s.gridSize.X >= size.X && s.gridSize.Y >= size.Y {
		// if grid size is already larger

		if size.X > s.size.X {
			// extend width
			for y := 0; y < s.size.Y; y++ {
				for x := s.size.X; x < size.X; x++ {
					s.chars[y*s.gridSize.X+x].char = 0
				}
			}
			s.size.X = size.X
		}

		if size.Y > s.size.Y {
			for y := s.size.Y; y < size.Y; y++ {
				for x := 0; x < s.size.X; x++ {
					s.chars[y*s.gridSize.X+x].char = 0
				}
			}

			s.size.Y = size.Y
		}

		return
	}
	// brute force
	s.size = size
	s.gridSize = size
	s.chars = make([]cell, size.X*size.Y)
}

func (s *screen) get(pos mathi.Vec2) cell {
	if pos.X < 0 || pos.Y < 0 || pos.X >= s.size.X || pos.Y >= s.size.Y {
		return cell{}
	}

	i := pos.Y*s.gridSize.X + pos.X
	return s.chars[i]
}

// set paints the character, foreground and background
// if the character is 0 the character will not be updated, old one remains.
// colors are not changed when transparent.
func (s *screen) set(pos mathi.Vec2, char rune, fg, bg RGBAColor) {
	if pos.X < 0 || pos.Y < 0 || pos.X >= s.size.X || pos.Y >= s.size.Y {
		return
	}

	i := pos.Y*s.gridSize.X + pos.X

	if char != 0 {
		s.chars[i].char = char
	}

	s.chars[i].fg = s.chars[i].fg.Blend(fg)
	s.chars[i].bg = s.chars[i].bg.Blend(bg)
}

// clear sets the whole screen to given values.
func (s *screen) clear(char rune, fg, bg RGBColor) {
	for y := 0; y < s.size.Y; y++ {
		for x := 0; x < s.size.X; x++ {
			i := y*s.gridSize.X + x

			s.chars[i].char = char
			s.chars[i].fg = fg
			s.chars[i].bg = bg
		}
	}
}

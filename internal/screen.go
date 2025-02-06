package internal

import (
	"gitlab.com/bytewright/gmath/mathi"
)

type Screen struct {
	Size     mathi.Vec2
	gridSize mathi.Vec2
	chars    []cell
}

type cell struct {
	// Char defines the character that should be painted. If it is 0 it means that the
	// character at the screen is currently undefined and must be updated.
	Char rune
	FG   uint32
	BG   uint32
}

func NewScreen(size mathi.Vec2) *Screen {
	return &Screen{
		Size:     size,
		gridSize: size,
		chars:    make([]cell, size.X*size.Y),
	}
}

func (s *Screen) Resize(size mathi.Vec2) {
	if s.Size == size {
		return
	}

	if s.Size.X >= size.X && s.Size.Y >= size.Y {
		// if we shrink the size, just set new size
		s.Size = size
		return
	}

	if s.gridSize.X >= size.X && s.gridSize.Y >= size.Y {
		// if grid size is already larger

		if size.X > s.Size.X {
			// extend width
			for y := 0; y < s.Size.Y; y++ {
				for x := s.Size.X; x < size.X; x++ {
					s.chars[y*s.gridSize.X+x].Char = 0
				}
			}
			s.Size.X = size.X
		}

		if size.Y > s.Size.Y {
			for y := s.Size.Y; y < size.Y; y++ {
				for x := 0; x < s.Size.X; x++ {
					s.chars[y*s.gridSize.X+x].Char = 0
				}
			}

			s.Size.Y = size.Y
		}

		return
	}
	// brute force
	s.Size = size
	s.gridSize = size
	s.chars = make([]cell, size.X*size.Y)
}

func (s *Screen) Get(pos mathi.Vec2) cell {
	if pos.X < 0 || pos.Y < 0 || pos.X >= s.Size.X || pos.Y >= s.Size.Y {
		return cell{}
	}

	i := pos.Y*s.gridSize.X + pos.X
	return s.chars[i]
}

// set paints the character, foreground and background
// if the character is 0 the character will not be updated, old one remains.
// colors are blended when transparent.
func (s *Screen) Set(pos mathi.Vec2, char rune, fg, bg uint32) {
	if pos.X < 0 || pos.Y < 0 || pos.X >= s.Size.X || pos.Y >= s.Size.Y {
		return
	}

	i := pos.Y*s.gridSize.X + pos.X

	if char != 0 {
		s.chars[i].Char = char
	}

	s.chars[i].FG = blendColor(s.chars[i].FG, fg)
	s.chars[i].BG = blendColor(s.chars[i].BG, bg)
}

// clear sets the whole screen to given values.
func (s *Screen) SetScreen(char rune, fg, bg uint32) {
	for y := 0; y < s.Size.Y; y++ {
		for x := 0; x < s.Size.X; x++ {
			i := y*s.gridSize.X + x

			s.chars[i].Char = char
			s.chars[i].FG = fg
			s.chars[i].BG = bg
		}
	}
}

func blendColor(base uint32, top uint32) uint32 {
	tr := (top >> 16) & 0xff
	tg := (top >> 8) & 0xff
	tb := (top) & 0xff
	a := (top >> 24) & 0xff

	br := (base >> 16) & 0xff
	bg := (base >> 8) & 0xff
	bb := (base) & 0xff

	nr := (br*(0xff-a) + tr*a) >> 8
	ng := (bg*(0xff-a) + tg*a) >> 8
	nb := (bb*(0xff-a) + tb*a) >> 8

	return uint32(nr<<16 | ng<<8 | nb)
}

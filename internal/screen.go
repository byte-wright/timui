package internal

import (
	"gitlab.com/bytewright/gmath/mathi"
)

type Screen struct {
	Size  mathi.Vec2
	Chars []cell
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
		Size:  size,
		Chars: make([]cell, size.X*size.Y),
	}
}

func (s *Screen) Resize(size mathi.Vec2) {
	if s.Size == size {
		return
	}

	s.Size = size
	s.Chars = make([]cell, size.X*size.Y)
}

func (s *Screen) Get(pos mathi.Vec2) cell {
	if pos.X < 0 || pos.Y < 0 || pos.X >= s.Size.X || pos.Y >= s.Size.Y {
		return cell{}
	}

	i := pos.Y*s.Size.X + pos.X
	return s.Chars[i]
}

// set paints the character, foreground and background
// if the character is 0 the character will not be updated, old one remains.
// colors are blended when transparent.
func (s *Screen) Set(pos mathi.Vec2, char rune, fg, bg uint32) {
	if pos.X < 0 || pos.Y < 0 || pos.X >= s.Size.X || pos.Y >= s.Size.Y {
		return
	}

	i := pos.Y*s.Size.X + pos.X

	if char != 0 {
		s.Chars[i].Char = char
	}

	s.Chars[i].FG = blendColor(s.Chars[i].FG, fg)
	s.Chars[i].BG = blendColor(s.Chars[i].BG, bg)
}

// SetScreen sets the whole screen to given values.
func (s *Screen) SetScreen(char rune, fg, bg uint32) {
	cell := cell{Char: char, FG: fg, BG: bg}
	for i := range s.Chars {
		s.Chars[i] = cell
	}
}

func blendColor(base uint32, top uint32) uint32 {
	a := (top >> 24) & 0xff
	if a == 0 {
		return base
	}
	if a == 255 {
		return top & 0xffffff
	}

	tr := (top >> 16) & 0xff
	tg := (top >> 8) & 0xff
	tb := (top) & 0xff

	br := (base >> 16) & 0xff
	bg := (base >> 8) & 0xff
	bb := (base) & 0xff

	nr := (br*(0xff-a) + tr*a) >> 8
	ng := (bg*(0xff-a) + tg*a) >> 8
	nb := (bb*(0xff-a) + tb*a) >> 8

	return uint32(nr<<16 | ng<<8 | nb)
}

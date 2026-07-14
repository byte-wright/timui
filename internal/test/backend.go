package test

import (
	"strings"

	"github.com/byte-wright/timui"
	"gitlab.com/bytewright/gmath/mathi"
)

// TestBackend accumulates the pushed cell diffs like a real terminal, so
// after any Finish its grid equals that frame's front buffer.
type TestBackend struct {
	size    mathi.Vec2
	chars   [][]rune
	Mouse   mathi.Vec2
	Pressed bool
}

func NewBackend(x, y int) *TestBackend {
	chars := make([][]rune, y)
	for i := range chars {
		chars[i] = make([]rune, x)
		for j := range chars[i] {
			chars[i][j] = ' '
		}
	}

	return &TestBackend{
		size:  mathi.Vec2{X: x, Y: y},
		chars: chars,
		Mouse: mathi.Vec2{X: -1, Y: -1},
	}
}

func (b *TestBackend) Size() mathi.Vec2                { return b.size }
func (b *TestBackend) MousePosition() mathi.Vec2       { return b.Mouse }
func (b *TestBackend) MousePressed(key timui.Key) bool { return b.Pressed }
func (b *TestBackend) Render()                         {}

func (b *TestBackend) Set(pos mathi.Vec2, char rune, fg, bg uint32) {
	if char != 0 {
		b.chars[pos.Y][pos.X] = char
	}
}

// String frames the screen so snapshot files keep their trailing spaces even
// through editors that trim whitespace.
func (b *TestBackend) String() string {
	border := "+" + strings.Repeat("-", b.size.X) + "+\n"

	sb := strings.Builder{}
	sb.WriteString(border)
	for _, row := range b.chars {
		sb.WriteString("|" + string(row) + "|\n")
	}
	sb.WriteString(border)

	return sb.String()
}

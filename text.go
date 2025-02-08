package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

func (t *Timui) Text(name string, pos mathi.Vec2, fg, bg RGBAColor) {
	p := t.CurrentArea().From.Add(pos)

	for i, r := range []rune(name) {
		t.front.Set(p.Add(mathi.Vec2{X: i}), r, uint32(fg), uint32(bg))
	}
}

func (t *Timui) Label(name string) {
	w := t.CurrentArea().Size().X

	if len(name) > w+1 {
		ne := w - 2
		if ne > 0 {
			name = name[:ne] + "..."
		}
	}

	t.Text(name, mathi.Vec2{}, Transparent, Transparent)
	t.moveCursor(mathi.Vec2{Y: 1})
}

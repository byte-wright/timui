package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

func (t *Timui[B]) Text(name string, pos mathi.Vec2, fg, bg RGBAColor) {
	p := t.CurrentArea().From.Add(pos)

	for i, r := range []rune(name) {
		t.front.Set(p.Add(mathi.Vec2{X: i}), r, uint32(fg), uint32(bg))
	}
}

func (t *Timui[B]) Label(name string) {
	t.Text(name, mathi.Vec2{}, 0xffdddddd, 0x00000000)
	t.moveCursor(mathi.Vec2{Y: 1})
}

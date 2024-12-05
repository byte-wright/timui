package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

func (t *Timui[B]) Text(name string, pos mathi.Vec2, fg, bg int32) {
	p := t.CurrentArea().From.Add(pos)

	for i, r := range []rune(name) {
		t.backend.Set(p.Add(mathi.Vec2{X: i}), r, fg, bg)
	}
}

func (t *Timui[B]) Label(name string) {
	t.Text(name, mathi.Vec2{}, 0xff9999, 0x000000)
	t.moveCursor(mathi.Vec2{Y: 1})
}

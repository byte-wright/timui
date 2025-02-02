package timui

import (
	"strings"

	"gitlab.com/bytewright/gmath/mathi"
)

func (g *Timui[B]) Checkbox(name string, checked *bool) bool {
	size := g.CurrentArea().Size()

	size.Y = 1

	mouse := g.MouseInput(name, mathi.Box2{To: size})

	bgCol := RGBA(0x11, 0x11, 0x77, 0xff)

	if mouse.Hovered() > 0 {
		bgCol = RGBA(0x33, 0x33, 0x99, 0xff)
	}

	if mouse.LeftPressed() > 0 {
		bgCol = RGBA(0x00, 0x00, 0x66, 0xff)
	}

	if *checked {
		g.Text("[X] ", mathi.Vec2{}, 0xffbbbbbb, bgCol)
	} else {
		g.Text("[ ] ", mathi.Vec2{}, 0xffbbbbbb, bgCol)
	}

	pad := size.X - len(name) - 4

	g.Text(name+strings.Repeat(" ", pad), mathi.Vec2{X: 4}, 0xffbbbbbb, bgCol)

	g.moveCursor(mathi.Vec2{Y: 1})

	clicked := mouse.LeftReleased()

	if clicked {
		*checked = !*checked
	}

	return clicked
}

package timui

import (
	"strings"

	"gitlab.com/bytewright/gmath/mathi"
)

func (g *Timui[B]) Button(name string) bool {
	size := g.CurrentArea().Size()

	size.Y = 1

	mouse := g.MouseInput(name, mathi.Box2{To: size})

	bgCol := RGB(0x11, 0x11, 0x77)

	if mouse.Hovered() > 0 {
		bgCol = RGB(0x33, 0x33, 0x99)
	}

	if mouse.LeftPressed() > 0 {
		bgCol = RGB(0x00, 0x00, 0x66)
	}

	pad := (size.X - len(name))
	padl := pad / 2
	padr := pad - padl
	padr -= 1
	padl -= 1

	if padr < 0 {
		padr = 0
	}

	if padl < 0 {
		padl = 0
	}

	g.Text("["+strings.Repeat(" ", padl)+name+strings.Repeat(" ", padr)+"]",
		mathi.Vec2{}, 0xbbbbbb, bgCol)

	g.moveCursor(mathi.Vec2{Y: 1})

	return mouse.LeftReleased()
}

package timui

import (
	"strings"

	"gitlab.com/bytewright/gmath/mathi"
)

func (g *Timui[B]) Button(name string) bool {
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
		mathi.Vec2{}, 0xffbbbbbb, bgCol)

	g.moveCursor(mathi.Vec2{Y: 1})

	return mouse.LeftReleased()
}

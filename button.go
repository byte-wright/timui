package timui

import (
	"strings"

	"gitlab.com/bytewright/gmath/mathi"
)

func (g *Timui[B]) Button(name string) bool {
	size := g.CurrentArea().Size()

	size.Y = 1

	mouse := g.MouseInput(name, mathi.Box2{To: size})

	bgCol := int32(0x111177)

	if mouse.Hovered() > 0 {
		bgCol = 0x333399
	}

	if mouse.LeftPressed() > 0 {
		bgCol = 0x000066
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

package timui

import (
	"strings"

	"gitlab.com/bytewright/gmath/mathi"
)

func (g *Timui[B]) Button(name string) bool {
	size := g.CurrentArea().Size()

	size.Y = 2

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
	padl -= 1
	padr -= 2

	g.Text(strings.Repeat(" ", padl)+name+strings.Repeat(" ", padr), mathi.Vec2{X: 1}, 0xbbbbbb, bgCol)

	pos := g.CurrentArea().From

	for x := 2; x < size.X-1; x++ {
		g.front.set(pos.Add(mathi.Vec2{X: x, Y: 1}), '▀', 0x333333, 0x000000)
	}

	g.front.set(pos.Add(mathi.Vec2{X: size.X - 2}), '▄', 0x333333, 0x000000)

	g.moveCursor(mathi.Vec2{Y: 2})

	return mouse.LeftReleased()
}

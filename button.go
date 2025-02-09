package timui

import (
	"strings"

	"gitlab.com/bytewright/gmath/mathi"
)

func (g *Timui) Button(name string) bool {
	area := g.CurrentArea()
	size := area.Size()

	size.Y = 1

	mouse := g.MouseInputForSize(name, size)

	bgCol := g.Theme.Widget.BG

	if mouse.Hovered() > 0 {
		bgCol = g.Theme.Widget.HoverBG
	}

	if mouse.LeftPressed() > 0 {
		bgCol = g.Theme.Widget.InteractBG
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

	g.Set(area.From, '[', g.Theme.Widget.Line, bgCol)

	g.Set(area.From.Add(mathi.Vec2{X: size.X - 1}), ']', g.Theme.Widget.Line, bgCol)

	g.Text(strings.Repeat(" ", padl)+name+strings.Repeat(" ", padr),
		mathi.Vec2{X: 1}, g.Theme.Widget.Text.RGBA(0xff), bgCol.RGBA(0xff))

	g.moveCursor(mathi.Vec2{Y: 1})

	return mouse.LeftReleased()
}

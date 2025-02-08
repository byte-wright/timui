package timui

import (
	"strings"

	"gitlab.com/bytewright/gmath/mathi"
)

func (g *Timui) Checkbox(name string, checked *bool) bool {
	size := g.CurrentArea().Size()

	size.Y = 1

	mouse := g.MouseInput(name, mathi.Box2{To: size})

	bgCol := g.Theme.Widget.BG

	if mouse.Hovered() > 0 {
		bgCol = g.Theme.Widget.HoverBG
	}

	if mouse.LeftPressed() > 0 {
		bgCol = g.Theme.Widget.InteractBG
	}

	bgColA := bgCol.RGBA(0xff)

	if *checked {
		g.Text("[X] ", mathi.Vec2{}, g.Theme.Widget.Line.RGBA(0xff), bgColA)
	} else {
		g.Text("[ ] ", mathi.Vec2{}, g.Theme.Widget.Line.RGBA(0xff), bgColA)
	}

	pad := size.X - len(name) - 4

	g.Text(name+strings.Repeat(" ", pad), mathi.Vec2{X: 4}, g.Theme.Widget.Text.RGBA(0xff), bgColA)

	g.moveCursor(mathi.Vec2{Y: 1})

	clicked := mouse.LeftReleased()

	if clicked {
		*checked = !*checked
	}

	return clicked
}

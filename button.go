package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

func (g *Timui[B]) Button(name string) bool {
	size := g.CurrentArea().Size()

	size.Y = 2

	mouse := g.MouseInput(name, mathi.Box2{To: size})

	if mouse.LeftPressed() > 0 {
	}

	pad := (size.X - len(name)) / 2

	g.Text(name, mathi.Vec2{X: pad})

	g.moveCursor(mathi.Vec2{Y: 2})

	return mouse.LeftReleased()
}

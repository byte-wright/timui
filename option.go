package timui

import (
	"strings"

	"gitlab.com/bytewright/gmath/mathi"
)

type OptionGroupElement[V comparable] struct {
	t        *Timui
	selected *V
}

func OptionGroup[V comparable](t *Timui, name string, selected *V, f func(*OptionGroupElement[V])) {
	t.id.Push(name)

	f(&OptionGroupElement[V]{
		t:        t,
		selected: selected,
	})

	t.id.Pop()
}

func (o OptionGroupElement[V]) Option(name string, value V) bool {
	size := o.t.CurrentArea().Size()

	size.Y = 1

	mouse := o.t.MouseInputForSize(name, size)

	bgCol := o.t.Theme.Widget.BG

	if mouse.Hovered() > 0 {
		bgCol = o.t.Theme.Widget.HoverBG
	}

	if mouse.LeftPressed() > 0 {
		bgCol = o.t.Theme.Widget.InteractBG
	}

	bgColA := bgCol.RGBA(0xff)

	if *o.selected == value {
		o.t.Text("(X) ", mathi.Vec2{}, o.t.Theme.Widget.Line.RGBA(0xff), bgColA)
	} else {
		o.t.Text("( ) ", mathi.Vec2{}, o.t.Theme.Widget.Line.RGBA(0xff), bgColA)
	}

	pad := size.X - len(name) - 4

	o.t.Text(name+strings.Repeat(" ", pad), mathi.Vec2{X: 4}, o.t.Theme.Widget.Text.RGBA(0xff), bgColA)

	o.t.moveCursor(mathi.Vec2{Y: 1})

	clicked := mouse.LeftReleased()

	if clicked {
		*o.selected = value
	}

	return clicked
}

package timui

import (
	"strconv"

	"gitlab.com/bytewright/gmath/mathi"
)

var (
	dropdownInputStyle     = [3]rune{'[', ' ', ']'}
	dropdownBordeStyle     = [6]rune{'-', '|', '/', '\\', '\\', '/'}
	dropdownBordeStyleLine = [6]rune{'─', '│', '┌', '┐', '└', '┘'}
)

type dropdown[B Backend] struct {
	g        *Timui[B]
	open     bool
	elements int
	selected int
	paint    func(i int, selected bool)
}

type dropdownManager[B Backend] struct {
	lastDropdowns map[ID]*dropdown[B]
	nextDropdowns map[ID]*dropdown[B]
}

func (g *Timui[B]) Dropdown(id string, elements int, selected *int, paint func(i int, s bool)) {
	cid := g.PushID(id)
	dd, has := g.dropdownManager.lastDropdowns[cid]
	if !has {
		dd = &dropdown[B]{g: g}
	}

	dd.elements = elements
	dd.paint = paint
	dd.selected = *selected

	g.dropdownManager.nextDropdowns[cid] = dd

	dd.paintSelection()

	g.PopID()

	area := *g.CurrentArea()

	if dd.open {
		g.runAfter(func() {
			g.PushID(cid.value)

			modal := g.MouseInput("modal", *g.CurrentArea())
			if modal.LeftReleased() {
				dd.open = false
			}

			height := dd.elements

			area.To.Y = area.From.Y + height + 2
			g.PushArea(area)

			g.Border(dropdownBordeStyleLine)

			pad := g.Pad(1, 1, 1, 1)

			for i := 0; i < dd.elements; i++ {
				ma := *g.CurrentArea()

				ma.To.Y = ma.From.Y + 1
				g.PushArea(ma)

				mi := g.MouseInput(strconv.Itoa(i), mathi.Box2{To: ma.Size()})

				if mi.Hovered() > 0 {
					fg := RGB(0xff, 0xff, 0xff)
					bg := RGB(0x00, 0x00, 0x66)
					g.HLine(dropdownInputStyle, fg, bg)
				}

				if mi.LeftReleased() {
					*selected = i
					dd.open = false
				}

				g.PopArea()

				ma.From.X += 8
				ma.To.X -= 8

				g.PushArea(ma)
				dd.paint(i, dd.selected == i)
				g.PopArea()

				g.moveCursor(mathi.Vec2{Y: 1})
			}

			pad.Finish()
			g.PopArea()
			g.PopID()
		})
	}
}

func (d *dropdown[B]) paintSelection() {
	remWidth := d.g.CurrentArea().Size().X

	size := mathi.Vec2{X: remWidth, Y: 1}

	mouse := d.g.MouseInput("area", mathi.Box2{To: size})

	if mouse.Hovered() > 0 {
		fg := RGB(0xff, 0xff, 0xff)
		bg := RGB(0x00, 0x00, 0x88)

		d.g.HLine(dropdownInputStyle, fg, bg)
	} else {
		fg := RGB(0xff, 0xff, 0xff)
		bg := RGB(0x00, 0x00, 0x66)
		d.g.HLine(dropdownInputStyle, fg, bg)
	}

	if mouse.LeftReleased() {
		d.open = !d.open
	}

	if d.open {
		d.g.Text("[v]", mathi.Vec2{X: remWidth - 3}, RGB(255, 255, 255), RGB(0x33, 0x33, 0x33))
	} else {
		d.g.Text("[^]", mathi.Vec2{X: remWidth - 3}, RGB(255, 255, 255), RGB(0x33, 0x33, 0x33))
	}

	a := *d.g.CurrentArea()
	a.To.Y = a.From.Y
	a.From.X += 1
	a.To.X -= 1
	d.g.PushArea(a)
	d.paint(d.selected, false)
	d.g.PopArea()
	d.g.moveCursor(mathi.Vec2{Y: 1})
}

func newDropdownManager[B Backend]() *dropdownManager[B] {
	return &dropdownManager[B]{
		lastDropdowns: map[ID]*dropdown[B]{},
		nextDropdowns: map[ID]*dropdown[B]{},
	}
}

func (m *dropdownManager[B]) finish(_ *Timui[B]) {
	m.lastDropdowns, m.nextDropdowns = m.nextDropdowns, m.lastDropdowns
	m.nextDropdowns = map[ID]*dropdown[B]{}
}

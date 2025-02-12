package timui

import (
	"strconv"

	"github.com/byte-wright/timui/internal"
	"gitlab.com/bytewright/gmath/mathi"
)

var (
	dropdownInputStyle     = [3]rune{'[', ' ', ']'}
	dropdownSelectionStyle = [3]rune{' ', ' ', ' '}
	dropdownBordeStyle     = [6]rune{'-', '|', '/', '\\', '\\', '/'}
	dropdownBordeStyleLine = [6]rune{'─', '│', '┌', '┐', '└', '┘'}
)

type dropdown struct {
	g        *Timui
	open     bool
	elements int
	selected int
	paint    func(i int, selected bool)
}

type dropdownManager struct {
	lastDropdowns map[internal.ID]*dropdown
	nextDropdowns map[internal.ID]*dropdown
}

func (g *Timui) Dropdown(id string, elements int, selected *int, paint func(i int, s bool)) {
	cid := g.id.Push(id)
	dd, has := g.dropdownManager.lastDropdowns[cid]
	if !has {
		dd = &dropdown{g: g}
	}

	dd.elements = elements
	dd.paint = paint
	dd.selected = *selected

	g.dropdownManager.nextDropdowns[cid] = dd

	dd.paintSelection(&g.Theme)

	g.id.Pop()

	area := *g.CurrentArea()

	if dd.open {
		g.runAfter(func() {
			g.id.PushID(cid)

			modal := g.MouseInput("modal")
			if modal.LeftReleased() {
				dd.open = false
			}

			height := dd.elements

			area.To.Y = area.From.Y + height + 2
			g.PushArea(area)

			g.Border(dropdownBordeStyleLine, g.Theme.Widget.Line, g.Theme.Widget.BG)

			g.Pad(1, 1, 1, 1, func() {
				g.SetArea(' ', g.Theme.Widget.Text, g.Theme.Widget.BG)

				g.id.Push("selection")

				for i := 0; i < dd.elements; i++ {
					ma := *g.CurrentArea()

					ma.To.Y = ma.From.Y + 1
					g.PushArea(ma)

					mi := g.MouseInputForSize(strconv.Itoa(i), ma.Size())

					if mi.Hovered() > 0 || i == *selected {
						g.HLine(dropdownSelectionStyle, g.Theme.Widget.Text, g.Theme.Widget.HoverBG)
					}

					if mi.LeftReleased() {
						*selected = i
						dd.open = false
					}

					g.PopArea()

					ma.From.X += 1
					ma.To.X -= 1

					g.PushArea(ma)
					dd.paint(i, dd.selected == i)
					g.PopArea()

					g.moveCursor(mathi.Vec2{Y: 1})
				}

				g.id.Pop()
			})
			g.PopArea()
			g.id.Pop()
		})
	}
}

func (d *dropdown) paintSelection(theme *Theme) {
	remWidth := d.g.CurrentArea().Size().X

	size := mathi.Vec2{X: remWidth, Y: 1}

	mouse := d.g.MouseInputForSize("area", size)
	if mouse.LeftReleased() {
		d.open = !d.open
	}

	bgCol := theme.Widget.BG

	if mouse.Hovered() > 0 || d.open {
		bgCol = theme.Widget.HoverBG
	}

	if mouse.LeftPressed() > 0 {
		bgCol = theme.Widget.InteractBG
	}

	d.g.HLine(dropdownInputStyle, theme.Widget.Line, bgCol)

	if d.open {
		d.g.Text("][ʌ]", mathi.Vec2{X: remWidth - 4}, theme.Widget.Line.RGBA(0xff), bgCol.RGBA(0xff))
	} else {
		d.g.Text("][v]", mathi.Vec2{X: remWidth - 4}, theme.Widget.Line.RGBA(0xff), bgCol.RGBA(0xff))
	}

	a := *d.g.CurrentArea()
	a.To.Y = a.From.Y
	a.From.X += 2
	a.To.X -= 5
	d.g.PushArea(a)
	d.g.SetArea(0, theme.Widget.Text, bgCol)
	d.paint(d.selected, false)
	d.g.PopArea()
	d.g.moveCursor(mathi.Vec2{Y: 1})
}

func newDropdownManager() *dropdownManager {
	return &dropdownManager{
		lastDropdowns: map[internal.ID]*dropdown{},
		nextDropdowns: map[internal.ID]*dropdown{},
	}
}

func (m *dropdownManager) finish(_ *Timui) {
	m.lastDropdowns, m.nextDropdowns = m.nextDropdowns, m.lastDropdowns
	m.nextDropdowns = map[internal.ID]*dropdown{}
}

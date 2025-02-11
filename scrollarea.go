package timui

import (
	"github.com/byte-wright/timui/internal"
	"gitlab.com/bytewright/gmath/mathi"
)

type scrollAreaManager struct {
	lastScrollAreas map[internal.ID]*scrollAreaV
	nextScrollAreas map[internal.ID]*scrollAreaV
}

type scrollAreaV struct {
	g             *Timui
	factor        float32
	contentHeight int
}

func (g *Timui) ScrollAreaV(id string, body func()) {
	uid := g.id.Push(id)

	scrollArea := g.scrollAreaManager.lastScrollAreas[uid]
	if scrollArea == nil {
		scrollArea = &scrollAreaV{
			g: g,
		}
	}

	g.scrollAreaManager.nextScrollAreas[uid] = scrollArea

	area := *g.CurrentArea()
	area.To.X -= 1

	g.PushArea(area)
	g.PushClip(area)
	overflow := scrollArea.contentHeight - area.Size().Y
	if overflow < 0 {
		overflow = 0
	}
	g.PushAreaTranslation(mathi.Vec2{Y: -int(float32(overflow) * scrollArea.factor)})

	topY := g.CurrentArea().From.Y

	body()

	scrollArea.contentHeight = g.CurrentArea().From.Y - topY

	barWidth := 1

	g.PopArea()
	g.PopClip()
	g.PopArea()

	areaHeight := g.CurrentArea().Size().Y
	minHeight := 1

	factor := float32(areaHeight) / float32(scrollArea.contentHeight)
	if factor > 1 {
		factor = 1
	}
	knobHeight := minHeight + int(float32(areaHeight-minHeight-2)*factor)

	// draw bar background
	barArea := *g.CurrentArea()
	barArea.From.X = barArea.To.X - barWidth
	size := barArea.Size()

	g.PushArea(barArea)
	g.SetArea(' ', MustRGBS("#bbb"), MustRGBS("#013"))

	g.Set(barArea.From, '▴', MustRGBS("#bbb"), MustRGBS("#135"))
	g.Set(mathi.Vec2{X: barArea.From.X, Y: barArea.To.Y - 1}, '▾', MustRGBS("#bbb"), MustRGBS("#135"))

	g.PopArea()

	// draw knob
	knobTop := int(float32(size.Y-knobHeight-2) * scrollArea.factor)
	knobArea := *g.CurrentArea()
	knobArea.From.X = knobArea.To.X - barWidth
	knobArea.From.Y = knobArea.From.Y + 1 + knobTop
	knobArea.To.Y = knobArea.From.Y + knobHeight

	if knobHeight < size.Y { // why?
		g.PushArea(knobArea)
		g.SetArea('░', MustRGBS("#46f"), MustRGBS("#00a"))
		g.PopArea()

		g.PushArea(barArea)

		dragAreaSize := g.CurrentArea().Size()

		dragArea := mathi.Box2{From: mathi.Vec2{Y: 1}, To: mathi.Vec2{X: barWidth, Y: dragAreaSize.Y - 1}}

		pos := mathi.Vec2{Y: knobTop}

		g.Draggable("draggable", dragArea, mathi.Vec2{X: barWidth, Y: knobHeight}, &pos)

		scrollArea.factor = 0
		if size.Y-knobHeight-2 > 0 {
			scrollArea.factor = float32(pos.Y) / float32(size.Y-knobHeight-2)
		}
		g.PopArea()
	}

	g.id.Pop()
}

func newScrollAreaManager() *scrollAreaManager {
	return &scrollAreaManager{
		lastScrollAreas: map[internal.ID]*scrollAreaV{},
		nextScrollAreas: map[internal.ID]*scrollAreaV{},
	}
}

func (m *scrollAreaManager) finish() {
	m.lastScrollAreas, m.nextScrollAreas = m.nextScrollAreas, m.lastScrollAreas
	m.nextScrollAreas = map[internal.ID]*scrollAreaV{}
}

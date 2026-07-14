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

	overflow := scrollArea.contentHeight - area.Size().Y
	if overflow < 0 {
		overflow = 0
	}

	g.WithArea(area, func() {
		g.WithClip(area, func() {
			g.WithAreaTranslation(mathi.Vec2{Y: -int(float32(overflow) * scrollArea.factor)}, func() {
				topY := g.CurrentArea().From.Y

				body()

				scrollArea.contentHeight = g.CurrentArea().From.Y - topY
			})
		})
	})

	barWidth := 1

	areaHeight := g.CurrentArea().Size().Y
	minHeight := 1

	factor := float32(areaHeight) / float32(scrollArea.contentHeight)
	if factor > 1 {
		factor = 1
	}

	trackHeight := areaHeight - 2
	knobHeight := minHeight + int(float32(trackHeight-minHeight)*factor)
	knobRange := trackHeight - knobHeight

	barArea := *g.CurrentArea()
	barArea.From.X = barArea.To.X - barWidth

	bar := g.Theme.ScrollBar

	g.WithArea(barArea, func() {
		g.SetArea(' ', bar.Text, bar.BG)

		g.Set(barArea.From, '▴', bar.Text, bar.ArrowBG)
		g.Set(mathi.Vec2{X: barArea.From.X, Y: barArea.To.Y - 1}, '▾', bar.Text, bar.ArrowBG)
	})

	knobTop := int(float32(knobRange) * scrollArea.factor)
	knobArea := barArea
	knobArea.From.Y = knobArea.From.Y + 1 + knobTop
	knobArea.To.Y = knobArea.From.Y + knobHeight

	if knobHeight > 0 && knobHeight <= trackHeight {
		g.WithArea(knobArea, func() {
			g.SetArea('░', bar.Knob, bar.KnobBG)
		})

		g.WithArea(barArea, func() {
			dragArea := mathi.Box2{From: mathi.Vec2{Y: 1}, To: mathi.Vec2{X: barWidth, Y: areaHeight - 1}}

			pos := mathi.Vec2{Y: knobTop}

			g.Draggable("draggable", dragArea, mathi.Vec2{X: barWidth, Y: knobHeight}, &pos)

			scrollArea.factor = 0
			if knobRange > 0 {
				scrollArea.factor = float32(pos.Y) / float32(knobRange)
			}
		})
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

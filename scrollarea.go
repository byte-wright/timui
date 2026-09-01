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

// wheelStep is how many rows one wheel notch moves.
const wheelStep = 3

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

	g.WithArea(area, func() {
		wheel := g.MouseInput("wheel").Wheel()
		if wheel.Y != 0 && overflow > 0 {
			scrollArea.factor = clampFactor(scrollArea.factor +
				float32(wheel.Y*wheelStep)/float32(overflow))
		}
	})

	scrollArea.factor, _ = g.scrollBarV(*g.CurrentArea(), scrollArea.contentHeight, scrollArea.factor)

	g.id.Pop()
}

// scrollBarV draws the vertical bar along the right edge of area and returns the
// position after any drag of its knob, 0 at the top and 1 at the bottom, plus
// whether the knob was actually moved. The position is only recomputed from the
// knob while it is dragged, because the knob sits on whole rows and feeding that
// back every frame would snap away any finer movement, like a wheel notch. The
// caller must have pushed an id scope, the bar claims a fixed id inside it.
func (g *Timui) scrollBarV(area mathi.Box2, contentHeight int, factor float32) (float32, bool) {
	const barWidth = 1
	const minKnobHeight = 1

	areaHeight := area.Size().Y

	visible := float32(1)
	if contentHeight > 0 {
		visible = float32(areaHeight) / float32(contentHeight)
	}

	if visible > 1 {
		visible = 1
	}

	trackHeight := areaHeight - 2
	knobHeight := minKnobHeight + int(float32(trackHeight-minKnobHeight)*visible)
	knobRange := trackHeight - knobHeight

	barArea := area
	barArea.From.X = barArea.To.X - barWidth

	bar := g.Theme.ScrollBar

	g.WithArea(barArea, func() {
		g.SetArea(' ', bar.Text, bar.BG)

		g.Set(barArea.From, '▴', bar.Text, bar.ArrowBG)
		g.Set(mathi.Vec2{X: barArea.From.X, Y: barArea.To.Y - 1}, '▾', bar.Text, bar.ArrowBG)
	})

	knobTop := int(float32(knobRange) * factor)
	knobArea := barArea
	knobArea.From.Y = knobArea.From.Y + 1 + knobTop
	knobArea.To.Y = knobArea.From.Y + knobHeight

	if knobHeight <= 0 || knobHeight > trackHeight {
		return factor, false
	}

	g.WithArea(knobArea, func() {
		g.SetArea('░', bar.Knob, bar.KnobBG)
	})

	moved := false

	g.WithArea(barArea, func() {
		dragArea := mathi.Box2{From: mathi.Vec2{Y: 1}, To: mathi.Vec2{X: barWidth, Y: areaHeight - 1}}

		pos := mathi.Vec2{Y: knobTop}

		g.Draggable("draggable", dragArea, mathi.Vec2{X: barWidth, Y: knobHeight}, &pos)

		if pos.Y == knobTop || knobRange <= 0 {
			return
		}

		factor = float32(pos.Y) / float32(knobRange)
		moved = true
	})

	return factor, moved
}

func clampFactor(f float32) float32 {
	if f < 0 {
		return 0
	}

	if f > 1 {
		return 1
	}

	return f
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

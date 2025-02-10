package timui

import (
	"github.com/byte-wright/timui/internal"
	"gitlab.com/bytewright/gmath/mathi"
)

type scrollAreaManager struct {
	lastScrollAreas map[internal.ID]*ScrollAreaV
	nextScrollAreas map[internal.ID]*ScrollAreaV
}

type ScrollAreaV struct {
	g             *Timui
	factor        float32
	contentHeight int
}

// type ScrollStyle[T Texture] struct {
// 	BackgroundTextures XRectStyle[T]
// 	KnobTextures       XRectStyle[T]
// 	Pad                int
// 	Width              int
// }

func (g *Timui) ScrollAreaV(id string) *ScrollAreaV {
	uid := g.id.Push(id)

	scrollArea := g.scrollAreaManager.lastScrollAreas[uid]
	if scrollArea == nil {
		scrollArea = &ScrollAreaV{
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

	scrollArea.contentHeight = g.CurrentArea().From.Y

	return scrollArea
}

func (s *ScrollAreaV) Finish() {
	s.contentHeight = s.g.CurrentArea().From.Y - s.contentHeight

	// st := s.g.Style.Scroll.Peek()
	barWidth := 1

	s.g.PopArea()
	s.g.PopClip()
	s.g.PopArea()

	areaHeight := s.g.CurrentArea().Size().Y
	minHeight := 1

	factor := float32(areaHeight-minHeight) / float32(s.contentHeight)
	knobHeight := minHeight + int(float32(areaHeight-minHeight)*factor)

	// draw bar background
	barArea := *s.g.CurrentArea()
	barArea.From.X = barArea.To.X - barWidth
	size := barArea.Size()

	s.g.PushArea(barArea)
	//	s.g.VRect(&st.BackgroundTextures, size, st.Pad, st.Pad, v41)
	s.g.SetArea('▴', MustRGBS("#bbb"), MustRGBS("#135")) // ▾

	s.g.PopArea()

	// draw knob
	// knob := s.g.Texture("scroll.png")

	knobTop := int(float32(size.Y-knobHeight) * s.factor)
	knobArea := *s.g.CurrentArea()
	knobArea.From.X = knobArea.To.X - barWidth
	knobArea.From.Y = knobArea.From.Y + knobTop
	knobArea.To.Y = knobArea.From.Y + knobHeight

	if knobHeight < size.Y {
		s.g.PushArea(knobArea)
		// s.g.VRect(&st.KnobTextures, mathi.Vec2{X: barWidth, Y: knobHeight}, st.Pad, st.Pad, v41)
		s.g.SetArea('░', MustRGBS("#4f4"), MustRGBS("#0a0"))
		s.g.PopArea()

		s.g.PushArea(barArea)

		dragAreaSize := s.g.CurrentArea().Size()

		dragArea := mathi.Box2{To: mathi.Vec2{X: barWidth, Y: dragAreaSize.Y}}

		pos := mathi.Vec2{Y: knobTop}

		s.g.Draggable("draggable", dragArea, mathi.Vec2{X: barWidth, Y: knobHeight}, &pos)

		s.factor = float32(pos.Y) / float32(size.Y-knobHeight)
		s.g.PopArea()
	}

	s.g.id.Pop()
}

func newScrollAreaManager() *scrollAreaManager {
	return &scrollAreaManager{
		lastScrollAreas: map[internal.ID]*ScrollAreaV{},
		nextScrollAreas: map[internal.ID]*ScrollAreaV{},
	}
}

func (m *scrollAreaManager) finish() {
	m.lastScrollAreas, m.nextScrollAreas = m.nextScrollAreas, m.lastScrollAreas
	m.nextScrollAreas = map[internal.ID]*ScrollAreaV{}
}

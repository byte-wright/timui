package timui

import (
	"github.com/byte-wright/timui/internal"
	"gitlab.com/bytewright/gmath/mathi"
)

type Draggable struct {
	dragEnd bool
}

type draggableManager struct {
	lastDraggables map[internal.ID]*Draggable
	nextDraggables map[internal.ID]*Draggable

	dragging      *Draggable
	startPos      mathi.Vec2
	mouseStartPos mathi.Vec2
}

func (t *Timui) Draggable(id string, area mathi.Box2, size mathi.Vec2, pos *mathi.Vec2) (bool, bool) {
	cid := t.id.Push(id)

	draggable, has := t.draggableManager.lastDraggables[cid]
	if !has {
		draggable = &Draggable{}
	}

	cursor := t.CurrentArea().From

	mouse := t.MouseInputForArea(id, mathi.NewBox2FromSize((area.From.Add(*pos)), size))

	dragStart := false

	if mouse.LeftPressed() > 0 && t.draggableManager.dragging == nil {
		t.draggableManager.dragging = draggable

		t.draggableManager.startPos = cursor.Add(area.From).Add(*pos)
		t.draggableManager.mouseStartPos = t.backend.MousePosition()

		dragStart = true
	}

	if t.draggableManager.dragging == draggable {
		globalPos := t.draggableManager.startPos.Add(t.backend.MousePosition().Sub(t.draggableManager.mouseStartPos))
		globalPos = mkInside(area.Translate(cursor), globalPos, size)
		*pos = globalPos.Sub(cursor).Sub(area.From)
	}

	t.draggableManager.nextDraggables[cid] = draggable

	t.id.Pop()

	dragEnd := draggable.dragEnd

	draggable.dragEnd = false
	return dragStart, dragEnd
}

func mkInside(area mathi.Box2, pos, size mathi.Vec2) mathi.Vec2 {
	if pos.X+size.X > area.To.X {
		pos.X = area.To.X - size.X
	}

	if pos.Y+size.Y > area.To.Y {
		pos.Y = area.To.Y - size.Y
	}

	if pos.X < area.From.X {
		pos.X = area.From.X
	}

	if pos.Y < area.From.Y {
		pos.Y = area.From.Y
	}

	return pos
}

func newDraggableManager() *draggableManager {
	return &draggableManager{
		lastDraggables: map[internal.ID]*Draggable{},
		nextDraggables: map[internal.ID]*Draggable{},
	}
}

func (m *draggableManager) finish(g *Timui) {
	if m.dragging != nil && !g.backend.MousePressed(MouseButtonLeft) {
		m.dragging.dragEnd = true
		m.dragging = nil
	}

	m.lastDraggables, m.nextDraggables = m.nextDraggables, m.lastDraggables
	m.nextDraggables = map[internal.ID]*Draggable{}
}

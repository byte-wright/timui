package timui

import (
	"github.com/byte-wright/timui/internal"
	"gitlab.com/bytewright/gmath/mathi"
)

type MouseInput struct {
	area       mathi.Box2
	hovered    bool
	hoverCount int

	leftPressed    bool
	leftPressCount int
	leftReleased   bool
}

type mouseInputManager struct {
	lastInputs map[internal.ID]*MouseInput
	nextInputs map[internal.ID]*MouseInput

	underCursor     *MouseInput
	lastUnderCursor *MouseInput
}

// MouseInputForSize creates a mouse area from given size relative to the current cursor area
func (t *Timui) MouseInputForSize(id string, size mathi.Vec2) *MouseInput {
	return t.MouseInputForArea(id, mathi.Box2{To: size})
}

// MouseInput creates a mouse area for the current cursor area
func (t *Timui) MouseInput(id string) *MouseInput {
	s := t.CurrentArea().Size()
	return t.MouseInputForArea(id, mathi.Box2{To: s})
}

// MouseInputForArea creates a mouse area from given area relative to the current cursor area
func (t *Timui) MouseInputForArea(id string, area mathi.Box2) *MouseInput {
	cid := t.id.Push(id)
	mouseInput, has := t.mouseInputManager.lastInputs[cid]
	if !has {
		mouseInput = &MouseInput{}
	}

	// update hover state
	if mouseInput.hovered {
		mouseInput.hoverCount++
		mouseInput.hovered = false
	} else {
		mouseInput.hoverCount = 0
	}

	// update leftPressed state
	mouseInput.leftReleased = false
	if mouseInput.leftPressed {
		mouseInput.leftPressCount++
		mouseInput.leftPressed = false
	} else {
		if mouseInput.leftPressCount > 0 && mouseInput.hoverCount > 0 {
			mouseInput.leftReleased = true
		}
		mouseInput.leftPressCount = 0
	}

	mouseInput.area = area.Translate(t.CurrentArea().From)
	mousePos := t.backend.MousePosition()
	if mouseInput.area.Contains(mousePos) && t.PeekClip().Contains(mousePos) {
		t.mouseInputManager.underCursor = mouseInput
	}

	t.mouseInputManager.nextInputs[cid] = mouseInput

	t.id.Pop()

	return mouseInput
}

func (m *MouseInput) Hovered() int {
	return m.hoverCount
}

func (m *MouseInput) LeftPressed() int {
	return m.leftPressCount
}

func (m *MouseInput) LeftReleased() bool {
	return m.leftReleased
}

func newMouseInputManager() *mouseInputManager {
	return &mouseInputManager{
		lastInputs: map[internal.ID]*MouseInput{},
		nextInputs: map[internal.ID]*MouseInput{},
	}
}

func (m *mouseInputManager) finish(g *Timui) {
	m.lastUnderCursor = m.underCursor

	if m.underCursor != nil {
		m.underCursor.hovered = true

		if g.backend.MousePressed(MouseButtonLeft) {
			m.underCursor.leftPressed = true
		}

		m.underCursor = nil
	}

	m.lastInputs, m.nextInputs = m.nextInputs, m.lastInputs
	m.nextInputs = map[internal.ID]*MouseInput{}
}

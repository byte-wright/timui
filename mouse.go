package timui

import (
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

type mouseInputManager[B Backend] struct {
	lastInputs map[ID]*MouseInput
	nextInputs map[ID]*MouseInput

	underCursor     *MouseInput
	lastUnderCursor *MouseInput
}

func (g *Timui[B]) MouseInput(id string, area mathi.Box2) *MouseInput {
	cid := g.PushID(id)
	mouseInput, has := g.mouseInputManager.lastInputs[cid]
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

	mouseInput.area = area.Translate(g.CurrentArea().From)
	mousePos := g.backend.MousePosition()
	if mouseInput.area.Contains(mousePos) && g.PeekClip().Contains(mousePos) {
		g.mouseInputManager.underCursor = mouseInput
	}

	g.mouseInputManager.nextInputs[cid] = mouseInput

	g.PopID()

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

func newMouseInputManager[B Backend]() *mouseInputManager[B] {
	return &mouseInputManager[B]{
		lastInputs: map[ID]*MouseInput{},
		nextInputs: map[ID]*MouseInput{},
	}
}

func (m *mouseInputManager[B]) finish(g *Timui[B]) {
	m.lastUnderCursor = m.underCursor

	if m.underCursor != nil {
		m.underCursor.hovered = true

		if g.backend.MousePressed(MouseButtonLeft) {
			m.underCursor.leftPressed = true
		}

		m.underCursor = nil
	}

	m.lastInputs, m.nextInputs = m.nextInputs, m.lastInputs
	m.nextInputs = map[ID]*MouseInput{}
}

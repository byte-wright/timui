package timui

import (
	"strings"

	"gitlab.com/bytewright/gmath/mathi"
)

// selectionManager owns the one text selection that can exist at a time - there
// is no focus model, so like draggableManager this is global rather than keyed
// per widget. Its rendering is handled by the library (draw), callers never draw
// it themselves.
type selectionManager struct {
	dragging bool
	active   bool

	rect    mathi.Box2
	anchor  mathi.Vec2
	current mathi.Vec2
}

func newSelectionManager() *selectionManager {
	return &selectionManager{}
}

// update runs once per frame after every widget (including deferred overlays)
// has had a chance to register a MouseInput and a CopyRect for this frame's
// mouse position. A press that lands where nothing claims the click starts a
// selection scoped to the innermost CopyRect under it; releasing the button
// copies a non-empty selection to the clipboard.
func (m *selectionManager) update(t *Timui) {
	mousePos := t.backend.MousePosition()
	pressed := t.backend.MousePressed(MouseButtonLeft)
	rect := t.copyRectManager.take(t)

	switch {
	case m.dragging:
		m.current = mousePos

		if !pressed {
			m.dragging = false

			if text := m.text(t); text != "" {
				t.backend.WriteClipboard(text)
			} else {
				m.active = false
			}
		}

	case pressed && !t.mouseInputManager.selectionBlocked():
		m.dragging = true
		m.active = true
		m.rect = rect
		m.anchor = mousePos
		m.current = mousePos

	case pressed:
		m.active = false
	}
}

// box is the normalized selection, clamped to the rect it was started in. The
// clamp can never empty it out: the drag anchor is only ever set to a position
// PushCopyRect already found inside rect.
func (m *selectionManager) box() mathi.Box2 {
	b := mathi.Box2{
		From: m.anchor.Min(m.current),
		To:   m.anchor.Max(m.current).Add(mathi.Vec2{X: 1, Y: 1}),
	}

	b.From = b.From.Max(m.rect.From)
	b.To = b.To.Min(m.rect.To)

	return b
}

// text reads the selected cells back out of the front buffer, trimming
// trailing spaces off each row.
func (m *selectionManager) text(t *Timui) string {
	box := m.box()

	lines := make([]string, 0, box.To.Y-box.From.Y)

	for y := box.From.Y; y < box.To.Y; y++ {
		var line strings.Builder

		for x := box.From.X; x < box.To.X; x++ {
			r := t.front.GetChar(mathi.Vec2{X: x, Y: y})
			if r == 0 {
				r = ' '
			}

			line.WriteRune(r)
		}

		lines = append(lines, strings.TrimRight(line.String(), " "))
	}

	return strings.Join(lines, "\n")
}

// draw paints the highlight over the active selection. It runs last, after
// every widget and overlay for the frame, so the highlight always sits on top.
func (m *selectionManager) draw(t *Timui) {
	if !m.active {
		return
	}

	box := m.box()
	fg := t.Theme.Selection.Text.RGBA(0xff)
	bg := t.Theme.Selection.BG.RGBA(0xff)

	for y := box.From.Y; y < box.To.Y; y++ {
		for x := box.From.X; x < box.To.X; x++ {
			t.front.Set(mathi.Vec2{X: x, Y: y}, 0, uint32(fg), uint32(bg))
		}
	}
}

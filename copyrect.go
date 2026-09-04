package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

// copyRectManager tracks the stack of copy rects pushed this frame and, as they
// are pushed, the innermost one the mouse is currently inside. Nested pushes
// happen later than their parent's in call order, and deferred overlays
// (dialog, dropdown) push even later than the main pass, so "last pushed match
// wins" gives the same z-order priority as mouseInputManager.underCursor.
type copyRectManager struct {
	rects []mathi.Box2

	current    mathi.Box2
	hasCurrent bool
}

func newCopyRectManager() *copyRectManager {
	return &copyRectManager{}
}

// PushCopyRect scopes text selection to area: a drag started inside it can only
// select within it. Grid/Rows/Columns cells, ScrollAreaV, ScrollListV and Dialog
// each push one around their content, so a selection can't run into a
// neighboring cell or off-screen content clipped away.
func (t *Timui) PushCopyRect(area mathi.Box2) {
	t.copyRectManager.rects = append(t.copyRectManager.rects, area)

	if area.Contains(t.backend.MousePosition()) {
		t.copyRectManager.current = area
		t.copyRectManager.hasCurrent = true
	}
}

func (t *Timui) PopCopyRect() {
	t.copyRectManager.rects = t.copyRectManager.rects[:len(t.copyRectManager.rects)-1]
}

func (t *Timui) PeekCopyRect() mathi.Box2 {
	if len(t.copyRectManager.rects) > 0 {
		return t.copyRectManager.rects[len(t.copyRectManager.rects)-1]
	}

	return mathi.Box2{To: t.backend.Size()}
}

// WithCopyRect runs body with the selection scoped to area.
func (t *Timui) WithCopyRect(area mathi.Box2, body func()) {
	t.PushCopyRect(area)
	body()
	t.PopCopyRect()
}

// take returns the innermost copy rect the mouse was inside this frame, or the
// whole screen if it wasn't inside any, and resets tracking for the next frame.
func (m *copyRectManager) take(t *Timui) mathi.Box2 {
	rect, has := m.current, m.hasCurrent

	m.current = mathi.Box2{}
	m.hasCurrent = false

	if !has {
		return mathi.Box2{To: t.backend.Size()}
	}

	return rect
}

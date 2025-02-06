package timui

import (
	"github.com/byte-wright/timui/internal"
	"gitlab.com/bytewright/gmath/mathi"
)

type Timui[B Backend] struct {
	backend B
	front   *internal.Screen
	back    *internal.Screen

	area []mathi.Box2

	after []func()

	id                internal.IDManager
	clipManager       clipManager
	mouseInputManager mouseInputManager[B]
	dropdownManager   dropdownManager[B]

	Theme Theme
}

type Key int

var (
	MouseButtonLeft  Key = 1_000_000
	MouseButtonRight Key = 1_000_001
)

type Backend interface {
	Size() mathi.Vec2
	MousePosition() mathi.Vec2
	MousePressed(key Key) bool
	Set(pos mathi.Vec2, char rune, fg, bg uint32)
	Render()
}

func New[B Backend](backend B) *Timui[B] {
	size := backend.Size()

	front := internal.NewScreen(size)
	front.SetScreen(' ', 0, 0)

	tui := &Timui[B]{
		backend: backend,
		front:   front,
		back:    internal.NewScreen(size),

		area:              []mathi.Box2{},
		id:                *internal.NewIDManager(),
		clipManager:       *newClipManager(),
		mouseInputManager: *newMouseInputManager[B](),
		dropdownManager:   *newDropdownManager[B](),

		Theme: DefaultTheme,
	}

	tui.finish()

	return tui
}

func (t *Timui[B]) runAfter(f func()) {
	t.after = append(t.after, f)
}

func (t *Timui[B]) Finish() {
	for _, f := range t.after {
		f()
	}

	for y := 0; y < t.front.Size.Y; y++ {
		for x := 0; x < t.front.Size.X; x++ {
			pos := mathi.Vec2{X: x, Y: y}

			if t.front.Get(pos) != t.back.Get(pos) {
				fc := t.front.Get(pos)
				t.backend.Set(pos, fc.Char, uint32(fc.FG), uint32(fc.BG))
			}
		}
	}

	t.backend.Render()

	// todo screen resizing
	t.back, t.front = t.front, t.back

	t.back.Resize(t.backend.Size())
	t.front.Resize(t.backend.Size())

	t.front.SetScreen(' ', uint32(t.Theme.Text), uint32(t.Theme.BG))

	t.finish()
}

// finish is called after each pass and also once before the first pass
func (t *Timui[B]) finish() {
	t.area = t.area[:0]
	t.area = append(t.area, mathi.Box2{To: t.backend.Size()})

	t.mouseInputManager.finish(t)
	t.dropdownManager.finish(t)

	t.after = t.after[:0]
}

func (t *Timui[B]) CurrentArea() *mathi.Box2 {
	return &t.area[len(t.area)-1]
}

func (t *Timui[B]) PushArea(area mathi.Box2) {
	t.area = append(t.area, area)
}

func (t *Timui[B]) PopArea() {
	t.area = t.area[:len(t.area)-1]
}

func (t *Timui[B]) GetMousePosition() mathi.Vec2 {
	return t.backend.MousePosition()
}

func (t *Timui[B]) moveCursor(delta mathi.Vec2) {
	t.CurrentArea().From = t.CurrentArea().From.Add(delta)
}

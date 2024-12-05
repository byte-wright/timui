package timui

import "gitlab.com/bytewright/gmath/mathi"

type Timui[B Backend] struct {
	backend B
	front   *screen
	back    *screen

	area []mathi.Box2

	idManager         idManager
	clipManager       clipManager
	mouseInputManager mouseInputManager[B]
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
	Set(pos mathi.Vec2, char rune, fg, bg int32)
	Render()
}

func New[B Backend](backend B) *Timui[B] {
	size := backend.Size()

	front := newScreen(size)
	front.clear(' ', 0, 0)

	tui := &Timui[B]{
		backend: backend,
		front:   front,
		back:    newScreen(size),

		area:              []mathi.Box2{},
		idManager:         *newIDManager(),
		clipManager:       *newClipManager(),
		mouseInputManager: *newMouseInputManager[B](),
	}

	tui.reset()

	return tui
}

func (t *Timui[B]) Finish() {
	for y := 0; y < t.front.size.Y; y++ {
		for x := 0; x < t.front.size.X; x++ {
			pos := mathi.Vec2{X: x, Y: y}

			if t.front.get(pos) != t.back.get(pos) {
				fc := t.front.get(pos)
				t.backend.Set(pos, fc.char, fc.fg, fc.bg)
			}
		}
	}

	t.backend.Render()

	// todo screen resizing
	t.back = t.front
	t.front.clear(' ', 0, 0)

	t.reset()
}

// reset is called after each pass and also once before the first pass
func (t *Timui[B]) reset() {
	t.area = t.area[:]
	t.area = append(t.area, mathi.Box2{To: t.backend.Size()})

	t.mouseInputManager.finish(t)
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

func (t *Timui[B]) moveCursor(delta mathi.Vec2) {
	t.CurrentArea().From = t.CurrentArea().From.Add(delta)
}

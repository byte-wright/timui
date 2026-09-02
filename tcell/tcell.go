package tcell

import (
	"github.com/byte-wright/timui"
	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
	"gitlab.com/bytewright/gmath/mathi"
)

type TCellBackend struct {
	screen     tcell.Screen
	mousePos   mathi.Vec2
	leftMouse  bool
	rightMouse bool
	input      timui.Input
	mods       timui.ModMask
}

func NewBackend() (*TCellBackend, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	err = screen.Init()
	if err != nil {
		return nil, err
	}

	screen.EnableMouse(tcell.MouseMotionEvents)

	defStyle := tcell.StyleDefault.
		Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)

	screen.SetStyle(defStyle)

	return &TCellBackend{
		screen: screen,
	}, nil
}

func (b *TCellBackend) Set(pos mathi.Vec2, char rune, fg, bg uint32) {
	st := tcell.StyleDefault.Foreground(tcell.NewHexColor(int32(fg))).Background(tcell.NewHexColor(int32(bg)))
	b.screen.SetContent(pos.X, pos.Y, char, nil, st)
}

func (b *TCellBackend) Render() {
	b.screen.ShowCursor(b.mousePos.X, b.mousePos.Y)
	b.screen.Show()
}

func (b *TCellBackend) MousePosition() mathi.Vec2 {
	return b.mousePos
}

func (b *TCellBackend) TakeInput() timui.Input {
	input := b.input
	b.input = timui.Input{}
	input.Mods = b.mods

	return input
}

func (b *TCellBackend) MousePressed(key timui.Key) bool {
	if key == timui.MouseButtonLeft {
		return b.leftMouse
	}

	if key == timui.MouseButtonRight {
		return b.rightMouse
	}

	return false
}

func (b *TCellBackend) Events() bool {
	exit := false

	for b.screen.HasPendingEvent() {
		ev := b.screen.PollEvent()

		switch evt := ev.(type) {
		case *tcell.EventKey:
			b.mods = mapMods(evt.Modifiers())

			if evt.Key() == tcell.KeyCtrlC {
				exit = true
			}

			key, ok := mapKey(evt)
			if ok {
				b.input.Keys = append(b.input.Keys, key)
			}

		case *tcell.EventMouse:
			b.mods = mapMods(evt.Modifiers())

			x, y := evt.Position()
			b.mousePos = mathi.Vec2{X: x, Y: y}

			bmsk := evt.Buttons()

			b.leftMouse = bmsk&tcell.Button1 == tcell.Button1
			b.rightMouse = bmsk&tcell.Button2 == tcell.Button2

			b.input.Wheel = b.input.Wheel.Add(wheelDelta(bmsk))
		}
	}

	return exit
}

// mapMods translates tcell's modifier mask. Only the modifiers timui knows are
// kept, the rest, like meta, are dropped.
func mapMods(mask tcell.ModMask) timui.ModMask {
	mods := timui.ModMask(0)

	if mask&tcell.ModShift != 0 {
		mods |= timui.ModShift
	}

	if mask&tcell.ModCtrl != 0 {
		mods |= timui.ModCtrl
	}

	if mask&tcell.ModAlt != 0 {
		mods |= timui.ModAlt
	}

	return mods
}

// wheelDelta reads the wheel bits out of a button mask. tcell reports wheel
// movement as momentary buttons rather than as its own event.
func wheelDelta(mask tcell.ButtonMask) mathi.Vec2 {
	d := mathi.Vec2{}

	if mask&tcell.WheelUp != 0 {
		d.Y--
	}

	if mask&tcell.WheelDown != 0 {
		d.Y++
	}

	if mask&tcell.WheelLeft != 0 {
		d.X--
	}

	if mask&tcell.WheelRight != 0 {
		d.X++
	}

	return d
}

var keyMap = map[tcell.Key]timui.Key{
	tcell.KeyUp:         timui.KeyUp,
	tcell.KeyDown:       timui.KeyDown,
	tcell.KeyLeft:       timui.KeyLeft,
	tcell.KeyRight:      timui.KeyRight,
	tcell.KeyPgUp:       timui.KeyPageUp,
	tcell.KeyPgDn:       timui.KeyPageDown,
	tcell.KeyHome:       timui.KeyHome,
	tcell.KeyEnd:        timui.KeyEnd,
	tcell.KeyEnter:      timui.KeyEnter,
	tcell.KeyEscape:     timui.KeyEscape,
	tcell.KeyTab:        timui.KeyTab,
	tcell.KeyBackspace:  timui.KeyBackspace,
	tcell.KeyBackspace2: timui.KeyBackspace,
	tcell.KeyDelete:     timui.KeyDelete,
}

func mapKey(evt *tcell.EventKey) (timui.KeyEvent, bool) {
	if evt.Key() == tcell.KeyRune {
		return timui.KeyEvent{Key: timui.KeyRune, Rune: evt.Rune()}, true
	}

	key, ok := keyMap[evt.Key()]
	if !ok {
		return timui.KeyEvent{}, false
	}

	return timui.KeyEvent{Key: key}, true
}

func (b *TCellBackend) Size() mathi.Vec2 {
	x, y := b.screen.Size()
	return mathi.Vec2{X: x, Y: y}
}

func (b *TCellBackend) Exit() {
	b.screen.Fini()
}

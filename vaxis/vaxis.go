package vaxis

import (
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/byte-wright/timui"
	"gitlab.com/bytewright/gmath/mathi"
	"go.rockorager.dev/vaxis"
)

type VaxisBackend struct {
	vx       *vaxis.Vaxis
	window   vaxis.Window
	mousePos mathi.Vec2

	leftMouse  bool
	rightMouse bool
	input      timui.Input
	mods       timui.ModMask
}

func NewBackend() (*VaxisBackend, error) {
	// vaxis renders truecolor with the strict T.416 colon form (38:2:r:g:b)
	// by default. Terminals without subparameter support read the first value
	// after the "2" as the colorspace id, shifting every channel. tcell uses
	// the legacy 38;2;r;g;b form, which every truecolor terminal accepts, so
	// make vaxis emit that too. Respect an explicit override.
	if os.Getenv("VAXIS_FORCE_LEGACY_SGR") == "" {
		os.Setenv("VAXIS_FORCE_LEGACY_SGR", "1")
	}

	vx, err := vaxis.New(vaxis.Options{})
	if err != nil {
		return nil, err
	}

	return &VaxisBackend{
		vx:     vx,
		window: vx.Window(),
	}, nil
}

func (b *VaxisBackend) Set(pos mathi.Vec2, char rune, fg, bg uint32) {
	grapheme := " "
	if char != 0 {
		grapheme = string(char)
	}

	b.window.SetCell(pos.X, pos.Y, vaxis.Cell{
		Character: vaxis.Character{Grapheme: grapheme, Width: 1},
		Style: vaxis.Style{
			Foreground: vaxis.RGBColor(uint8(fg>>16), uint8(fg>>8), uint8(fg)),
			Background: vaxis.RGBColor(uint8(bg>>16), uint8(bg>>8), uint8(bg)),
		},
	})
}

func (b *VaxisBackend) Render() {
	b.vx.ShowCursor(b.mousePos.X, b.mousePos.Y, vaxis.CursorBlock)
	b.vx.Render()
}

func (b *VaxisBackend) MousePosition() mathi.Vec2 {
	return b.mousePos
}

func (b *VaxisBackend) TakeInput() timui.Input {
	input := b.input
	b.input = timui.Input{}
	input.Mods = b.mods

	return input
}

func (b *VaxisBackend) MousePressed(key timui.Key) bool {
	if key == timui.MouseButtonLeft {
		return b.leftMouse
	}

	if key == timui.MouseButtonRight {
		return b.rightMouse
	}

	return false
}

func (b *VaxisBackend) Events() bool {
	exit := false

	for {
		select {
		case ev := <-b.vx.Events():
			switch ev := ev.(type) {
			case vaxis.Key:
				b.applyMods(ev)

				if ev.Keycode == 'c' && ev.Modifiers&vaxis.ModCtrl != 0 {
					exit = true
				}

				key, ok := mapKey(ev)
				if ok {
					b.input.Keys = append(b.input.Keys, key)
				}

			case vaxis.Mouse:
				b.mods = mapMods(ev.Modifiers)

				b.mousePos = mathi.Vec2{X: ev.Col, Y: ev.Row}

				switch ev.EventType {
				case vaxis.EventPress:
					switch ev.Button {
					case vaxis.MouseLeftButton:
						b.leftMouse = true
					case vaxis.MouseRightButton:
						b.rightMouse = true
					case vaxis.MouseWheelUp:
						b.input.Wheel.Y--
					case vaxis.MouseWheelDown:
						b.input.Wheel.Y++
					case vaxis.MouseWheelLeft:
						b.input.Wheel.X--
					case vaxis.MouseWheelRight:
						b.input.Wheel.X++
					}
				case vaxis.EventRelease:
					switch ev.Button {
					case vaxis.MouseLeftButton:
						b.leftMouse = false
					case vaxis.MouseRightButton:
						b.rightMouse = false
					}
				}

			case vaxis.Resize:
				b.vx.Resize(ev)
			}
		default:
			return exit
		}
	}
}

// applyMods syncs the held modifier state with a key event. The kitty keyboard
// protocol reports modifier keys as their own press/release events, which is
// how a lone Ctrl press is seen; each such event only touches its own bit.
// Regular events report the whole set at once via ev.Modifiers.
func (b *VaxisBackend) applyMods(ev vaxis.Key) {
	if mod, ok := modifierKeys[ev.Keycode]; ok {
		if ev.EventType == vaxis.EventRelease {
			b.mods &^= mod
		} else {
			b.mods |= mod
		}
		return
	}

	b.mods = mapMods(ev.Modifiers)
}

// mapMods translates vaxis' modifier mask. Only the modifiers timui knows are
// kept, the rest, like super and meta, are dropped.
func mapMods(mask vaxis.ModifierMask) timui.ModMask {
	mods := timui.ModMask(0)

	if mask&vaxis.ModShift != 0 {
		mods |= timui.ModShift
	}

	if mask&vaxis.ModCtrl != 0 {
		mods |= timui.ModCtrl
	}

	if mask&vaxis.ModAlt != 0 {
		mods |= timui.ModAlt
	}

	return mods
}

var keyMap = map[rune]timui.Key{
	vaxis.KeyUp:        timui.KeyUp,
	vaxis.KeyDown:      timui.KeyDown,
	vaxis.KeyLeft:      timui.KeyLeft,
	vaxis.KeyRight:     timui.KeyRight,
	vaxis.KeyPgUp:      timui.KeyPageUp,
	vaxis.KeyPgDown:    timui.KeyPageDown,
	vaxis.KeyHome:      timui.KeyHome,
	vaxis.KeyEnd:       timui.KeyEnd,
	vaxis.KeyEnter:     timui.KeyEnter,
	vaxis.KeyEsc:       timui.KeyEscape,
	vaxis.KeyTab:       timui.KeyTab,
	vaxis.KeyBackspace: timui.KeyBackspace,
	vaxis.KeyDelete:    timui.KeyDelete,
}

// modifierKeys maps the kitty keyboard protocol's standalone modifier key
// events onto the modifiers timui knows.
var modifierKeys = map[rune]timui.ModMask{
	vaxis.KeyLeftShift:    timui.ModShift,
	vaxis.KeyRightShift:   timui.ModShift,
	vaxis.KeyLeftControl:  timui.ModCtrl,
	vaxis.KeyRightControl: timui.ModCtrl,
	vaxis.KeyLeftAlt:      timui.ModAlt,
	vaxis.KeyRightAlt:     timui.ModAlt,
}

// mapKey translates a vaxis key press. Special keys are matched against the
// known key names, printable keys keep the character that was actually typed
// (vaxis lowercases the keycode, the shifted character lives in Text).
func mapKey(k vaxis.Key) (timui.KeyEvent, bool) {
	if key, ok := keyMap[k.Keycode]; ok {
		return timui.KeyEvent{Key: key}, true
	}

	if k.Keycode > unicode.MaxRune {
		return timui.KeyEvent{}, false
	}

	r := k.Keycode
	if k.Text != "" {
		r, _ = utf8.DecodeRuneInString(k.Text)
	}

	return timui.KeyEvent{Key: timui.KeyRune, Rune: r}, true
}

func (b *VaxisBackend) Size() mathi.Vec2 {
	size := b.vx.Size()
	b.window = b.vx.Window()

	return mathi.Vec2{X: size.Cols, Y: size.Rows}
}

func (b *VaxisBackend) Exit() {
	b.vx.Close()
}

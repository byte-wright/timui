package vaxis

import (
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui"
	"go.rockorager.dev/vaxis"
)

func TestMapKey(t *testing.T) {
	tests := []struct {
		name string
		in   vaxis.Key
		out  timui.KeyEvent
		ok   bool
	}{
		{"up", vaxis.Key{Keycode: vaxis.KeyUp}, timui.KeyEvent{Key: timui.KeyUp}, true},
		{"down", vaxis.Key{Keycode: vaxis.KeyDown}, timui.KeyEvent{Key: timui.KeyDown}, true},
		{"left", vaxis.Key{Keycode: vaxis.KeyLeft}, timui.KeyEvent{Key: timui.KeyLeft}, true},
		{"right", vaxis.Key{Keycode: vaxis.KeyRight}, timui.KeyEvent{Key: timui.KeyRight}, true},
		{"page up", vaxis.Key{Keycode: vaxis.KeyPgUp}, timui.KeyEvent{Key: timui.KeyPageUp}, true},
		{"page down", vaxis.Key{Keycode: vaxis.KeyPgDown}, timui.KeyEvent{Key: timui.KeyPageDown}, true},
		{"home", vaxis.Key{Keycode: vaxis.KeyHome}, timui.KeyEvent{Key: timui.KeyHome}, true},
		{"end", vaxis.Key{Keycode: vaxis.KeyEnd}, timui.KeyEvent{Key: timui.KeyEnd}, true},
		{"enter", vaxis.Key{Keycode: vaxis.KeyEnter}, timui.KeyEvent{Key: timui.KeyEnter}, true},
		{"escape", vaxis.Key{Keycode: vaxis.KeyEsc}, timui.KeyEvent{Key: timui.KeyEscape}, true},
		{"tab", vaxis.Key{Keycode: vaxis.KeyTab}, timui.KeyEvent{Key: timui.KeyTab}, true},
		{"backspace", vaxis.Key{Keycode: vaxis.KeyBackspace}, timui.KeyEvent{Key: timui.KeyBackspace}, true},
		{"delete", vaxis.Key{Keycode: vaxis.KeyDelete}, timui.KeyEvent{Key: timui.KeyDelete}, true},
		{"rune", vaxis.Key{Keycode: 'a', Text: "a"}, timui.KeyEvent{Key: timui.KeyRune, Rune: 'a'}, true},
		{"shifted rune", vaxis.Key{Keycode: 'a', Text: "A"}, timui.KeyEvent{Key: timui.KeyRune, Rune: 'A'}, true},
		{"unknown special", vaxis.Key{Keycode: vaxis.KeyF01}, timui.KeyEvent{}, false},
	}

	for _, tc := range tests {
		out, ok := mapKey(tc.in)
		expect.Value(t, tc.name, ok).ToBe(tc.ok)
		expect.Value(t, tc.name+" key", out.Key).ToBe(tc.out.Key)
		expect.Value(t, tc.name+" rune", out.Rune).ToBe(tc.out.Rune)
	}
}

func TestMapMods(t *testing.T) {
	mods := mapMods(vaxis.ModShift | vaxis.ModCtrl | vaxis.ModAlt | vaxis.ModSuper)
	expect.Value(t, "known kept", mods).ToBe(timui.ModShift | timui.ModCtrl | timui.ModAlt)
	expect.Value(t, "shift kept", mods.Has(timui.ModShift)).ToBe(true)
	expect.Value(t, "ctrl kept", mods.Has(timui.ModCtrl)).ToBe(true)
	expect.Value(t, "alt kept", mods.Has(timui.ModAlt)).ToBe(true)
}

func TestMapModsEmpty(t *testing.T) {
	expect.Value(t, "none", mapMods(0)).ToBe(timui.ModMask(0))
}

func TestModifierPressRelease(t *testing.T) {
	be := &VaxisBackend{}

	// a lone Ctrl press is seen as its own event
	be.applyMods(vaxis.Key{Keycode: vaxis.KeyLeftControl, EventType: vaxis.EventPress})
	expect.Value(t, "ctrl held", be.mods.Has(timui.ModCtrl)).ToBe(true)

	// holding another modifier with ctrl keeps both
	be.applyMods(vaxis.Key{Keycode: vaxis.KeyRightShift, EventType: vaxis.EventPress})
	expect.Value(t, "ctrl+shift held", be.mods.Has(timui.ModCtrl|timui.ModShift)).ToBe(true)

	// regular keys report the full modifier set
	be.applyMods(vaxis.Key{Keycode: 'a', Text: "a", Modifiers: vaxis.ModCtrl | vaxis.ModShift})
	expect.Value(t, "still ctrl+shift", be.mods.Has(timui.ModCtrl|timui.ModShift)).ToBe(true)

	// releasing just ctrl keeps shift
	be.applyMods(vaxis.Key{Keycode: vaxis.KeyLeftControl, EventType: vaxis.EventRelease})
	expect.Value(t, "ctrl released", be.mods.Has(timui.ModCtrl)).ToBe(false)
	expect.Value(t, "shift kept", be.mods.Has(timui.ModShift)).ToBe(true)

	be.applyMods(vaxis.Key{Keycode: vaxis.KeyRightShift, EventType: vaxis.EventRelease})
	expect.Value(t, "shift released", be.mods).ToBe(timui.ModMask(0))
}

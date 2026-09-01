package timui_test

import (
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/internal/test"
)

func TestKeysArriveOnTheNextFrame(t *testing.T) {
	tui, be := test.New(t, 20, 5)

	be.Input.Keys = []timui.KeyEvent{
		{Key: timui.KeyPageUp},
		{Key: timui.KeyRune, Rune: 'f'},
	}

	tui.Finish()

	expect.Value(t, "page up", tui.KeyPressed(timui.KeyPageUp)).ToBe(true)
	expect.Value(t, "f", tui.RunePressed('f')).ToBe(true)
	expect.Value(t, "home", tui.KeyPressed(timui.KeyHome)).ToBe(false)

	tui.Finish()

	expect.Value(t, "keys are not repeated", len(tui.Keys())).ToBe(0)
}

package tcell

import (
	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
	"gitlab.com/bytewright/gmath/mathi"
)

type TCellBackend struct {
	screen   tcell.Screen
	mousePos mathi.Vec2
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

func (b *TCellBackend) Set(pos mathi.Vec2, char rune) {
	b.screen.SetContent(pos.X, pos.Y, char, nil, tcell.StyleDefault)
}

func (b *TCellBackend) Render() {
	b.screen.ShowCursor(b.mousePos.X, b.mousePos.Y)
	b.screen.Show()
}

func (b *TCellBackend) Events() bool {
	exit := false

	for b.screen.HasPendingEvent() {
		ev := b.screen.PollEvent()

		switch evt := ev.(type) {
		case *tcell.EventKey:
			if evt.Key() == tcell.KeyEscape || evt.Key() == tcell.KeyCtrlC {
				exit = true
			}

		case *tcell.EventMouse:
			x, y := evt.Position()
			b.mousePos = mathi.Vec2{X: x, Y: y}
		}
	}

	return exit
}

func (b *TCellBackend) Size() mathi.Vec2 {
	x, y := b.screen.Size()
	return mathi.Vec2{X: x, Y: y}
}

func (b *TCellBackend) Exit() {
	b.screen.Fini()
}

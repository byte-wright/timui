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
			if evt.Key() == tcell.KeyEscape || evt.Key() == tcell.KeyCtrlC {
				exit = true
			}

		case *tcell.EventMouse:
			x, y := evt.Position()
			b.mousePos = mathi.Vec2{X: x, Y: y}

			bmsk := evt.Buttons()

			b.leftMouse = bmsk&tcell.Button1 == tcell.Button1
			b.rightMouse = bmsk&tcell.Button2 == tcell.Button2
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

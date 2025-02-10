package timui

import (
	"strings"

	"gitlab.com/bytewright/gmath/mathi"
)

func (t *Timui) Text(name string, pos mathi.Vec2, fg, bg RGBAColor) {
	p := t.CurrentArea().From.Add(pos)

	for i, r := range []rune(name) {
		pos := p.Add(mathi.Vec2{X: i})
		t.SetAlpha(pos, r, fg, bg)
	}
}

func (t *Timui) Label(name string) {
	w := t.CurrentArea().Size().X

	t.Text(cutText(name, w), mathi.Vec2{}, Transparent, Transparent)
	t.moveCursor(mathi.Vec2{Y: 1})
}

func cutText(text string, width int) string {
	runes := []rune(text)

	if len(runes) <= width {
		return text
	}

	if width <= 3 {
		return string(runes[:width])
	}

	return string(runes[:width-3]) + "..."
}

func cutTextAndPad(text string, width int) string {
	runes := []rune(text)

	if len(runes) == width {
		return text
	}

	if len(runes) <= width {
		pad := width - len(runes)
		return text + strings.Repeat(" ", pad)
	}

	if width <= 3 {
		return string(runes[:width])
	}

	return string(runes[:width-3]) + "..."
}

package timui

type Theme struct {
	BG     RGBColor
	Text   RGBColor
	Widget WidgetTheme
	Border BorderStyle
}

type WidgetTheme struct {
	BG         RGBColor
	Text       RGBColor
	Line       RGBColor
	HoverBG    RGBColor
	InteractBG RGBColor
	FocusLine  RGBColor
}

type BorderStyle struct {
	Rect       [6]rune
	Horizontal [3]rune
	Vertical   [3]rune
}

var DefaultTheme = Theme{
	BG:   MustRGBS("#000"),
	Text: MustRGBS("#faa"),
	Widget: WidgetTheme{
		BG:         MustRGBS("#004"),
		Text:       MustRGBS("#bbb"),
		Line:       MustRGBS("#22f"),
		HoverBG:    MustRGBS("#22a"),
		InteractBG: MustRGBS("#008"),
		FocusLine:  MustRGBS("#ffa"),
	},
	Border: BorderDouble,
}

var BorderDouble = BorderStyle{
	Rect:       [6]rune{'═', '║', '╔', '╗', '╚', '╝'},
	Vertical:   [3]rune{'╦', '║', '╩'},
	Horizontal: [3]rune{'╠', '═', '╣'},
}

var BorderSingle = BorderStyle{
	Rect:       [6]rune{'─', '│', '┌', '┐', '└', '┘'},
	Vertical:   [3]rune{'┬', '│', '┴'},
	Horizontal: [3]rune{'├', '─', '┤'},
}

var BorderRoundSingle = BorderStyle{
	Rect:       [6]rune{'─', '│', '╭', '╮', '╰', '╯'},
	Vertical:   [3]rune{'┬', '│', '┴'},
	Horizontal: [3]rune{'├', '─', '┤'},
}

func (t *Theme) UseBorder(b BorderStyle) func() {
	before := t.Border
	t.Border = b

	return func() {
		t.Border = before
	}
}

package timui

type Theme struct {
	BG          RGBColor
	Text        RGBColor
	BorderLine  RGBColor
	BorderBG    RGBColor
	Widget      WidgetTheme
	BorderStyle BorderStyle
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
	BG:         MustRGBS("#000"),
	Text:       MustRGBS("#f33"),
	BorderLine: MustRGBS("#3ff"),
	BorderBG:   MustRGBS("#000"),
	Widget: WidgetTheme{
		BG:         MustRGBS("#004"),
		Text:       MustRGBS("#bbb"),
		Line:       MustRGBS("#a0a"),
		HoverBG:    MustRGBS("#22a"),
		InteractBG: MustRGBS("#008"),
		FocusLine:  MustRGBS("#ffa"),
	},
	BorderStyle: BorderDouble,
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

var BorderNone = BorderStyle{
	Rect:       [6]rune{' ', ' ', ' ', ' ', ' ', ' '},
	Vertical:   [3]rune{' ', ' ', ' '},
	Horizontal: [3]rune{' ', ' ', ' '},
}

var BorderBasic = BorderStyle{
	Rect:       [6]rune{'-', '|', '/', '\\', '\\', '/'},
	Vertical:   [3]rune{'+', '|', '+'},
	Horizontal: [3]rune{'+', '-', '+'},
}

func (t *Theme) WithBorder(b BorderStyle, content func()) {
	before := t.BorderStyle
	t.BorderStyle = b

	content()

	t.BorderStyle = before
}

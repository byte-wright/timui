package timui

type Theme struct {
	BG     RGBColor
	Text   RGBColor
	Widget WidgetTheme
}

type WidgetTheme struct {
	BG         RGBColor
	Text       RGBColor
	Line       RGBColor
	HoverBG    RGBColor
	InteractBG RGBColor
	FocusLine  RGBColor
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
}

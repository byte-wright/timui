package timui

import "fmt"

type (
	RGBColor  uint32
	RGBAColor uint32
)

func RGBA(r, g, b, a int) RGBAColor {
	return RGBAColor((a & 0xff << 24) | (r & 0xff << 16) | (g & 0xff << 8) | (b & 0xff))
}

func (rgba RGBAColor) RGBA() (int, int, int, int) {
	i := int(rgba)
	return (i >> 16) & 0xff, (i >> 8) & 0xff, (i) & 0xff, (i >> 24) & 0xff
}

func (rgb RGBAColor) String() string {
	hs := fmt.Sprintf("%08X", uint32(rgb))
	return "#" + hs
}

func RGB(r, g, b int) RGBColor {
	return RGBColor((r & 0xff << 16) | (g & 0xff << 8) | (b & 0xff))
}

func (rgb RGBColor) RGB() (int, int, int) {
	i := int(rgb)
	return (i >> 16) & 0xff, (i >> 8) & 0xff, (i) & 0xff
}

func (rgb RGBColor) MulDiv(factor, div int) RGBColor {
	r, g, b := rgb.RGB()

	return RGB(r*factor/div, g*factor/div, b*factor/div)
}

func (rgb RGBColor) Add(o RGBColor) RGBColor {
	r, g, b := rgb.RGB()
	or, og, ob := o.RGB()

	return RGB(r+or, g+og, b+ob)
}

func (rgb RGBColor) Blend(rgba RGBAColor) RGBColor {
	r, g, b, a := rgba.RGBA()

	c := RGB(r, g, b).MulDiv(a, 0xff)
	o := rgb.MulDiv(0xff-a, 0xff)

	return c.Add(o)
}

func (rgb RGBColor) String() string {
	hs := fmt.Sprintf("%08X", uint32(rgb))
	return "#" + hs[2:]
}

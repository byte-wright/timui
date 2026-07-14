package timui

import "fmt"

type (
	RGBColor  uint32
	RGBAColor uint32
)

var (
	Transparent = RGBA(0, 0, 0, 0)
	White       = RGBA(0xff, 0xff, 0xff, 0xff)
	Black       = RGBA(0x0, 0x0, 0x0, 0xff)
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

func MustRGBAS(s string) RGBAColor {
	col, err := RGBAS(s)
	if err != nil {
		panic(err)
	}

	return col
}

// RGBAS parses "#rgba" or "#rrggbbaa" into an RGBAColor.
func RGBAS(s string) (RGBAColor, error) {
	c, err := parseHexChannels(s, 4)
	if err != nil {
		return RGBA(0, 0, 0, 0x88), err
	}

	return RGBA(c[0], c[1], c[2], c[3]), nil
}

func RGB(r, g, b int) RGBColor {
	return RGBColor((r & 0xff << 16) | (g & 0xff << 8) | (b & 0xff))
}

func MustRGBS(s string) RGBColor {
	col, err := RGBS(s)
	if err != nil {
		panic(err)
	}

	return col
}

// RGBS parses "#rgb" or "#rrggbb" into an RGBColor.
func RGBS(s string) (RGBColor, error) {
	c, err := parseHexChannels(s, 3)
	if err != nil {
		return RGB(0, 0, 0), err
	}

	return RGB(c[0], c[1], c[2]), nil
}

// parseHexChannels reads "#" followed by either count hex digits (short form,
// each digit doubled) or count*2 digits (one byte per channel).
func parseHexChannels(s string, count int) ([4]int, error) {
	var c [4]int

	if len(s) == 0 || s[0] != '#' {
		return c, fmt.Errorf("color must begin with '#'")
	}

	digits := s[1:]

	nibble := func(b byte) (int, bool) {
		switch {
		case b >= '0' && b <= '9':
			return int(b - '0'), true
		case b >= 'a' && b <= 'f':
			return int(b-'a') + 10, true
		case b >= 'A' && b <= 'F':
			return int(b-'A') + 10, true
		}

		return 0, false
	}

	switch len(digits) {
	case count:
		for i := 0; i < count; i++ {
			v, ok := nibble(digits[i])
			if !ok {
				return c, fmt.Errorf("invalid color format '%v'", s)
			}

			c[i] = v * 17
		}
	case count * 2:
		for i := 0; i < count; i++ {
			hi, hiOK := nibble(digits[i*2])
			lo, loOK := nibble(digits[i*2+1])
			if !hiOK || !loOK {
				return c, fmt.Errorf("invalid color format '%v'", s)
			}

			c[i] = hi<<4 + lo
		}
	default:
		return c, fmt.Errorf("invalid color format '%v'", s)
	}

	return c, nil
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

func (rgb RGBColor) RGBA(a int) RGBAColor {
	return RGBAColor(int32(a&0xff<<24) | int32(rgb))
}

func (rgb RGBColor) String() string {
	hs := fmt.Sprintf("%08X", uint32(rgb))
	return "#" + hs[2:]
}

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

func RGBAS(s string) (RGBAColor, error) {
	if len(s) == 0 || s[0] != '#' {
		return RGBA(0, 0, 0, 0x88), fmt.Errorf("Color must begin with '#'")
	}

	var err error

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = fmt.Errorf("invalid color format '%v'", s)

		return 0
	}

	r, g, b, a := 0, 0, 0, 0

	switch len(s) {
	case 9:
		r = int(hexToByte(s[1])<<4 + hexToByte(s[2]))
		g = int(hexToByte(s[3])<<4 + hexToByte(s[4]))
		b = int(hexToByte(s[5])<<4 + hexToByte(s[6]))
		a = int(hexToByte(s[7])<<4 + hexToByte(s[8]))
	case 5:
		r = int(hexToByte(s[1]) * 17)
		g = int(hexToByte(s[2]) * 17)
		b = int(hexToByte(s[3]) * 17)
		a = int(hexToByte(s[4]) * 17)
	default:
		return RGBA(0, 0, 0, 0x88), fmt.Errorf("invalid color format '%v'", s)
	}

	if err != nil {
		return RGBA(0, 0, 0, 0x88), err
	}

	return RGBA(r, g, b, a), err
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

func RGBS(s string) (RGBColor, error) {
	if len(s) == 0 || s[0] != '#' {
		return RGB(0, 0, 0), fmt.Errorf("Color must begin with '#'")
	}

	var err error

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = fmt.Errorf("invalid color format '%v'", s)

		return 0
	}

	r, g, b := 0, 0, 0

	switch len(s) {
	case 7:
		r = int(hexToByte(s[1])<<4 + hexToByte(s[2]))
		g = int(hexToByte(s[3])<<4 + hexToByte(s[4]))
		b = int(hexToByte(s[5])<<4 + hexToByte(s[6]))
	case 4:
		r = int(hexToByte(s[1]) * 17)
		g = int(hexToByte(s[2]) * 17)
		b = int(hexToByte(s[3]) * 17)
	default:
		return RGB(0, 0, 0), fmt.Errorf("invalid color format '%v'", s)
	}

	if err != nil {
		return RGB(0, 0, 0), err
	}

	return RGB(r, g, b), err
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

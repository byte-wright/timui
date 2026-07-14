package timui

import (
	"testing"

	"github.com/byte-wright/expect"
)

func TestRGBColors(t *testing.T) {
	red := RGB(0xff, 0x00, 0x00)
	white50 := RGBA(0xff, 0xff, 0xff, 0xff/2)

	expect.Value(t, "red", red.String()).ToBe("#FF0000")

	expect.Value(t, "red blended", red.Blend(white50).String()).ToBe("#FF7F7F")
}

func TestParseHexColors(t *testing.T) {
	expect.Value(t, "short rgb", MustRGBS("#f80")).ToBe(RGB(0xff, 0x88, 0x00))
	expect.Value(t, "long rgb", MustRGBS("#12aBcD")).ToBe(RGB(0x12, 0xab, 0xcd))

	expect.Value(t, "short rgba", MustRGBAS("#f804")).ToBe(RGBA(0xff, 0x88, 0x00, 0x44))
	expect.Value(t, "long rgba", MustRGBAS("#12aBcDeF")).ToBe(RGBA(0x12, 0xab, 0xcd, 0xef))

	for _, s := range []string{"", "f80", "#f8", "#ff80", "#fg0", "#ff88zz"} {
		_, err := RGBS(s)
		expect.Value(t, "rgb error for "+s, err != nil).ToBe(true)
	}

	_, err := RGBAS("#ff8800")
	expect.Value(t, "rgba error for rgb-length input", err != nil).ToBe(true)
}

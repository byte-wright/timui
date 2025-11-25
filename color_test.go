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

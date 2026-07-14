package timui

import (
	"testing"

	"github.com/byte-wright/expect"
	"gitlab.com/bytewright/gmath/mathi"
)

func TestFinishRunsNestedDeferred(t *testing.T) {
	tui := New(&testBackend{size: mathi.Vec2{X: 20, Y: 10}})

	order := []string{}
	tui.runAfter(func() {
		order = append(order, "outer")
		tui.runAfter(func() {
			order = append(order, "nested")
		})
	})

	tui.Finish()

	expect.Value(t, "deferred run order", order).ToBe([]string{"outer", "nested"})
}

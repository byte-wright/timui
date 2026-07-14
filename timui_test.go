package timui_test

import (
	"testing"

	"github.com/byte-wright/expect"
	"github.com/byte-wright/timui"
	"github.com/byte-wright/timui/internal/test"
)

func TestFinishRunsNestedDeferred(t *testing.T) {
	tui := timui.New(test.NewBackend(20, 10))

	order := []string{}
	tui.RunAfterForTest(func() {
		order = append(order, "outer")
		tui.RunAfterForTest(func() {
			order = append(order, "nested")
		})
	})

	tui.Finish()

	expect.Value(t, "deferred run order", order).ToBe([]string{"outer", "nested"})
}

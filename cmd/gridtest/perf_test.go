package main

import (
	"testing"

	"github.com/byte-wright/timui"
	"gitlab.com/bytewright/gmath/mathi"
)

type nopBackend struct{}

func (n *nopBackend) MousePosition() mathi.Vec2 {
	return mathi.Vec2{X: 15, Y: 3}
}

func (n *nopBackend) MousePressed(key timui.Key) bool {
	return false
}

func (n *nopBackend) Render() {
}

func (n *nopBackend) Set(pos mathi.Vec2, char rune, fg uint32, bg uint32) {
}

func (n *nopBackend) Size() mathi.Vec2 {
	return mathi.Vec2{X: 120, Y: 60}
}

func BenchmarkRenderLoop(b *testing.B) {
	tui := timui.New(&nopBackend{})

	for i := 0; i < b.N; i++ {
		render(tui)
	}
}

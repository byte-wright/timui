package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/byte-wright/timui"
	"gitlab.com/bytewright/gmath/mathi"
)

type recordBackend struct {
	chars  [][]rune
	colors uint64
}

func newRecordBackend() *recordBackend {
	chars := make([][]rune, 60)
	for i := range chars {
		chars[i] = make([]rune, 180)
		for j := range chars[i] {
			chars[i][j] = ' '
		}
	}

	return &recordBackend{chars: chars}
}

func (n *recordBackend) MousePosition() mathi.Vec2 { return mathi.Vec2{X: 15, Y: 3} }

func (n *recordBackend) MousePressed(key timui.Key) bool { return false }

func (n *recordBackend) Render() {}

func (n *recordBackend) Set(pos mathi.Vec2, char rune, fg uint32, bg uint32) {
	if char != 0 {
		n.chars[pos.Y][pos.X] = char
	}

	n.colors = n.colors*31 + uint64(pos.X)<<32 + uint64(pos.Y)<<16 + uint64(fg) + uint64(bg)
}

func (n *recordBackend) Size() mathi.Vec2 { return mathi.Vec2{X: 180, Y: 60} }

func TestDumpFrame(t *testing.T) {
	out := os.Getenv("FRAMEDUMP_OUT")
	if out == "" {
		t.Skip("FRAMEDUMP_OUT not set")
	}

	be := newRecordBackend()
	tui := timui.New(be)

	render(tui)

	dump := ""
	for _, row := range be.chars {
		dump += string(row) + "\n"
	}
	dump += fmt.Sprintf("colors: %x\n", be.colors)

	err := os.WriteFile(out, []byte(dump), 0o644)
	if err != nil {
		t.Fatal(err)
	}
}

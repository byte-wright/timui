package internal

import (
	"testing"

	"github.com/byte-wright/expect"
)

func TestIDStableForSamePath(t *testing.T) {
	a := NewIDManager()
	b := NewIDManager()

	a.Push("panel")
	b.Push("panel")
	expect.Value(t, "same path gives same id", a.Push("button") == b.Push("button")).ToBe(true)
}

func TestIDDiffersByLabelAndScope(t *testing.T) {
	m := NewIDManager()

	root := m.Peek()
	first := m.Push("a")
	m.Pop()
	second := m.Push("b")
	m.Pop()

	expect.Value(t, "different labels give different ids", first != second).ToBe(true)

	m.Push("scope")
	scoped := m.Push("a")
	expect.Value(t, "same label in different scope gives different id", scoped != first).ToBe(true)

	m.Pop()
	m.Pop()
	expect.Value(t, "pop restores parent scope", m.Peek()).ToBe(root)
}

func TestPushIDMixesIntoScope(t *testing.T) {
	m := NewIDManager()

	cid := m.Push("dropdown")
	m.Pop()

	restored := m.PushID(cid)
	expect.Value(t, "pushed id is mixed, not verbatim", restored != cid).ToBe(true)

	child := m.Push("selection")
	m.Pop()
	m.Pop()

	m.PushID(cid)
	expect.Value(t, "children under re-pushed id are stable", m.Push("selection")).ToBe(child)
}

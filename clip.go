package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

type clipManager struct {
	clips []mathi.Box2
}

func newClipManager() *clipManager {
	return &clipManager{}
}

func (g *Timui) PushClip(area mathi.Box2) {
	g.clipManager.clips = append(g.clipManager.clips, area)
}

func (g *Timui) PeekClip() mathi.Box2 {
	if len(g.clipManager.clips) > 0 {
		return g.clipManager.clips[len(g.clipManager.clips)-1]
	}

	return mathi.Box2{To: g.backend.Size()}
}

func (t *Timui) ClipContains(pos mathi.Vec2) bool {
	if len(t.clipManager.clips) == 0 {
		size := t.backend.Size()
		return pos.X >= 0 && pos.X < size.X && pos.Y >= 0 && pos.Y < size.Y
	}

	return t.clipManager.clips[len(t.clipManager.clips)-1].Contains(pos)
}

func (g *Timui) PopClip() {
	g.clipManager.clips = g.clipManager.clips[:len(g.clipManager.clips)-1]
}

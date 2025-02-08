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

func (g *Timui) PopClip() {
	g.clipManager.clips = g.clipManager.clips[:len(g.clipManager.clips)-1]
}

package timui

import "strings"

type idManager struct {
	ids []string
}

type ID struct {
	value string
}

func newIDManager() *idManager {
	return &idManager{
		ids: []string{},
	}
}

func (g *Timui[B]) PushID(id string) ID {
	g.idManager.ids = append(g.idManager.ids, id)
	return g.PeekID()
}

func (g *Timui[B]) PeekID() ID {
	return ID{value: strings.Join(g.idManager.ids, "-")}
}

func (g *Timui[B]) PopID() {
	g.idManager.ids = g.idManager.ids[:len(g.idManager.ids)-1]
}

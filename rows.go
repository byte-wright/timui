package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

type Rows struct {
	g         *Timui
	positions []splitRange
	area      mathi.Box2
	row       int
}

func (g *Timui) Rows(pos *SplitOptions) *Rows {
	positions := pos.calculatePositions(g.CurrentArea().Size().Y)

	area := g.CurrentArea()
	firstArea := *area
	yStart := firstArea.From.Y
	firstArea.From.Y = yStart + positions[0].from
	firstArea.To.Y = yStart + positions[0].to

	g.PushArea(firstArea)

	return &Rows{
		g:         g,
		positions: positions,
		area:      *area,
	}
}

func (s *Rows) Next() {
	s.g.PopArea()
	s.row += 1

	area := s.area
	yStart := area.From.Y
	area.From.Y = yStart + s.positions[s.row].from
	area.To.Y = yStart + s.positions[s.row].to

	s.g.PushArea(area)
}

func (s *Rows) Finish() {
	s.g.PopArea()
}

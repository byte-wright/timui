package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

type Columns[B Backend] struct {
	g         *Timui[B]
	positions []splitRange
	area      mathi.Box2
	column    int
	maxCursor int
}

func (g *Timui[B]) Columns(opts *SplitOptions) *Columns[B] {
	positions := opts.calculatePositions(g.CurrentArea().Size().X)

	area := g.CurrentArea()
	firstArea := *area
	firstArea.From.X = area.From.X + positions[0].from
	firstArea.To.X = area.From.X + positions[0].to

	g.PushArea(firstArea)

	return &Columns[B]{
		g:         g,
		positions: positions,
		area:      *area,
	}
}

func (s *Columns[B]) Next() {
	if s.g.CurrentArea().From.Y > s.maxCursor {
		s.maxCursor = s.g.CurrentArea().From.Y
	}
	s.g.PopArea()
	s.column += 1

	area := s.area
	area.From.X = s.area.From.X + s.positions[s.column].from
	area.To.X = s.area.From.X + s.positions[s.column].to

	s.g.PushArea(area)
}

func (s *Columns[B]) Finish() {
	if s.g.CurrentArea().From.Y > s.maxCursor {
		s.maxCursor = s.g.CurrentArea().From.Y
	}
	s.g.PopArea()

	s.g.CurrentArea().From.Y = s.maxCursor
}

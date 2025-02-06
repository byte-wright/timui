package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

type Grid[B Backend] struct {
	t    *Timui[B]
	area mathi.Box2
}

func (t *Timui[B]) Grid() *Grid[B] {
	t.Border(t.Theme.BorderStyle.Rect, t.Theme.BorderLine, t.Theme.BorderBG)

	area := *t.CurrentArea()
	originalArea := area

	pad := mathi.Vec2{X: 2, Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	t.PushArea(area)

	return &Grid[B]{
		t:    t,
		area: originalArea,
	}
}

func (g *Grid[B]) Finish() {
	g.t.PopArea()
}

type GridRows[B Backend] struct {
	t         *Timui[B]
	positions []splitRange
	area      mathi.Box2
	row       int
}

func (g *Grid[B]) Rows(pos *SplitOptions) *GridRows[B] {
	if pos.padding != 0 {
		panic("grid rows padding must be zero")
	}

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(g.area.Size().Y)

	gridRows := &GridRows[B]{
		t:         g.t,
		positions: positions,
		area:      g.area,
		row:       1,
	}

	area := gridRows.currentArea()

	pad := mathi.Vec2{X: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)

	return gridRows
}

func (g *Grid[B]) Columns(pos *SplitOptions) *GridColumns[B] {
	if pos.padding != 0 {
		panic("grid rows padding must be zero")
	}

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(g.area.Size().X)

	gridRows := &GridColumns[B]{
		t:         g.t,
		positions: positions,
		area:      g.area,
		column:    1,
	}

	area := gridRows.currentArea()

	pad := mathi.Vec2{Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)

	return gridRows
}

func (g *GridRows[B]) Next() {
	g.t.PopArea()
	g.row += 1

	g.t.PushArea(g.currentArea())

	g.t.HLine(g.t.Theme.BorderStyle.Horizontal, g.t.Theme.BorderLine, g.t.Theme.BorderBG)
	g.t.PopArea()

	g.row += 1

	area := g.currentArea()

	pad := mathi.Vec2{X: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)
}

func (g *GridRows[B]) currentArea() mathi.Box2 {
	area := g.area
	yStart := area.From.Y
	area.From.Y = yStart + g.positions[g.row].from
	area.To.Y = yStart + g.positions[g.row].to

	return area
}

func (g *GridRows[B]) currentCompleteArea() mathi.Box2 {
	area := g.currentArea()
	area.From.Y -= 1
	area.To.Y += 1

	return area
}

func (g *GridRows[B]) Finish() {
	g.t.PopArea()
}

type GridColumns[B Backend] struct {
	t         *Timui[B]
	positions []splitRange
	area      mathi.Box2
	column    int
}

func (g *GridRows[B]) Columns(pos *SplitOptions) *GridColumns[B] {
	if pos.padding != 0 {
		panic("grid columns padding must be zero")
	}

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(g.area.Size().X)

	gridColumns := &GridColumns[B]{
		t:         g.t,
		positions: positions,
		area:      g.currentCompleteArea(),
		column:    1,
	}

	area := gridColumns.currentArea()

	pad := mathi.Vec2{Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)

	return gridColumns
}

func (g *GridRows[B]) Rows(pos *SplitOptions) *GridRows[B] {
	if pos.padding != 0 {
		panic("grid columns padding must be zero")
	}

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(g.currentCompleteArea().Size().Y)

	gridRows := &GridRows[B]{
		t:         g.t,
		positions: positions,
		area:      g.currentCompleteArea(),
		row:       1,
	}

	area := gridRows.currentArea()

	pad := mathi.Vec2{X: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)

	return gridRows
}

func (g *GridColumns[B]) Next() {
	g.t.PopArea()
	g.column += 1

	area := g.currentArea()

	g.t.PushArea(area)

	g.t.VLine(g.t.Theme.BorderStyle.Vertical, g.t.Theme.BorderLine, g.t.Theme.BorderBG)
	g.t.PopArea()

	g.column += 1

	area = g.currentArea()
	pad := mathi.Vec2{Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)
}

func (g *GridColumns[B]) Rows(pos *SplitOptions) *GridRows[B] {
	if pos.padding != 0 {
		panic("grid rows padding must be zero")
	}

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(g.area.Size().Y)

	gridRows := &GridRows[B]{
		t:         g.t,
		positions: positions,
		area:      g.currentCompleteArea(),
		row:       1,
	}

	area := gridRows.currentArea()

	pad := mathi.Vec2{X: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)

	return gridRows
}

func (g *GridColumns[B]) Columns(pos *SplitOptions) *GridColumns[B] {
	if pos.padding != 0 {
		panic("grid rows padding must be zero")
	}

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(g.area.Size().X)

	gridColumns := &GridColumns[B]{
		t:         g.t,
		positions: positions,
		area:      g.currentCompleteArea(),
		column:    1,
	}

	area := gridColumns.currentArea()

	pad := mathi.Vec2{Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)

	return gridColumns
}

func (g *GridColumns[B]) currentArea() mathi.Box2 {
	area := g.area
	area.From.X = g.area.From.X + g.positions[g.column].from
	area.To.X = g.area.From.X + g.positions[g.column].to

	return area
}

func (g *GridColumns[B]) currentCompleteArea() mathi.Box2 {
	area := g.currentArea()
	area.From.X -= 1
	area.To.X += 1

	return area
}

func (g *GridColumns[B]) Finish() {
	g.t.PopArea()
}

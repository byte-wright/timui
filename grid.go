package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

type Grid struct {
	t    *Timui
	area mathi.Box2
}

func (t *Timui) Grid(body func(grid *Grid)) {
	t.Border(t.Theme.BorderStyle.Rect, t.Theme.BorderLine, t.Theme.BorderBG)

	area := *t.CurrentArea()
	originalArea := area

	pad := mathi.Vec2{X: 1, Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	t.PushArea(area)

	grid := &Grid{
		t:    t,
		area: originalArea,
	}

	body(grid)

	grid.finish()
}

func (g *Grid) finish() {
	g.t.PopArea()
}

type GridRows struct {
	t         *Timui
	positions []splitRange
	area      mathi.Box2
	row       int
}

func (g *Grid) Rows(pos *SplitOptions, body func(rows *GridRows)) {
	if pos.padding != 0 {
		panic("grid rows padding must be zero")
	}

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(g.area.Size().Y)

	gridRows := &GridRows{
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

	body(gridRows)

	gridRows.finish()
}

func (g *Grid) Columns(pos *SplitOptions, body func(rows *GridColumns)) {
	if pos.padding != 0 {
		panic("grid rows padding must be zero")
	}

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(g.area.Size().X)

	gridRows := &GridColumns{
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

	body(gridRows)

	gridRows.finish()
}

func (g *GridRows) Next() {
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

func (g *GridRows) currentArea() mathi.Box2 {
	area := g.area
	yStart := area.From.Y
	area.From.Y = yStart + g.positions[g.row].from
	area.To.Y = yStart + g.positions[g.row].to

	return area
}

func (g *GridRows) currentCompleteArea() mathi.Box2 {
	area := g.currentArea()
	area.From.Y -= 1
	area.To.Y += 1

	return area
}

func (g *GridRows) finish() {
	g.t.PopArea()
}

type GridColumns struct {
	t         *Timui
	positions []splitRange
	area      mathi.Box2
	column    int
}

func (g *GridRows) Columns(pos *SplitOptions, body func(columns *GridColumns)) {
	if pos.padding != 0 {
		panic("grid columns padding must be zero")
	}

	completeArea := g.currentCompleteArea()

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(completeArea.Size().X)

	gridColumns := &GridColumns{
		t:         g.t,
		positions: positions,
		area:      completeArea,
		column:    1,
	}

	area := gridColumns.currentArea()

	pad := mathi.Vec2{Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)

	body(gridColumns)

	gridColumns.finish()
}

func (g *GridRows) Rows(pos *SplitOptions, body func(rows *GridRows)) {
	if pos.padding != 0 {
		panic("grid columns padding must be zero")
	}

	completeArea := g.currentCompleteArea()

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(completeArea.Size().Y)

	gridRows := &GridRows{
		t:         g.t,
		positions: positions,
		area:      completeArea,
		row:       1,
	}

	area := gridRows.currentArea()

	pad := mathi.Vec2{X: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)

	body(gridRows)

	gridRows.finish()
}

func (g *GridColumns) Next() {
	g.t.PopArea()
	g.column += 1

	g.t.PushArea(g.currentArea())

	g.t.VLine(g.t.Theme.BorderStyle.Vertical, g.t.Theme.BorderLine, g.t.Theme.BorderBG)
	g.t.PopArea()

	g.column += 1

	area := g.currentArea()
	pad := mathi.Vec2{Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)
}

func (g *GridColumns) Rows(pos *SplitOptions, body func(rows *GridRows)) {
	if pos.padding != 0 {
		panic("grid rows padding must be zero")
	}

	completeArea := g.currentCompleteArea()

	pos.insertFixedBetween(1)

	positions := pos.calculatePositions(completeArea.Size().Y)

	gridRows := &GridRows{
		t:         g.t,
		positions: positions,
		area:      completeArea,
		row:       1,
	}

	area := gridRows.currentArea()

	pad := mathi.Vec2{X: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)

	body(gridRows)

	gridRows.finish()
}

func (g *GridColumns) Columns(pos *SplitOptions, body func(columns *GridColumns)) {
	if pos.padding != 0 {
		panic("grid rows padding must be zero")
	}

	pos.insertFixedBetween(1)
	completeArea := g.currentCompleteArea()
	positions := pos.calculatePositions(completeArea.Size().X)

	gridColumns := &GridColumns{
		t:         g.t,
		positions: positions,
		area:      completeArea,
		column:    1,
	}

	area := gridColumns.currentArea()

	pad := mathi.Vec2{Y: 1}
	area.From = area.From.Add(pad)
	area.To = area.To.Sub(pad)

	g.t.PushArea(area)

	body(gridColumns)

	gridColumns.finish()
}

func (g *GridColumns) currentArea() mathi.Box2 {
	area := g.area
	area.From.X = g.area.From.X + g.positions[g.column].from
	area.To.X = g.area.From.X + g.positions[g.column].to

	return area
}

func (g *GridColumns) currentCompleteArea() mathi.Box2 {
	area := g.currentArea()
	area.From.X -= 1
	area.To.X += 1

	return area
}

func (g *GridColumns) finish() {
	g.t.PopArea()
}

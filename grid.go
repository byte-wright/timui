package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

type Grid struct {
	t    *Timui
	area mathi.Box2
}

// Grid draws a bordered container and runs body inside its padded area.
// Rows and Columns on the grid subdivide it with divider lines that join
// the surrounding border.
func (t *Timui) Grid(body func(grid *Grid)) {
	t.Border(t.Theme.BorderStyle.Rect, t.Theme.BorderLine, t.Theme.BorderBG)

	area := *t.CurrentArea()
	grid := &Grid{t: t, area: area}

	pad := mathi.Vec2{X: 1, Y: 1}
	padded := area
	padded.From = padded.From.Add(pad)
	padded.To = padded.To.Sub(pad)

	t.WithArea(padded, func() {
		body(grid)
	})
}

func (g *Grid) Rows(opts *SplitOptions, cells ...func(cell *GridCell)) {
	g.t.gridSplit(axisRows, g.area, opts, cells)
}

func (g *Grid) Columns(opts *SplitOptions, cells ...func(cell *GridCell)) {
	g.t.gridSplit(axisColumns, g.area, opts, cells)
}

// GridCell is the context of one grid cell. Its area is the complete cell
// including the border lines around it, so nested splits overlap the parent's
// lines and their divider glyphs merge into junctions.
type GridCell struct {
	t    *Timui
	area mathi.Box2
}

func (c *GridCell) Rows(opts *SplitOptions, cells ...func(cell *GridCell)) {
	c.t.gridSplit(axisRows, c.area, opts, cells)
}

func (c *GridCell) Columns(opts *SplitOptions, cells ...func(cell *GridCell)) {
	c.t.gridSplit(axisColumns, c.area, opts, cells)
}

type gridAxis int

const (
	axisRows gridAxis = iota
	axisColumns
)

// gridSplit runs one cell func per split entry inside area. The split gets a
// one-cell gutter interleaved around every entry; the outermost gutters land
// on area's border lines, the inner ones get divider lines drawn into them.
func (t *Timui) gridSplit(a gridAxis, area mathi.Box2, opts *SplitOptions, cells []func(cell *GridCell)) {
	if opts.padding != 0 {
		panic("grid " + a.name() + " padding must be zero")
	}

	positions := opts.withFixedBetween(1).calculatePositions(a.extent(area.Size()))

	if len(cells)*2+1 != len(positions) {
		panic("grid " + a.name() + " cell count must match split count")
	}

	for i, cell := range cells {
		if i > 0 {
			gutter := a.cellArea(area, positions[i*2])
			t.WithArea(gutter, func() {
				a.divider(t)
			})
		}

		cellArea := a.cellArea(area, positions[i*2+1])
		c := &GridCell{t: t, area: a.expand(cellArea, 1)}

		t.WithArea(a.crossPad(cellArea, 1), func() {
			cell(c)
		})
	}
}

func (a gridAxis) name() string {
	if a == axisRows {
		return "rows"
	}

	return "columns"
}

func (a gridAxis) extent(size mathi.Vec2) int {
	if a == axisRows {
		return size.Y
	}

	return size.X
}

func (a gridAxis) cellArea(area mathi.Box2, r splitRange) mathi.Box2 {
	if a == axisRows {
		start := area.From.Y
		area.From.Y = start + r.from
		area.To.Y = start + r.to
	} else {
		start := area.From.X
		area.From.X = start + r.from
		area.To.X = start + r.to
	}

	return area
}

func (a gridAxis) expand(area mathi.Box2, n int) mathi.Box2 {
	if a == axisRows {
		area.From.Y -= n
		area.To.Y += n
	} else {
		area.From.X -= n
		area.To.X += n
	}

	return area
}

func (a gridAxis) crossPad(area mathi.Box2, n int) mathi.Box2 {
	if a == axisRows {
		area.From.X += n
		area.To.X -= n
	} else {
		area.From.Y += n
		area.To.Y -= n
	}

	return area
}

func (a gridAxis) divider(t *Timui) {
	if a == axisRows {
		t.HLine(t.Theme.BorderStyle.Horizontal, t.Theme.BorderLine, t.Theme.BorderBG)
	} else {
		t.VLine(t.Theme.BorderStyle.Vertical, t.Theme.BorderLine, t.Theme.BorderBG)
	}
}

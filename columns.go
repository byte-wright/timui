package timui

// Columns splits the current area horizontally and runs one cell func per
// split entry, each inside its column area. Afterwards the parent cursor is
// advanced to the bottom of the tallest column. Panics if the cell count
// does not match the split count.
func (g *Timui) Columns(opts *SplitOptions, cells ...func()) {
	positions := opts.calculatePositions(g.CurrentArea().Size().X)
	if len(cells) != len(positions) {
		panic("columns cell count must match split count")
	}

	area := *g.CurrentArea()
	maxCursor := area.From.Y

	for i, cell := range cells {
		cellArea := area
		cellArea.From.X = area.From.X + positions[i].from
		cellArea.To.X = area.From.X + positions[i].to

		g.WithArea(cellArea, func() {
			cell()
			if g.CurrentArea().From.Y > maxCursor {
				maxCursor = g.CurrentArea().From.Y
			}
		})
	}

	g.CurrentArea().From.Y = maxCursor
}

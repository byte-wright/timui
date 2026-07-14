package timui

// Rows splits the current area vertically and runs one cell func per split
// entry, each inside its row area. Panics if the cell count does not match
// the split count.
func (g *Timui) Rows(opts *SplitOptions, cells ...func()) {
	positions := opts.calculatePositions(g.CurrentArea().Size().Y)
	if len(cells) != len(positions) {
		panic("rows cell count must match split count")
	}

	area := *g.CurrentArea()
	yStart := area.From.Y

	for i, cell := range cells {
		cellArea := area
		cellArea.From.Y = yStart + positions[i].from
		cellArea.To.Y = yStart + positions[i].to

		g.WithArea(cellArea, cell)
	}
}

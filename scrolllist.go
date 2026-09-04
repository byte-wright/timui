package timui

// ScrollState is the caller owned position of a scroll list.
type ScrollState struct {
	// Offset is the index of the first visible row.
	Offset int

	// Follow pins the view to the last row while rows are appended. It is cleared
	// when the user scrolls away from the bottom and set again when they return to
	// it, so a log tails by default but stays put once it is read.
	Follow bool

	// ViewHeight is how many rows fit, written on every render so the caller can
	// page by a screenful.
	ViewHeight int
}

// ScrollListV renders count rows of height one and calls row only for the visible
// ones. Unlike ScrollAreaV, which draws all of its content and clips it away, the
// cost of a list does not grow with the number of rows.
func (g *Timui) ScrollListV(id string, state *ScrollState, count int, row func(i int)) {
	g.id.Push(id)
	defer g.id.Pop()

	area := *g.CurrentArea()

	listArea := area
	listArea.To.X -= 1

	height := listArea.Size().Y
	state.ViewHeight = height

	max := count - height
	if max < 0 {
		max = 0
	}

	g.WithArea(area, func() {
		wheelInput := g.MouseInput("wheel")
		wheelInput.passive = true

		wheel := wheelInput.Wheel()
		if wheel.Y != 0 {
			state.Offset += wheel.Y * wheelStep
			state.Follow = state.Offset >= max
		}
	})

	factor := float32(0)
	if max > 0 {
		factor = float32(clamp(state.Offset, 0, max)) / float32(max)
	}

	dragged, moved := g.scrollBarV(area, count, factor)
	if moved {
		state.Offset = int(float32(max)*dragged + 0.5)
		state.Follow = state.Offset >= max
	}

	if state.Follow {
		state.Offset = max
	}

	state.Offset = clamp(state.Offset, 0, max)

	g.WithArea(listArea, func() {
		g.WithClip(listArea, func() {
			g.WithCopyRect(listArea, func() {
				for i := 0; i < height && state.Offset+i < count; i++ {
					rowArea := listArea
					rowArea.From.Y = listArea.From.Y + i
					rowArea.To.Y = rowArea.From.Y + 1

					g.WithArea(rowArea, func() {
						row(state.Offset + i)
					})
				}
			})
		})
	})
}

// PageUp moves the view up by a screenful and stops following.
func (s *ScrollState) PageUp() {
	s.Scroll(-s.ViewHeight)
}

// PageDown moves the view down by a screenful.
func (s *ScrollState) PageDown() {
	s.Scroll(s.ViewHeight)
}

// Scroll moves the view by rows. The next render clamps the offset and decides
// whether the view ended up at the bottom again.
func (s *ScrollState) Scroll(rows int) {
	s.Offset += rows
	s.Follow = false

	if s.Offset < 0 {
		s.Offset = 0
	}
}

// Top jumps to the first row.
func (s *ScrollState) Top() {
	s.Offset = 0
	s.Follow = false
}

// Bottom jumps to the last row and follows from there.
func (s *ScrollState) Bottom() {
	s.Follow = true
}

func clamp(v, low, high int) int {
	if v < low {
		return low
	}

	if v > high {
		return high
	}

	return v
}

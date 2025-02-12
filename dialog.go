package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

func (t *Timui) Dialog(title string, visible *bool, content func()) {
	did := t.id.Push(title)

	if *visible {
		t.runAfter(func() {
			t.id.PushID(did)

			modal := t.MouseInput("modal")
			if modal.LeftReleased() {
				*visible = false
			}

			t.SetAreaAlpha(0, RGBA(0, 0, 0, 0xaa), RGBA(0, 0, 0, 0xaa))

			area := *t.CurrentArea()

			area.From = area.To.Mul(1).Div(4)
			area.To = area.To.Mul(3).Div(4)

			t.PushArea(area)

			t.MouseInput("dialog") // to cancel the modal backdrop close click

			t.Grid(func(grid *Grid) {
				t.SetArea(' ', t.Theme.Text, t.Theme.BG)

				grid.Rows(Split().Fixed(1).Factor(1), func(rows *GridRows) {
					rows.Columns(Split().Factor(1).Fixed(3), func(title *GridColumns) {
						t.Pad(0, 1, 0, 1, func() {
							t.Label("title")
						})

						title.Next()

						close := t.MouseInputForSize("close", mathi.Vec2{X: 3, Y: 1})
						if close.LeftReleased() {
							*visible = false
						}
						if close.Hovered() > 0 {
							t.SetAreaAlpha(0, Transparent, RGBA(0x55, 0, 0, 0xff))
						}
						t.Label(" X ")
					})

					rows.Next()
					t.Pad(0, 1, 0, 1, content)
				})
			})

			t.PopArea()

			t.id.Pop()
		})
	}

	t.id.Pop()
}

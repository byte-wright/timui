package timui

import (
	"gitlab.com/bytewright/gmath/mathi"
)

func (t *Timui) Dialog(title string, visible *bool) {
	did := t.id.Push(title)

	if *visible {
		t.runAfter(func() {
			t.id.PushID(did)

			modal := t.MouseInput("modal", *t.CurrentArea())
			if modal.LeftReleased() {
				*visible = false
			}

			t.SetAreaAlpha(0, RGBA(0, 0, 0, 0xaa), RGBA(0, 0, 0, 0xaa))

			area := *t.CurrentArea()

			area.From = area.To.Mul(1).Div(4)
			area.To = area.To.Mul(3).Div(4)

			t.PushArea(area)

			t.MouseInput("dialog", mathi.NewBox2FromSize(mathi.Vec2{}, area.Size()))

			grid := t.Grid()

			t.SetArea(' ', t.Theme.Text, t.Theme.BG)

			rows := grid.Rows(Split().Fixed(1).Factor(1))

			title := rows.Columns(Split().Factor(1).Fixed(3))

			tp := t.Pad(0, 1, 0, 1)
			t.Label("FISDH!")
			tp.Finish()

			title.Next()

			close := t.MouseInput("close", mathi.Box2{From: mathi.Vec2{}, To: mathi.Vec2{X: 3, Y: 1}})
			if close.LeftReleased() {
				*visible = false
			}
			if close.Hovered() > 0 {
				t.SetAreaAlpha(0, Transparent, RGBA(0x55, 0, 0, 0xff))
			}
			t.Label(" X ")
			title.Finish()

			rows.Next()
			cp := t.Pad(0, 1, 0, 1)
			t.Label("cntnt")
			cp.Finish()

			rows.Finish()

			grid.Finish()

			t.PopArea()

			t.id.Pop()
		})
	}

	t.id.Pop()
}

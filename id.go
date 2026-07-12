package timui

// PushID pushes an id scope onto the stack. Widgets created until the matching
// PopID derive their id relative to this scope, so identically-labelled widgets
// (for example two "+" buttons in different rows) stay distinct.
func (t *Timui) PushID(id string) {
	t.id.Push(id)
}

// PopID removes the most recently pushed id scope.
func (t *Timui) PopID() {
	t.id.Pop()
}

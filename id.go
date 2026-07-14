package timui

// WithID runs body inside an id scope. Widgets created by body derive their id
// relative to this scope, so identically-labelled widgets (for example two "+"
// buttons in different rows) stay distinct.
func (t *Timui) WithID(id string, body func()) {
	t.id.Push(id)
	body()
	t.id.Pop()
}

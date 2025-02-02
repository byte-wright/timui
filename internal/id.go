package internal

type IDManager struct {
	ids []ID
}

type ID struct {
	value string
}

func NewIDManager() *IDManager {
	return &IDManager{
		ids: []ID{},
	}
}

func (m *IDManager) Push(id string) ID {
	m.ids = append(m.ids,
		ID{m.Peek().value + "---" + id},
	)

	return m.Peek()
}

func (m *IDManager) PushID(id ID) ID {
	m.ids = append(m.ids, id)
	return m.Peek()
}

func (m *IDManager) Peek() ID {
	if len(m.ids) == 0 {
		return ID{value: "root"}
	}

	return m.ids[len(m.ids)-1]
}

func (m *IDManager) Pop() {
	m.ids = m.ids[:len(m.ids)-1]
}

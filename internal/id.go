package internal

const (
	fnvOffset64 = 1469598103934665603
	fnvPrime64  = 1099511628211
)

type IDManager struct {
	ids []ID
}

type ID uint64

func NewIDManager() *IDManager {
	return &IDManager{
		ids: []ID{0},
	}
}

func hashPushString(prev ID, s string) ID {
	h := prev ^ fnvOffset64
	for i := 0; i < len(s); i++ {
		h ^= ID(s[i])
		h *= fnvPrime64
	}

	return h
}

func hashPushID(prev ID, c ID) ID {
	h := prev ^ fnvOffset64

	// Mix the 8 bytes of c into the hash, just like string bytes.
	v := uint64(c)
	for i := 0; i < 8; i++ {
		h ^= ID(v & 0xFF)
		h *= fnvPrime64
		v >>= 8
	}

	return h
}

func (m *IDManager) Push(id string) ID {
	m.ids = append(m.ids, hashPushString(m.Peek(), id))
	return m.Peek()
}

func (m *IDManager) PushID(id ID) ID {
	m.ids = append(m.ids, hashPushID(m.Peek(), id))
	return m.Peek()
}

func (m *IDManager) Peek() ID {
	return m.ids[len(m.ids)-1]
}

func (m *IDManager) Pop() {
	m.ids = m.ids[:len(m.ids)-1]
}

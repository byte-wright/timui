package timui

type SplitOptions struct {
	splits  []split
	padding int
}

func (s *SplitOptions) Pad(padding int) *SplitOptions {
	s.padding = padding
	return s
}

type split struct {
	factor float32
	fixed  int
}

func ByFactors(factors []float32) *SplitOptions {
	so := &SplitOptions{}
	for _, f := range factors {
		so.splits = append(so.splits, split{factor: f})
	}

	return so
}

type splitRange struct {
	from int
	to   int
}

func (so *SplitOptions) calculatePositions(width int) []splitRange {
	positions := make([]splitRange, len(so.splits))
	remaining := width - (len(so.splits)-1)*so.padding

	totalF := float32(0)

	for _, f := range so.splits {
		totalF += f.factor
	}

	pos := 0
	for i, f := range so.splits {
		positions[i].from = pos

		w := int(float32(remaining) * f.factor / totalF)
		pos += w
		positions[i].to = pos
		pos += so.padding

		totalF -= f.factor
		remaining -= w
	}

	return positions
}

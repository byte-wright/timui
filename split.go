package timui

type SplitOptions struct {
	splits  []split
	padding int
}

type split struct {
	factor float32
	fixed  int
}

type splitRange struct {
	from int
	to   int
}

func Split() *SplitOptions {
	return &SplitOptions{}
}

func (s *SplitOptions) Factors(factors ...float32) *SplitOptions {
	for _, f := range factors {
		s.splits = append(s.splits, split{factor: f})
	}

	return s
}

func (s *SplitOptions) Fixed(fixed ...int) *SplitOptions {
	for _, f := range fixed {
		s.splits = append(s.splits, split{fixed: f})
	}

	return s
}

func (s *SplitOptions) Pad(padding int) *SplitOptions {
	s.padding = padding
	return s
}

func (s *SplitOptions) Add(factor float32, fixed int) *SplitOptions {
	s.splits = append(s.splits, split{factor: factor, fixed: fixed})
	return s
}

func (so *SplitOptions) calculatePositions(width int) []splitRange {
	totalFactor := float32(0)
	totalFixed := 0

	for _, f := range so.splits {
		totalFactor += f.factor
		totalFixed += f.fixed
	}

	positions := make([]splitRange, len(so.splits))
	remaining := width - (len(so.splits)-1)*so.padding - totalFixed

	pos := 0
	for i, f := range so.splits {
		positions[i].from = pos

		w := int(float32(remaining) * f.factor / totalFactor)
		pos += w
		pos += f.fixed
		positions[i].to = pos
		pos += so.padding

		totalFactor -= f.factor
		remaining -= w
	}

	return positions
}

package timui

import (
	"math"
)

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

func (s *SplitOptions) Factor(factors ...float32) *SplitOptions {
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

// withFixedBetween returns a copy with a fixed-size entry interleaved around
// every split; the receiver is left untouched so it can be reused across
// frames.
func (s *SplitOptions) withFixedBetween(v int) *SplitOptions {
	ns := &SplitOptions{padding: s.padding}
	ns.splits = make([]split, 0, len(s.splits)*2+1)

	ns.splits = append(ns.splits, split{fixed: v})
	for _, f := range s.splits {
		ns.splits = append(ns.splits, f, split{fixed: v})
	}

	return ns
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

		w := 0
		if remaining > 0 {
			w = int(math.Round(float64(float32(remaining) * f.factor / totalFactor)))
		}

		pos += w
		pos += f.fixed
		positions[i].to = pos
		pos += so.padding

		totalFactor -= f.factor
		remaining -= w
	}

	return positions
}

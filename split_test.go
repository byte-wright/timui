package timui

import (
	"testing"

	"github.com/byte-wright/expect"
)

func TestSplitFactors(t *testing.T) {
	expect.Value(t, "split factors",
		Split().Factor(0.25, 0.25, 1).calculatePositions(100),
	).ToBe([]splitRange{
		{from: 0, to: 17},
		{from: 17, to: 34},
		{from: 34, to: 100},
	})

	expect.Value(t, "split factors with pad",
		Split().Factor(0.25, 0.25, 1).Pad(3).calculatePositions(100),
	).ToBe([]splitRange{
		{from: 0, to: 16},
		{from: 19, to: 35},
		{from: 38, to: 100},
	})

	expect.Value(t, "split factors with fixed",
		Split().Factor(0.25, 0.25, 0.25).Add(0.25, 40).calculatePositions(100),
	).ToBe([]splitRange{
		{from: 0, to: 15},
		{from: 15, to: 30},
		{from: 30, to: 45},
		{from: 45, to: 100},
	})

	expect.Value(t, "split with fixed",
		Split().Fixed(10).Factor(0.5, 1).Fixed(10).Add(0.25, 40).calculatePositions(100),
	).ToBe([]splitRange{
		{from: 0, to: 10},
		{from: 10, to: 21},
		{from: 21, to: 44},
		{from: 44, to: 54},
		{from: 54, to: 100},
	})

	expect.Value(t, "split with fixed",
		Split().Fixed(1, 1, 1, 7, 1).Factor(1).Fixed(1, 1, 1).calculatePositions(27),
	).ToBe([]splitRange{
		{from: 0, to: 1},
		{from: 1, to: 2},
		{from: 2, to: 3},
		{from: 3, to: 10},
		{from: 10, to: 11},
		{from: 11, to: 24},
		{from: 24, to: 25},
		{from: 25, to: 26},
		{from: 26, to: 27},
	})

	expect.Value(t, "split case",
		Split().Fixed(12).Factor(0.3, 0.4, 0.5, 0.6, 0.5).insertFixedBetween(1).calculatePositions(200),
	).ToBe([]splitRange{
		{from: 0, to: 1},
		{from: 1, to: 13},
		{from: 13, to: 14},
		{from: 14, to: 38},
		{from: 38, to: 39},
		{from: 39, to: 70},
		{from: 70, to: 71},
		{from: 71, to: 110},
		{from: 110, to: 111},
		{from: 111, to: 158},
		{from: 158, to: 159},
		{from: 159, to: 199},
		{from: 199, to: 200},
	})
}

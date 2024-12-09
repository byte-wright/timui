package timui

import (
	"testing"

	"github.com/akabio/expect"
)

func TestSplitFactors(t *testing.T) {
	expect.Value(t, "split factors",
		Split().Factors(0.25, 0.25, 1).calculatePositions(100),
	).ToBe([]splitRange{
		{from: 0, to: 16},
		{from: 16, to: 32},
		{from: 32, to: 100},
	})

	expect.Value(t, "split factors with pad",
		Split().Factors(0.25, 0.25, 1).Pad(3).calculatePositions(100),
	).ToBe([]splitRange{
		{from: 0, to: 15},
		{from: 18, to: 33},
		{from: 36, to: 100},
	})

	expect.Value(t, "split factors with fixed",
		Split().Factors(0.25, 0.25, 0.25).Add(0.25, 40).calculatePositions(100),
	).ToBe([]splitRange{
		{from: 0, to: 15},
		{from: 15, to: 30},
		{from: 30, to: 45},
		{from: 45, to: 100},
	})

	expect.Value(t, "split with fixed",
		Split().Fixed(10).Factors(0.5, 1).Fixed(10).Add(0.25, 40).calculatePositions(100),
	).ToBe([]splitRange{
		{from: 0, to: 10},
		{from: 10, to: 21},
		{from: 21, to: 44},
		{from: 44, to: 54},
		{from: 54, to: 100},
	})
}

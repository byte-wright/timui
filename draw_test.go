package timui

import (
	"testing"
)

func BenchmarkMapBorderRune(b *testing.B) {
	acc := ' '
	for i := 0; i < b.N; i++ {
		a := mapBorderRune('╣', '|')
		b := mapBorderRune('║', 'a')
		c := mapBorderRune('╩', 'f')

		acc += a + b + c
	}
}

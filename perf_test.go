package timui

import "testing"

var (
	cnt = 0
	a   = 0
)

func WithFuncParam(f func()) {
	cnt++
	f()
	cnt--
}

func CntUp() {
	cnt++
}

func CntDown() {
	cnt--
}

func BenchmarkFuncParam(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WithFuncParam(func() {
			a += cnt
		})
	}
}

func BenchmarkEncParam(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CntUp()
		a += cnt
		CntDown()
	}
}

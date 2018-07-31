package history

import "testing"

func BenchmarkAdd(b *testing.B) {
	t := NewTree()
	b.N = 100000

	for i := 0; i < b.N; i++ {
		t.Add(string(uint64(i)))
	}
}

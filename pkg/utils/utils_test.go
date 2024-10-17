package utils

import (
	"testing"
)

func BenchmarkFileOpen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		WriteResultToFile("test.txt", "test")
	}
}

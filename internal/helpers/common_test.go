package helpers_test

import (
	"go-exposed-config-scanner/internal/helpers"
	"strings"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	arrays := make(chan string, 100)
	for i := 0; i < b.N; i++ {
		helpers.MergeURLAndPaths([]string{"http://test.com"}, []string{"1", "3"}, arrays)
	}

}

func BenchmarkTestString(b *testing.B) {

	for i := 0; i < b.N; i++ {
		x := strings.Builder{}
		x.Grow(1024)
		_ = x
	}

}

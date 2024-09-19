package helpers_test

import (
	"go-exposed-config-scanner/internal/helpers"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	arrays := make(chan string, 100)
	for i := 0; i < b.N; i++ {
		helpers.MergeURLAndPaths([]string{"http://test.com"}, []string{"1", "3"}, arrays)
	}

}

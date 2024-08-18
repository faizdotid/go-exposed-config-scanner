package helpers_test

import (
	"go-exposed-config-scanner/internal/helpers"
	"testing"
)

func BenchmarkTest(b *testing.B) {
	arrays := make([]string, 1000000)
	for i := 0; i < b.N; i++ {
		helpers.MergeURLAndPaths("http://test.com", []string{"1", "3"}, &arrays)
	}

}

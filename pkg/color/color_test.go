package color_test

import (
	"fmt"
	"go-exposed-config-scanner/pkg/color"
	"testing"
)

func TestColoring(t *testing.T) {
	t.Run("all color", func(t *testing.T) {
		for i := 30; i < 38; i++ {
			fmt.Print(string(color.Coloring("test", color.Color(i))))
			fmt.Print(" ")
		}
	})

	t.Run("all color with bold", func(t *testing.T) {
		for i := 30; i < 38; i++ {
			fmt.Print(string(color.Coloring("test", color.Color(i), color.Bold)))
			fmt.Print(" ")
		}
	})

	t.Run("all color with bold and background", func(t *testing.T) {
		for i := 30; i < 38; i++ {
			fmt.Print(string(color.Coloring("test", color.Color(i), color.Bold, color.Color(i+10))))
			fmt.Print(" ")
		}
	})

}

func BenchmarkColor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		color.Coloring("test", color.Red, color.Bold)
	}
}

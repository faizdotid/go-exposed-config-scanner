package matcher_test

import (
	"go-exposed-config-scanner/pkg/matcher"
	"regexp"
	"testing"
)

func TestMatcher(t *testing.T) {
	t.Run("TestWordMatcher", func(t *testing.T) {
		var matcher matcher.Matcher = regexp.MustCompile(`AKIA[0-9A-Z]{16}`)

		match := matcher.Match([]byte(`AKIAJGJGJGJGJGJGJGHVGVGVGVG`))

		if !match {

			t.Errorf("Expected match, got no match")
		}

	})

	t.Run("TestJsonMatcher", func(t *testing.T) {
		var matcher matcher.Matcher = matcher.NewJsonMatcher()
		match := matcher.Match([]byte(`{"key": ["value"]}`))
		if !match {
			t.Errorf("Expected match, got no match")
		}
	})

}

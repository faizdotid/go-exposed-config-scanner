package matcher_test

import (
	// "go-exposed-config-scanner/pkg/matcher"
	// "regexp"
	"fmt"
	"go-exposed-config-scanner/pkg/matcher"
	"net/http"
	"testing"
)

func TestMatcher(t *testing.T) {

	// 	var matcher matcher.Matcher = regexp.MustCompile(`AKIA[0-9A-Z]{16}`)

	// 	match := matcher.Match([]byte(`AKIAJGJGJGJGJGJGJGHVGVGVGVG`))

	// 	if !match {

	// 		t.Errorf("Expected match, got no match")
	// 	}

	// })

	// t.Run("TestJsonMatcher", func(t *testing.T) {
	// 	var matcher matcher.Matcher = matcher.NewJsonMatcher()
	// 	match := matcher.Match([]byte(`{"key": ["value"]}`))
	// 	if !match {
	// 		t.Errorf("Expected match, got no match")
	// 	}
	// })

	t.Run("TestHeaderMatcher", func(t *testing.T) {
		var h = http.Header{
			"Content-Type": []string{"application/json"},
		}
		m := matcher.NewWordMatcher("application/json")
		match := m.Match([]byte(fmt.Sprintf("%v", h)))
		if !match {
			t.Errorf("Expected match, got no match")
		}
	})

	t.Run("TestBinaryMatcher", func(t *testing.T) {
		m := matcher.NewBinaryMatcher("616263")
		match := m.Match([]byte("abc"))
		if !match {
			t.Errorf("Expected match, got no match")
		}
	})

}

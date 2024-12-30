package matcher_test

import (
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
		m := matcher.NewMatcher(
			matcher.NewWordMatcher("value"),
			[]int{123, 200},
			matcher.Headers,
		)
		resp := &http.Response{
			Header: http.Header{
				"key": []string{"value"},
			},
			StatusCode: 200,
		}
		match, _ := m.Match(resp)
		if !match {
			t.Errorf("Expected match, got no match")
		}

		t.Logf("Matched: %v", match)

	})
	t.Run("TestBinaryMatcher", func(t *testing.T) {
		m := matcher.NewBinaryMatcher("616263")
		match := m.Match([]byte("abc"))
		if !match {
			t.Errorf("Expected match, got no match")
		}
	})
}

func BenchmarkTest(b *testing.B) {
	// test word matcher
	m := matcher.NewMatcher(
		matcher.NewWordMatcher("value"),
		[]int{123, 200},
		matcher.Headers,
	)
	resp := &http.Response{
		Header: http.Header{
			"key": []string{"value"},
		},
		StatusCode: 200,
	}
	for i := 0; i < b.N; i++ {
		m.Match(resp)
	}

}

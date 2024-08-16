package matcher

type Matcher interface {
	Match([]byte) bool
}

type WordMatcher struct {
	contains string
}

type JsonMatcher struct{}

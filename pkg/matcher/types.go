package matcher

type Matcher interface {
	Match([]byte) bool
}

type WordMatcher struct {
	contains []string
}

type BinaryMatcher struct {
	contains [][]byte
}

type JSONMatcher struct{}

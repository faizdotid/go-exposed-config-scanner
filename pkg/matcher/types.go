package matcher

import (
	"net/http"
	"regexp"
)

type MatchFrom string

const (
	Body    MatchFrom = "body"
	Headers MatchFrom = "headers"
)

type Matcher struct {
	IMatch
	*StatusCodeMatch
	MatchFrom
}

type IMatch interface {
	Match([]byte) bool
}

type IMatcher interface {
	Match(*http.Response) (bool, error)
}

type WordMatch struct {
	contains []string
}

type BinaryMatch struct {
	contains [][]byte
}

type JSONMatch struct{}

type StatusCodeMatch struct {
	statusCodes map[int]struct{}
}

type RegexMatch struct {
	reg *regexp.Regexp
}

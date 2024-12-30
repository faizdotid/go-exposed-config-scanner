package matcher

import (
	"net/http"
	"regexp"
	"sync"
)

// Constants
const (
	Body    MatchFrom = "body"
	Headers MatchFrom = "headers"

	MaxBodySize int64 = 10 << 20 // 10MB
	BufferSize  int64 = 32 << 10 // 32KB
)

var (
	bufferPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, BufferSize)
			return &b
		},
	}
)

type MatchFrom string

type IMatch interface {
	Match([]byte) bool
}

type IMatcher interface {
	Match(*http.Response) (bool, error)
}

type Matcher struct {
	IMatch
	*StatusCodeMatch
	MatchFrom
}

type StatusCodeMatch struct {
	statusCodes map[int]struct{}
}

type WordMatch struct {
	contains []string
}

type BinaryMatch struct {
	contains [][]byte
}

type JSONMatch struct{}

type RegexMatch struct {
	reg *regexp.Regexp
}

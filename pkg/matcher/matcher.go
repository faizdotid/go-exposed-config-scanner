package matcher

import (
	"encoding/json"
	"strings"
)

func (m *WordMatcher) Match(data []byte) bool {
	return strings.Contains(string(data), m.contains)
}

func (m *JsonMatcher) Match(data []byte) bool {
	return json.Valid(data)
}

func NewWordMatcher(contains string) *WordMatcher {
	return &WordMatcher{contains: contains}
}

func NewJsonMatcher() *JsonMatcher {
	return &JsonMatcher{}
}

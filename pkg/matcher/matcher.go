package matcher

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"strings"
)

func (m *WordMatcher) Match(data []byte) bool {
	for _, word := range m.contains {
		if strings.Contains(string(data), word) {
			return true
		}
	}
	return false
}

func (m *JSONMatcher) Match(data []byte) bool {
	return json.Valid(data)
}

func NewWordMatcher(contains string) *WordMatcher {
	var each []string
	var words []string = strings.Split(contains, ",")
	for index, word := range words {
		curr := strings.TrimSpace(word)
		if curr[len(curr)-1] == '\\' {
			curr = curr[:len(curr)-1] + "," + strings.TrimSpace(words[index+1])
		}
		each = append(each, curr)
	}
	return &WordMatcher{contains: each}
}

func NewJSONMatcher() *JSONMatcher {
	return &JSONMatcher{}
}

func NewBinaryMatcher(contains string) *BinaryMatcher {
	var each [][]byte
	var words []string = strings.Split(contains, ",")
	for index, word := range words {
		curr := strings.TrimSpace(word)
		if curr[len(curr)-1] == '\\' {
			curr = curr[:len(curr)-1] + "," + strings.TrimSpace(words[index+1])
		}
		bin, _ := hex.DecodeString(curr)
		each = append(each, bin)
	}
	return &BinaryMatcher{contains: each}
}

func (m *BinaryMatcher) Match(data []byte) bool {
	for _, bin := range m.contains {
		if bytes.Contains(data, bin) {
			return true
		}
	}
	return false
}

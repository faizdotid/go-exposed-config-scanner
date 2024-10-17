package matcher

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
)

// STATUS CODE
func (m *StatusCodeMatch) Match(data *http.Response) bool {
	if len(m.statusCodes) == 0 || m.statusCodes == nil {
		return true
	}
	_, ok := m.statusCodes[data.StatusCode]
	return ok
}

func NewStatusCodesMatcher(statusCodes []int) *StatusCodeMatch {
	return &StatusCodeMatch{
		statusCodes: func() map[int]struct{} {
			m := make(map[int]struct{})
			for _, code := range statusCodes {
				m[code] = struct{}{}
			}
			return m
		}(),
	}
}

// BINARY
func (m *BinaryMatch) Match(data []byte) bool {
	for _, bin := range m.contains {
		if bytes.Contains(data, bin) {
			return true
		}
	}
	return false
}

func NewBinaryMatcher(contains string) *BinaryMatch {
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
	return &BinaryMatch{contains: each}
}

// WORDS
func (m *WordMatch) Match(data []byte) bool {
	for _, word := range m.contains {
		if strings.Contains(string(data), word) {
			return true
		}
	}
	return false
}

func NewWordMatcher(contains string) *WordMatch {
	var each []string
	var words []string = strings.Split(contains, ",")
	for index, word := range words {
		curr := strings.TrimSpace(word)
		if curr[len(curr)-1] == '\\' {
			curr = curr[:len(curr)-1] + "," + strings.TrimSpace(words[index+1])
		}
		each = append(each, curr)
	}
	return &WordMatch{contains: each}
}

// JSON
func (m *JSONMatch) Match(data []byte) bool {
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return false
	}
	return true
}

func NewJSONMatcher() *JSONMatch {
	return &JSONMatch{}
}

// REGEX
func (m *RegexMatch) Match(data []byte) bool {
	return m.reg.Match(data)
}

func NewRegexMatcher(value string) *RegexMatch {
	return &RegexMatch{reg: regexp.MustCompile(value)}
}

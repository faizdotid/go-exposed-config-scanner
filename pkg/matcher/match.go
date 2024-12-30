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
	if contains == "" {
		return &WordMatch{contains: nil}
	}

	words := strings.Split(contains, ",")
	each := make([]string, 0, len(words))
	var builder strings.Builder
	builder.Grow(len(contains)) // Pre-allocate builder capacity

	for i := 0; i < len(words); i++ {
		curr := strings.TrimSpace(words[i])
		if curr == "" {
			continue
		}

		if curr[len(curr)-1] == '\\' && i+1 < len(words) {
			builder.Reset()
			builder.WriteString(curr[:len(curr)-1])
			builder.WriteRune(',')
			builder.WriteString(strings.TrimSpace(words[i+1]))
			each = append(each, builder.String())
			i++
		} else {
			each = append(each, curr)
		}
	}

	return &WordMatch{contains: each}
}

// JSON
func (m *JSONMatch) Match(data []byte) bool {
	return json.Valid(data)
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

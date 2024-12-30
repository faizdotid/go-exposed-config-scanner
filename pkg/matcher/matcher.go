package matcher

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func NewMatcher(m IMatch, statusCodes []int, matchFrom MatchFrom) *Matcher {
	return &Matcher{
		IMatch:          m,
		StatusCodeMatch: NewStatusCodesMatcher(statusCodes),
		MatchFrom:       matchFrom,
	}
}

func (m *Matcher) matchBody(data *http.Response) (bool, error) {
	if data.Body == nil {
		return false, nil
	}
	defer data.Body.Close()

	bufPtr := bufferPool.Get().(*[]byte)
	buf := *bufPtr
	defer func() {
		buf = buf[:0]
		bufferPool.Put(bufPtr)
	}()

	reader := io.LimitReader(data.Body, MaxBodySize)
	body, err := io.ReadAll(reader)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %w", err)
	}

	return m.IMatch.Match(body), nil
}

func (m *Matcher) matchHeaders(resp *http.Response) (bool, error) {
	var builder strings.Builder
	for key, values := range resp.Header {
		builder.Reset()
		builder.WriteString(key)
		builder.WriteString(": ")

		for i, v := range values {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(v)
		}

		if m.IMatch.Match([]byte(builder.String())) {
			return true, nil
		}
	}
	return false, nil
}

func (m *Matcher) Match(data *http.Response) (bool, error) {
	if data == nil {
		return false, fmt.Errorf("response is nil")
	}

	if !m.StatusCodeMatch.Match(data) {
		return false, nil
	}

	matchers := map[MatchFrom]func(*http.Response) (bool, error){
		Body:    m.matchBody,
		Headers: m.matchHeaders,
	}

	if matcher, ok := matchers[m.MatchFrom]; ok {
		return matcher(data)
	}

	return false, fmt.Errorf("unsupported match type: %s", m.MatchFrom)
}

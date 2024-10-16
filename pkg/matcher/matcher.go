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
	buf, err := io.ReadAll(data.Body)
	if err != nil {
		return false, err
	}
	match := m.IMatch.Match(buf)
	return match, nil
}

func (m *Matcher) matchHeaders(data *http.Response) (bool, error) {
	for _, value := range data.Header {
		if m.IMatch.Match([]byte(strings.Join(value, ","))) {
			return true, nil
		}
	}
	return false, nil
}

func (m *Matcher) Match(data *http.Response) (bool, error) {
	if !m.StatusCodeMatch.Match(data) {
		return false, nil
	}

	switch m.MatchFrom {
	case Body:
		return m.matchBody(data)
	case Headers:
		return m.matchHeaders(data)
	default:
		return false, fmt.Errorf("invalid matchFrom value: %s", m.MatchFrom)
	}

}

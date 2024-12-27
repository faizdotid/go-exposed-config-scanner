package templates

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-exposed-config-scanner/pkg/matcher"
	"go-exposed-config-scanner/pkg/utils"
	"io"
	"net/http"
	"strings"
	"time"
)

// UnmarshalJSON is a custom unmarshaler for Template struct
func (t *Template) UnmarshalJSON(data []byte) error {
	if t.Request == nil {
		t.Request = &Request{}
	}

	var raw struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Output  string `json:"output"`
		Request struct {
			Method  string            `json:"method"`
			Timeout int               `json:"timeout"`
			Headers map[string]string `json:"headers"`
			Body    string            `json:"body"`
		}
		Match struct {
			StatusCode json.RawMessage `json:"status_code"`
			From       string          `json:"from"`
			Type       string          `json:"type"`
			Value      string          `json:"value"`
		}
		Paths []string `json:"paths"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	t.ID = raw.ID
	t.Name = raw.Name
	t.Output = raw.Output

	// Validate request components
	if err := t.validateMethodRequest(raw.Request.Method); err != nil {
		return err
	}

	if err := t.validateBodyRequest(raw.Request.Body); err != nil {
		return err
	}

	if err := t.validateTimeoutRequest(raw.Request.Timeout); err != nil {
		return err
	}

	if err := t.validatePathsRequest(raw.Paths); err != nil {
		return err
	}

	t.setHeadersRequest(raw.Request.Headers)

	// Validate matcher components
	statusCodes := t.validateStatusCodeRequest(string(raw.Match.StatusCode))

	match, err := t.validateMatchRequest(raw.Match.Type, raw.Match.Value)
	if err != nil {
		return err
	}

	matchFrom, err := t.validateMatchFrom(raw.Match.From)
	if err != nil {
		return err
	}

	t.Matcher = matcher.NewMatcher(
		match,
		statusCodes,
		matchFrom,
	)

	return nil
}

// Request validation methods
func (t *Template) validateMethodRequest(data string) error {
	validMethods := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
		"HEAD":   true,
		"PATCH":  true,
	}

	if validMethods[data] {
		t.Request.Method = data
		return nil
	}
	return errors.New("invalid method")
}

func (t *Template) validateBodyRequest(data string) error {
	if data != "" {
		t.Request.Body = io.NopCloser(strings.NewReader(data))
	} else {
		t.Request.Body = nil
	}
	return nil
}

func (t *Template) validateTimeoutRequest(data int) error {
	if data < 0 || data > 60 {
		return errors.New("timeout must be greater than 0 and less than 60")
	}

	if data == 0 {
		t.Request.Timeout = defaultHttpTimeout
	} else {
		t.Request.Timeout = time.Duration(data) * time.Second
	}
	return nil
}

func (t *Template) validatePathsRequest(data []string) error {
	if len(data) == 0 {
		return errors.New("paths must not be empty")
	}

	for _, path := range data {
		t.Paths = append(t.Paths, strings.TrimPrefix(strings.TrimSpace(path), "/"))
	}
	return nil
}

func (t *Template) setHeadersRequest(data map[string]string) {
	t.Request.Headers = make(http.Header)
	for k, v := range defaultHttpHeaders {
		t.Request.Headers[k] = v
	}
	for k, v := range data {
		t.Request.Headers.Set(k, v)
	}
}

// Matcher validation methods
func (t Template) validateStatusCodeRequest(data string) []int {
	if data == "" || data == "*" || data == "all" || data == "0" || data == "any" {
		return []int{}
	}
	return utils.ExplodeString[int](data)
}

func (t Template) validateMatchRequest(types, value string) (matcher.IMatch, error) {
	switch types {
	case "regex":
		return matcher.NewRegexMatcher(value), nil
	case "word", "words":
		return matcher.NewWordMatcher(value), nil
	case "json":
		return matcher.NewJSONMatcher(), nil
	case "binary":
		return matcher.NewBinaryMatcher(value), nil
	default:
		return nil, errors.New("invalid match type")
	}
}

func (t Template) validateMatchFrom(data string) (matcher.MatchFrom, error) {
	switch data {
	case "body":
		return matcher.Body, nil
	case "headers", "header":
		return matcher.Headers, nil
	default:
		return "", errors.New("invalid match from")
	}
}

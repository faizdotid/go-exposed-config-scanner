package templates

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-exposed-config-scanner/pkg/matcher"
	"io"
	"regexp"
	"strings"
	"time"
)

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

func (t *Template) validateMatchRequest(from, types, value string) error {
	if from != "body" && from != "headers" {
		return errors.New("invalid from value")
	}

	switch types {
	case "regex":
		t.Match = regexp.MustCompile(value)
	case "word", "words":
		t.Match = matcher.NewWordMatcher(value)
	case "json":
		t.Match = matcher.NewJSONMatcher()
	default:
		return errors.New("invalid match type")
	}
	return nil
}

func (t *Template) validateMatchFrom(data string) error {
	if data != "body" && data != "headers" {
		return errors.New("invalid matchFrom value")
	}
	t.MatchFrom = data
	return nil
}

func (t *Template) setHeadersRequest(data map[string]string) {
	t.Request.Headers = defaultHttpHeaders
	for k, v := range data {
		t.Request.Headers.Set(k, v)
	}
}

func (t *Template) validatePathsRequest(data []string) error {
	if len(data) == 0 {
		return errors.New("paths must not be empty")
	}
	t.Paths = data
	return nil
}

// custom unmarshaler for Template struct
func (t *Template) UnmarshalJSON(data []byte) error {
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
			From  string `json:"from"`
			Type  string `json:"type"`
			Value string `json:"value"`
		}
		Paths []string `json:"paths"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	t.ID = raw.ID
	t.Name = raw.Name
	t.Output = raw.Output

	// validate method
	if err := t.validateMethodRequest(raw.Request.Method); err != nil {
		return err
	}

	// validate body
	if err := t.validateBodyRequest(raw.Request.Body); err != nil {
		return err
	}

	// validate timeout
	if err := t.validateTimeoutRequest(raw.Request.Timeout); err != nil {
		return err
	}

	// set headers
	t.setHeadersRequest(raw.Request.Headers)

	// validate match
	if err := t.validateMatchRequest(raw.Match.From, raw.Match.Type, raw.Match.Value); err != nil {
		return err
	}

	// validate matchFrom
	if err := t.validateMatchFrom(raw.Match.From); err != nil {
		return err
	}

	// validate paths
	if err := t.validatePathsRequest(raw.Paths); err != nil {
		return err
	}

	return nil
}

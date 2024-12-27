package templates

import (
	"errors"
	"go-exposed-config-scanner/pkg/matcher"
	"io"
	"net/http"
	"time"
)

// Errors
var (
	ErrFolderNotFound   = errors.New("folder not found")
	ErrTemplateNotFound = errors.New("template not found")
	ErrInvalidTemplate  = errors.New("invalid template")
)

// Default HTTP configurations
var (
	defaultHttpTimeout = 7 * time.Second
	defaultHttpHeaders = http.Header{
		"User-Agent": []string{
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 " +
				"(KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		},
	}
)

// Template represents a configuration template for scanning
type Template struct {
	ID        string
	Name      string
	Output    string
	Paths     []string
	MatchFrom string
	Matcher   matcher.IMatcher
	Request   *Request
}

// Request contains HTTP request configuration
type Request struct {
	Method  string
	Body    io.ReadCloser
	Timeout time.Duration
	Headers http.Header
}

// Templates is a collection of Template pointers
type Templates []*Template

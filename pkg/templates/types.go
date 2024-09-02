package templates

import (
	"errors"
	"go-exposed-config-scanner/pkg/matcher"
	"io"
	"net/http"
	"time"
)

type Template struct {
	ID        string
	Name      string
	Output    string
	Request   *Request
	Match     matcher.Matcher
	MatchFrom string
	Paths     []string
}

type Request struct {
	Method  string
	Timeout time.Duration
	Headers http.Header
	Body    io.ReadCloser
}

type Templates []*Template

var (
	// err variables
	ErrFolderNotFound   = errors.New("folder not found")
	ErrTemplateNotFound = errors.New("template not found")
	ErrInvalidTemplate  = errors.New("invalid template")

	// default values
	defaultHttpTimeout = 7 * time.Second
	defaultHttpHeaders = http.Header{
		"User-Agent": []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3"},
	}
)

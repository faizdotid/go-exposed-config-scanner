package templates

import (
	"errors"
	"go-exposed-config-scanner/pkg/matcher"
)

type Template struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Output    string          `json:"output"`
	MatchFrom string          `json:"match_from"`
	Matcher   matcher.Matcher `json:"match"`
	Paths     []string        `json:"paths"`
}

// Type for slice of templates
type Templates []*Template

var (
	// Err variable
	ErrTemplateNotFound = errors.New("template not found")
	ErrTemplateFormat   = errors.New("template format error")
)

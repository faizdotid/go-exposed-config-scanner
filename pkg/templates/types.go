package templates

import (
	"errors"
	"go-exposed-config-scanner/pkg/validators"
)

type Template struct {
	ID        string
	Name      string
	Paths     []string
	Output    string
	Validator validators.ValidatorFunction
	MatchFrom string
}

// Type for slice of templates
type Templates []*Template

var (
	// Err variable
	ErrTemplateNotFound = errors.New("template not found")
	ErrTemplateFormat   = errors.New("template format error")
)

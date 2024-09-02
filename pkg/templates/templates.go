package templates

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// return template by id
func (t Templates) GetTemplateByID(id string) (*Template, error) {
	for _, template := range t {
		if template.ID == id {
			return template, nil
		}
	}
	return nil, ErrTemplateNotFound
}

// load template from a specific directory
func (t *Templates) LoadTemplate(dir string) error {
	if dir == "" {
		dir = "templates"
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" && !strings.Contains(path, "example.json") && !strings.Contains(path, "README.md") {
			if err := t.readFileTemplate(path); err != nil {
				return fmt.Errorf("error reading template file %s: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking the path %s: %w", dir, err)
	}
	return nil
}

// read file template
func (t *Templates) readFileTemplate(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}

	template := &Template{
		Request: &Request{},
	}
	if err := json.Unmarshal(file, template); err != nil {
		return fmt.Errorf("error unmarshalling JSON from file %s: %w", path, err)
	}

	*t = append(*t, template)
	return nil
}

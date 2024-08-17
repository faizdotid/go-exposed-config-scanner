package templates

import (
	"encoding/json"
	"fmt"
	"go-exposed-config-scanner/pkg/matcher"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Load template config from directory
func (t *Templates) LoadTemplates(dir string) error {
	if dir == "" {
		dir = "configs"
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %w", path, err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" && !strings.Contains(path, "example.json") {
			if err := t.readFromFile(path); err != nil {
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

// readFromFile reads a template from a JSON file and appends it to the Templates slice.
func (t *Templates) readFromFile(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}

	var template Template
	if err := json.Unmarshal(file, &template); err != nil {
		return fmt.Errorf("error unmarshalling JSON from file %s: %w", path, err)
	}

	*t = append(*t, &template)
	return nil
}

// GetTemplateByID returns a template by its ID.
func (t Templates) GetTemplateByID(id string) (*Template, error) {
	for _, template := range t {
		if template.ID == id {
			return template, nil
		}
	}
	return nil, ErrTemplateNotFound
}

// UnmarshalJSON is a custom JSON unmarshalling function for Template.
func (t *Template) UnmarshalJSON(data []byte) error {
	var raw struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Output    string `json:"output"`
		MatchFrom string `json:"match_from"`
		Matcher   struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"match"`
		Paths []string `json:"paths"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	t.ID = raw.ID
	t.Name = raw.Name
	t.Output = raw.Output
	t.Paths = raw.Paths
	t.MatchFrom = raw.MatchFrom

	// Handle the Matcher based on its type
	switch raw.Matcher.Type {
	case "regex":
		regex, err := regexp.Compile(raw.Matcher.Value)
		if err != nil {
			return fmt.Errorf("invalid regexp in template: %w", err)
		}
		t.Matcher = regex
	case "word":
		t.Matcher = matcher.NewWordMatcher(raw.Matcher.Value)
	case "json":
		t.Matcher = matcher.NewJsonMatcher()
	default:
		return fmt.Errorf("unknown matcher type %q in template", raw.Matcher.Type)
	}

	return nil
}

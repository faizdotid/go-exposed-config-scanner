package templates

import (
	"encoding/json"
	"fmt"
	"go-exposed-config-scanner/pkg/validators"
	"os"
	"path/filepath"
)


func (t *Templates) LoadTemplates(dir string) error {
	if dir == "" {
		dir = "configs"
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			if err := t.readFromFile(path); err != nil {
				return fmt.Errorf("error reading template file %s: %w", path, err)
			}
		}
		return nil
	})
}

func (t *Templates) registerNewTemplate(id, name, output, validatorName, matchFrom string, paths []string) error {
	validator, err := validators.GetValidator(validatorName)
	if err != nil {
		return fmt.Errorf("error getting validator: %w", err)
	}
	*t = append(*t, &Template{ID: id, Name: name, Paths: paths, Output: output, Validator: validator, MatchFrom: matchFrom})
	return nil
}

func (t *Templates) readFromFile(path string) error {

	file, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(file, &data); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	id, ok := data["id"].(string)
	if !ok {
		return fmt.Errorf("invalid template format for ID %s", id)
	}

	name, ok := data["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name for template ID %s", id)
	}

	output, ok := data["output"].(string)
	if !ok {
		return fmt.Errorf("invalid output for template ID %s", id)
	}

	validatorName, ok := data["validator"].(string)
	if !ok {
		return fmt.Errorf("invalid validator for template ID %s", id)
	}

	paths, ok := data["paths"].([]interface{})
	if !ok {
		return fmt.Errorf("invalid paths for template ID %s", id)
	}
	matchFrom, ok := data["match_from"].(string)
	if !ok {
		matchFrom = ""
	}
	if matchFrom != "" {
		matchFrom = "body"
	}


	stringPaths := make([]string, len(paths))
	for i, p := range paths {
		stringPaths[i], ok = p.(string)
		if !ok {
			return fmt.Errorf("invalid path format for template ID %s", id)
		}
	}

	if err := t.registerNewTemplate(id, name, output, validatorName, matchFrom, stringPaths); err != nil {
		return fmt.Errorf("error registering template %s: %w", id, err)
	}

	return nil
}

func (t Templates) GetTemplateByID(id string) (*Template, error) {
	for _, template := range t {
		if template.ID == id {
			return template, nil
		}
	}
	return nil, ErrTemplateNotFound
}


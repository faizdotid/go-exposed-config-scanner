package helpers

import (
	"fmt"
	"go-exposed-config-scanner/pkg/templates"
	"strings"
)

func MergeURLAndPaths(urls []string, paths []string, result chan<- string) {
	initialSize := 1024

	builder := strings.Builder{}

	for _, url := range urls {
		baseURL := strings.TrimSpace(url)

		requiredCapacity := len(baseURL)
		if !strings.Contains(baseURL, "http") {
			requiredCapacity += len("http://")
		}
		if !strings.HasSuffix(baseURL, "/") {
			requiredCapacity++
		}

		maxPathLen := 0
		for _, path := range paths {
			if len(path) > maxPathLen {
				maxPathLen = len(path)
			}
		}

		totalRequired := requiredCapacity + maxPathLen

		if totalRequired > initialSize {
			builder.Grow(totalRequired)
		} else {
			builder.Grow(initialSize)
		}

		if !strings.Contains(baseURL, "http") {
			builder.WriteString("http://")
			builder.WriteString(baseURL)
		} else {
			builder.WriteString(baseURL)
		}

		if !strings.HasSuffix(builder.String(), "/") {
			builder.WriteString("/")
		}

		baseURLWithSlash := builder.String()
		builder.Reset()

		for _, path := range paths {
			builder.WriteString(baseURLWithSlash)
			builder.WriteString(strings.TrimSpace(path))
			result <- builder.String()
			builder.Reset()
		}
	}
}

// ParseArgsForTemplates parses the given ID or retrieves all templates based on the 'all' flag.
func ParseArgsForTemplates(id string, all bool, t *templates.Templates) ([]*templates.Template, error) {
	if all {
		return *t, nil
	}

	if id == "" && !all {
		return nil, fmt.Errorf("you must provide a template ID or use the -all flag")
	}
	if !strings.Contains(id, ",") {
		template, err := t.GetTemplateByID(id)
		if err != nil {
			return nil, err
		}
		return []*templates.Template{template}, nil
	}

	ids := strings.Split(id, ",")
	var selectedTemplates []*templates.Template
	for _, id := range ids {
		template, err := t.GetTemplateByID(id)
		if err != nil {
			return nil, err
		}
		selectedTemplates = append(selectedTemplates, template)
	}

	return selectedTemplates, nil
}

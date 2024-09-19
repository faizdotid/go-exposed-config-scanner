package helpers

import (
	"fmt"
	"go-exposed-config-scanner/pkg/templates"
	// "reflect"
	"strings"
)

// MergeURLAndPaths merges a base URL with a list of paths and appends the resulting URLs to the provided slice.
func MergeURLAndPaths(urls []string, paths []string, result chan<- string) {
	// ptr := reflect.ValueOf(result)

	// if ptr.Kind() != reflect.Ptr || ptr.Elem().Kind() != reflect.Slice {
	// 	panic("result must be a pointer to a slice")
	// }
	for _, url := range urls {
		url = strings.TrimSpace(url)
		if !strings.Contains(url, "http") {
			url = "http://" + url
		}

		for _, path := range paths {
			path = strings.TrimSpace(path)
			if !strings.HasSuffix(url, "/") {
				url += "/"
			}

			finalURL := url + path
			result <- finalURL
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

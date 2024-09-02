package templates_test

import (
	// "go-scanner/pkg/templates"
	"encoding/json"
	"go-exposed-config-scanner/pkg/templates"
	"testing"
)

func TestLoadTemplates(t *testing.T) {
	var mytemp templates.Templates
	err := mytemp.LoadTemplate("../../configs")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	t.Logf("Templates: %v", mytemp[0])
}

func BenchmarkTemplates(b *testing.B) {
	var mytemp templates.Templates
	err := mytemp.LoadTemplate("../../configs")
	if err != nil {
		b.Errorf("Expected no error, got %v", err)
	}
}

func TestJson(t *testing.T) {
	jsonraw := `{"x": "value"}`
	var s struct {
		Key string `json:"key"`
	}
	err := json.Unmarshal([]byte(jsonraw), &s)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if s.Key == "" {
		t.Errorf("Expected value, got empty")
	}
	t.Logf("Struct: %T", s.Key)

}

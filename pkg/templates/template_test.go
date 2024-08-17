package templates_test

import (
	// "go-scanner/pkg/templates"
	"go-exposed-config-scanner/pkg/templates"
	"testing"
)

func TestLoadTemplates(t *testing.T) {
	var mytemp templates.Templates
	err := mytemp.LoadTemplates("../../configs")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	t.Logf("Templates: %v", mytemp[0])
}

func BenchmarkTemplates(b *testing.B) {
	var mytemp templates.Templates
	err := mytemp.LoadTemplates("../../configs")
	if err != nil {
		b.Errorf("Expected no error, got %v", err)
	}
}

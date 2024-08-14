package validators_test

import (
	"go-exposed-config-scanner/pkg/validators"
	"testing"
)

func TestNotFoundValidator(t *testing.T) {
	_, err := validators.GetValidator("notfound")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestFoundValidator(t *testing.T) {
	v, err := validators.GetValidator("is_json")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	t.Logf("Found validator: %v", v)
}

func TestAddValidator(t *testing.T) {
	validators.RegisterValidator("test_validator", func(data []byte) (bool, error) {
		return true, nil
	})
	v, err := validators.GetValidator("test_validator")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	t.Logf("Found validator: %v", v)
}

func TestAll(t *testing.T) {
	t.Run("TestIsPhpIinfo", func(t *testing.T) {
		valid, err := validators.IsPhpIinfo([]byte("phpinfo()"))
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !valid {
			t.Errorf("Expected valid, got invalid")
		}
	})
	t.Run("TestIsWordressBackupConfig", func(t *testing.T) {
		valid, err := validators.IsWordressBackupConfig([]byte("$table_prefix"))
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !valid {
			t.Errorf("Expected valid, got invalid")
		}
	})

	t.Run("TestIsYiiDebugger", func(t *testing.T) {
		valid, err := validators.IsYiiDebugger([]byte("Yii Debugger"))
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !valid {
			t.Errorf("Expected valid, got invalid")
		}
	})
	t.Run("TestIsDotEnv", func(t *testing.T) {
		valid, err := validators.IsDotEnv([]byte("APP_KEY"))
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if !valid {
			t.Errorf("Expected valid, got invalid")
		}
	})
}

func BenchmarkFunction(b *testing.B) {
	for i := 0; i < b.N; i++ {
		validators.IsPhpIinfo([]byte("phpinfo()"))
		validators.IsWordressBackupConfig([]byte("$table_prefix"))
		validators.IsYiiDebugger([]byte("Yii Debugger"))
		validators.IsDotEnv([]byte("APP_KEY"))
		validators.IsAwsCredential([]byte("AKIAXXXXXXXXXXXXXX"))
		validators.IsJson([]byte(`{"key": "value"}`))
		validators.IsJavascript([]byte("text/javascript"))

		// All OK Have 0/1 allocs and 0/1 B/op
	}
}
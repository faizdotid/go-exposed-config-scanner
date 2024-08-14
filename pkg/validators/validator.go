package validators

import (
	"encoding/json"
	"strings"
)

func IsJson(data []byte) (bool, error) {
	return json.Valid(data), nil
}

func IsJavascript(data []byte) (bool, error) {
	content := string(data)
	return strings.Contains(content, "text/javascript") ||
		strings.Contains(content, "application/javascript") ||
		strings.Contains(content, "application/x-javascript"), nil
}

func IsAwsCredential(data []byte) (bool, error) {
	return RegexAwsCredential.Match(data), nil
}

func IsDotEnv(data []byte) (bool, error) {
	return strings.Contains(string(data), "APP_KEY"), nil
}

func IsPhpIinfo(data []byte) (bool, error) {
	return strings.Contains(string(data), "phpinfo()"), nil
}

func IsYiiDebugger(data []byte) (bool, error) {
	return strings.Contains(string(data), "Yii Debugger"), nil
}

func IsWordressBackupConfig(data []byte) (bool, error) {
	return strings.Contains(string(data), "$table_prefix"), nil
}


func RegisterValidator(name string, validator ValidatorFunction) {
	validators = append(validators, map[string]ValidatorFunction{name: validator})
}

// get validator function by name
func GetValidator(name string) (ValidatorFunction, error) {
	for _, v := range validators {
		if _, ok := v[name]; ok {
			return v[name], nil
		}
	}
	return nil, ErrValidatorNotFound
}

// Init built-in validators
func init() {
	RegisterValidator("is_json", IsJson)
	RegisterValidator("is_javascript", IsJavascript)
	RegisterValidator("is_aws_credential", IsAwsCredential)
	RegisterValidator("is_dotenv", IsDotEnv)
	RegisterValidator("is_php_info", IsPhpIinfo)
	RegisterValidator("is_yii_debugger", IsYiiDebugger)
	RegisterValidator("is_wordpress_backup_config", IsWordressBackupConfig)
}

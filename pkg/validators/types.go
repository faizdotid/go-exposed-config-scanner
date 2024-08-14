package validators

import (
	"errors"
	"regexp"
)

// type for validator function
type ValidatorFunction func([]byte) (bool, error)

var (
	validators           []map[string]ValidatorFunction
	ErrValidatorNotFound = errors.New("validator not found")

	RegexAwsCredential = regexp.MustCompile(`AKIA[0-9A-Z]{16}`)
)

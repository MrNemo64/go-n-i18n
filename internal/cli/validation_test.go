package cli_test

import (
	"testing"

	"github.com/MrNemo64/go-n-i18n/internal/cli"
)

func TestKeyValidator_Validate(t *testing.T) {
	t.Parallel()
	validator := cli.KeyValidator()
	notIfFalse := func(v bool) string {
		if v {
			return ""
		} else {
			return "not "
		}
	}

	tests := map[string]bool{
		"first":                true,
		"one.two":              true,
		"one_string.then2":     true,
		"A":                    true,
		"1_string":             false,
		"one_string.2_strings": false,
		"_then":                false,
		"":                     false,
		".":                    false,
	}

	for key, isValid := range tests {
		validatorResult := validator.IsValidKey(key)
		if isValid != validatorResult {
			t.Errorf("The key '%s' is %svalid, but the validator returned %v", key, notIfFalse(isValid), validatorResult)
		}
	}
}

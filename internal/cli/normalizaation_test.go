package cli_test

import (
	"testing"

	"github.com/MrNemo64/go-n-i18n/internal/cli"
)

func TestKeyNormalizer_Normalize(t *testing.T) {
	t.Parallel()
	normalizer := cli.KeyNormalizer()

	tests := map[string]string{
		"first":            "First",
		"one.two":          "One.Two",
		"one_string.then2": "OneString.Then2",
		"A":                "A",
	}

	for key, normalized := range tests {
		normalizatorResult := normalizer.Normalize(key)
		if normalized != normalizatorResult {
			t.Errorf("The key '%s' was normalized to '%s' but '%s' was expected", key, normalizatorResult, normalized)
		}
	}
}

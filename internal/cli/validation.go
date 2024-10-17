package cli

import "regexp"

type keyValidator struct {
	re *regexp.Regexp
}

func KeyValidator() keyValidator {
	return keyValidator{
		re: regexp.MustCompile(`^([a-zA-Z][a-zA-Z0-9_-]*)(\.[a-zA-Z][a-zA-Z0-9_-]*)*$`),
	}
}

func (kv keyValidator) IsValidKey(key string) bool {
	return kv.re.MatchString(key)
}

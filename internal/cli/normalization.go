package cli

import "strings"

type keyNormalizer struct{}

func KeyNormalizer() keyNormalizer {
	return keyNormalizer{}
}

func (keyNormalizer) Normalize(key string) string {
	// TODO this function can be heavily optimized (probably)
	parts := strings.Split(key, ".")

	sb := strings.Builder{}
	for i, part := range parts {
		sb.Reset()

		sb.WriteString(strings.ToUpper(part[:1]))
		for j := 1; j < len(part); j++ {
			if part[j] == '-' || part[j] == '_' {
				if j == len(part)-1 {
					break // last char is a _ or a - so just ignore it
				}
				j++
				sb.WriteString(strings.ToUpper(part[j : j+1]))
			} else {
				sb.WriteString(part[j : j+1])
			}
		}

		parts[i] = sb.String()
	}

	return strings.Join(parts, ".")
}

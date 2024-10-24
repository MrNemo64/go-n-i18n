package types

import "strings"

func copyAndAdd[T any](ori []T, new ...T) []T {
	ret := make([]T, len(ori))
	copy(ret, ori)
	return append(ret, new...)
}

func ResolveFullPath(parent *MessageBag, child string) []string {
	if parent == nil {
		if child == "" {
			return []string{}
		}
		return []string{child}
	}
	return copyAndAdd(parent.Path(), child)
}

func PathAsStr(path []string) string {
	return strings.Join(path, ".")
}

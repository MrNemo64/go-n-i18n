package assert

import "fmt"

func NonNil(v any, msg string) {
	if v == nil {
		panic(fmt.Errorf("expected non nil but was nil: %s", msg))
	}
}

func NoError(err error) {
	if err != nil {
		panic(fmt.Errorf("reached unreachable error: %w", err))
	}
}

func Has[T comparable](arr []T, el T) {
	for _, e := range arr {
		if e == el {
			return
		}
	}
	panic(fmt.Errorf("expected slice %+v to have the element %+v", arr, el))
}

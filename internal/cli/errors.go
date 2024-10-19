package cli

import "fmt"

type Error struct {
	msg  string
	args []any
}

func (err Error) WithArgs(arg ...any) Error {
	copy := err
	copy.args = arg
	return copy
}

func (err Error) Error() string {
	return fmt.Sprintf(err.msg, err.args...)
}

func (err Error) Is(other error) bool {
	casted, ok := other.(Error)
	if !ok {
		return false
	}
	return casted.msg == err.msg
}

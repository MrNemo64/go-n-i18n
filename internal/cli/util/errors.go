package util

import (
	"fmt"
)

type Error struct {
	msg  string
	args []any
}

func MakeError(msg string) Error {
	return Error{msg: msg}
}

func (err Error) WithArgs(arg ...any) Error {
	copy := err
	copy.args = arg
	return copy
}

func (err Error) Error() string {
	return fmt.Errorf(err.msg, err.args...).Error()
}

func (err Error) Is(other error) bool {
	casted, ok := other.(Error)
	return ok && casted.msg == err.msg
}

func (err Error) Unwrap() []error {
	if err.args == nil {
		return []error{}
	}
	var errors []error
	for _, arg := range err.args {
		if e, ok := arg.(error); ok {
			errors = append(errors, e)
		}
	}
	return errors
}

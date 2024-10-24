package types

type ValueString struct {
	message string
}

func NewStringLiteralValue(message string) *ValueString {
	return &ValueString{message: message}
}

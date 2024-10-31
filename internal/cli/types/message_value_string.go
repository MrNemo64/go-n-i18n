package types

import "strings"

type ValueString struct {
	message string
}

func NewStringLiteralValue(message string) *ValueString {
	return &ValueString{message: message}
}

func (*ValueString) multilineMarker()              {}
func (*ValueString) conditionableMarker()          {}
func (s *ValueString) AsValueString() *ValueString { return s }
func (*ValueString) AsValueParametrized() *ValueParametrized {
	panic("called AsValueParametrized on a ValueString")
}
func (*ValueString) AsMultiline() *ValueMultiline {
	panic("called AsMultiline on a ValueString")
}
func (*ValueString) AsConditional() *ValueConditional {
	panic("called AsConditional on a ValueString")
}
func (s *ValueString) Escaped(quote string) string {
	return strings.ReplaceAll(s.message, quote, "\\\"")
}

package types

import "github.com/MrNemo64/go-n-i18n/internal/cli/util"

var (
	ErrInsuficientLines util.Error = util.MakeError("there must be at least one line")
)

type Multilineable interface {
	multilineMarker()
}

type ValueMultiline struct {
	Lines []Multilineable
}

func NewMultilineValue(lines []Multilineable) (*ValueMultiline, error) {
	if len(lines) == 0 {
		return nil, ErrInsuficientLines
	}
	return &ValueMultiline{Lines: lines}, nil
}

func (*ValueMultiline) conditionableMarker()           {}
func (s *ValueMultiline) AsMultiline() *ValueMultiline { return s }
func (*ValueMultiline) AsValueString() *ValueString {
	panic("called AsValueString on a ValueMultiline")
}
func (*ValueMultiline) AsValueParametrized() *ValueParametrized {
	panic("called AsValueParametrized on a ValueMultiline")
}
func (*ValueMultiline) AsConditional() *ValueConditional {
	panic("called AsConditional on a ValueMultiline")
}

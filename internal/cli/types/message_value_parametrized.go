package types

import (
	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
)

var (
	ErrInvalidAmountOfTextSegmentsAndArguments util.Error = util.MakeError("the amount of text segments (%d) is not the amount of arguments (%d) + 1")
)

type ValueParametrized struct {
	TextSegments []*ValueString
	Args         []*UsedArgument
}

type UsedArgument struct {
	Argument *MessageArgument
	Format   string
}

func NewParametrizedStringValue(textSegments []*ValueString, args []*UsedArgument) (*ValueParametrized, error) {
	if len(textSegments) != len(args)+1 {
		return nil, ErrInvalidAmountOfTextSegmentsAndArguments.WithArgs(len(textSegments), len(args))
	}
	return &ValueParametrized{
		TextSegments: textSegments,
		Args:         args,
	}, nil
}

func (*ValueParametrized) AsValueString() *ValueString {
	panic("called AsValueString on a ParametrizedString")
}
func (s *ValueParametrized) AsValueParametrized() *ValueParametrized { return s }

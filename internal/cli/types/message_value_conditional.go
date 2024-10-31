package types

type Conditionable interface {
	conditionableMarker()
}

type Condition struct {
	Condition string
	Value     Conditionable
}

type ValueConditional struct {
	Conditions []Condition
	Else       Conditionable
}

func NewConditionalValue(conditions []Condition, elseCondition Conditionable) (*ValueConditional, error) {
	return &ValueConditional{
		Conditions: conditions,
		Else:       elseCondition,
	}, nil
}

func (c *ValueConditional) AsConditional() *ValueConditional { return c }
func (*ValueConditional) AsValueString() *ValueString {
	panic("called AsValueString on a ValueConditional")
}
func (*ValueConditional) AsValueParametrized() *ValueParametrized {
	panic("called AsValueParametrized on a ValueConditional")
}
func (*ValueConditional) AsMultiline() *ValueMultiline {
	panic("called AsMultiline on a ValueConditional")
}

package types

type MessageValue interface {
	AsValueString() *ValueString
	AsValueParametrized() *ValueParametrized
	AsMultiline() *ValueMultiline
}

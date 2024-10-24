package types

import "github.com/MrNemo64/go-n-i18n/internal/cli/util"

var (
	ErrCannotMergeValues util.Error = util.MakeError("cannot merge message value %d with %d")
)

type MessageValue interface {
	AsValueString() *ValueString
}

package types

import "github.com/MrNemo64/go-n-i18n/internal/cli/util"

var (
	ErrArgumentCollition util.Error = util.MakeError("argumtn %s has a type colition: %s != %s")
)

type MessageArgument struct {
	Name string
	Type string
}

type ArgumentList struct {
	args []MessageArgument
}

func NewArgumentList() *ArgumentList {
	return &ArgumentList{args: make([]MessageArgument, 0)}
}

func (l *ArgumentList) AddArgument(arg MessageArgument) error {
	existing, found := l.GetArgument(arg.Name)
	if !found {
		l.args = append(l.args, arg)
		return nil
	}
	if existing.Type != arg.Type {
		return ErrArgumentCollition.WithArgs(arg.Name, arg.Type, existing.Type)
	}
	return nil
}

func (l *ArgumentList) GetArgument(name string) (*MessageArgument, bool) {
	for i := range l.args {
		if l.args[i].Name == name {
			return &l.args[i], true
		}
	}
	return nil, false
}

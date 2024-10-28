package types

import (
	"errors"

	"github.com/MrNemo64/go-n-i18n/internal/cli/assert"
	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
)

var (
	ErrArgumentCollition util.Error = util.MakeError("argument %s has a type colition: %s != %s")
	ErrMergeArgumentList            = util.MakeError("could not merge argument lists: %w")
)

type ArgumentType struct {
	Name          string
	Aliases       []string
	Type          string
	DefaultFormat string
	IsUnknown     bool
}

func (t *ArgumentType) Is(name string) bool {
	if name == t.Name {
		return true
	}
	for _, alias := range t.Aliases {
		if alias == name {
			return true
		}
	}
	return false
}

type ArgumentProvider struct {
	types []*ArgumentType
}

func NewArgumentProvider() *ArgumentProvider {
	p := &ArgumentProvider{}
	p.Register(&ArgumentType{
		Name:          "any",
		Aliases:       []string{"any", "unknown"},
		Type:          "any",
		DefaultFormat: "v",
		IsUnknown:     true,
	})
	p.Register(&ArgumentType{
		Name:          "string",
		Aliases:       []string{"string", "str"},
		Type:          "string",
		DefaultFormat: "s",
	})
	p.Register(&ArgumentType{
		Name:          "boolean",
		Aliases:       []string{"boolean", "bool"},
		Type:          "bool",
		DefaultFormat: "t",
	})
	p.Register(&ArgumentType{
		Name:          "integer",
		Aliases:       []string{"integer", "int"},
		Type:          "int",
		DefaultFormat: "d",
	})
	p.Register(&ArgumentType{
		Name:          "float64",
		Aliases:       []string{"float64", "f64", "f", "double"},
		Type:          "float64",
		DefaultFormat: "g",
	})
	return p
}

func (p *ArgumentProvider) UnknwonType() *ArgumentType {
	return p.types[0]
}

func (p *ArgumentProvider) Register(arg *ArgumentType) bool {
	if _, found := p.FindArgument(arg.Name); found {
		return false
	}
	for _, alias := range arg.Aliases {
		if _, found := p.FindArgument(alias); found {
			return false
		}
	}
	p.types = append(p.types, arg)
	return true
}

func (p *ArgumentProvider) FindArgument(name string) (*ArgumentType, bool) {
	for _, arg := range p.types {
		if arg.Is(name) {
			return arg, true
		}
	}
	return nil, false
}

func (p *ArgumentProvider) FindArgumentOrUnknwonType(name string) *ArgumentType {
	if arg, found := p.FindArgument(name); found {
		return arg
	}
	return p.UnknwonType()
}

type MessageArgument struct {
	Name string
	Type *ArgumentType
}

type ArgumentList struct {
	Args []*MessageArgument
}

func NewArgumentList() *ArgumentList {
	return &ArgumentList{Args: make([]*MessageArgument, 0)}
}

func (l *ArgumentList) Merge(other *ArgumentList) error {
	var errs []error
	for _, arg := range other.Args {
		if _, err := l.AddArgument(arg); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return ErrMergeArgumentList.WithArgs(errors.Join(errs...))
}

func (l *ArgumentList) AddArgument(arg *MessageArgument) (*MessageArgument, error) {
	assert.NonNil(arg.Type, "arg.Type")
	existing, found := l.GetArgument(arg.Name)
	if !found {
		l.Args = append(l.Args, arg)
		return arg, nil
	}

	if arg.Type.IsUnknown {
		return existing, nil
	}
	if existing.Type.IsUnknown {
		existing.Type = arg.Type // arg.Type != unknown and existing.Type == unknown -> specify the type
		return existing, nil
	}
	if existing.Type != arg.Type {
		return nil, ErrArgumentCollition.WithArgs(arg.Name, arg.Type.Name, existing.Type.Name)
	}
	return existing, nil
}

func (l *ArgumentList) GetArgument(name string) (*MessageArgument, bool) {
	for i := range l.Args {
		if l.Args[i].Name == name {
			return l.Args[i], true
		}
	}
	return nil, false
}

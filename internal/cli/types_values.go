package cli

import "strings"

type StringLiteral struct {
	String string
}

func (l *StringLiteral) Part() string { return l.String }

type ArgumentInstance struct {
	Name   string
	Format string
}

func (a *ArgumentInstance) Part() string { return "%" + a.Format }

type MessagePart interface {
	Part() string
}

type ParametrizedString struct {
	Parts []MessagePart
}

func (p *ParametrizedString) String() string {
	return strings.Join(Map(p.Parts, func(t *MessagePart) string { return (*t).Part() }), "")
}

func (p *ParametrizedString) Args() []*ArgumentInstance {
	args := []*ArgumentInstance{}
	for _, v := range p.Parts {
		if arg, ok := v.(*ArgumentInstance); ok {
			args = append(args, arg)
		}
	}
	return args
}

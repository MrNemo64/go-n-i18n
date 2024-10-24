package types

type FunctionDefinition interface {
	Name() string
	ReturnType() string
}

type MessageInstanceFunctionDefinition struct {
	source *MessageInstance
	name   string
}

func (f *MessageInstanceFunctionDefinition) Source() *MessageInstance { return f.source }
func (f *MessageInstanceFunctionDefinition) Name() string             { return f.name }
func (*MessageInstanceFunctionDefinition) ReturnType() string         { return "string" }

type MessageBagFunctionDefinition struct {
	source *MessageBag
	name   string
}

func (f *MessageBagFunctionDefinition) Source() *MessageBag { return f.source }
func (f *MessageBagFunctionDefinition) Name() string        { return f.name }
func (f *MessageBagFunctionDefinition) ReturnType() string  { return f.name }

type InterfaceDefinition struct {
	Source     *MessageBag
	Name       string
	Functions  []FunctionDefinition
	Interfaces []*InterfaceDefinition
}

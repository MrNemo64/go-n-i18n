package cli

var (
	ErrEntryNotFound                          Error = Error{msg: "entry '%s' not found"}
	ErrEntryParentIsNotBag                          = Error{msg: "could not get entry '%s' because the entry '%s' is not a message bag (kind: %d)"}
	ErrEntryNotFoundBecauseParentDoesNotExist       = Error{msg: "could not get entry '%s' because the entry '%s' does not exist"}
	ErrAddedEntryIsNotTheSameKind                   = Error{msg: "the added entry has kind %d but there is already an entry with kind %d"}

	ErrLiteralMessageRedefinition = Error{msg: `the literal message with key '%s' was already defined as "%s" but it got redefined as "%s"`}
)

type MessageEntryKind int

const (
	MessageEntryLiteral MessageEntryKind = iota
	MessageEntryParametrized
	MessageEntryBag
)

type MessageEntryValue interface {
}

type MessageEntry interface {
	Key() string
	Kind() MessageEntryKind
	Languages() *Set[string]
	EnsureAllLanguagesPresent(defLang string, languages []string) bool
	FullPath() []string
	AssignParent(*MessageEntryMessageBag)

	AsLiteral() *MessageEntryLiteralString
	AsParametrized() *MessageEntryParametrizedString
	AsBag() *MessageEntryMessageBag
}

type MessageFunctionDefinition struct {
	name    string
	Message *MessageEntryLiteralString
}

func (f *MessageFunctionDefinition) Name() string       { return f.name }
func (f *MessageFunctionDefinition) ReturnType() string { return "string" }

type BagFunctionDefinition struct {
	name       string
	returnType string
}

func (f *BagFunctionDefinition) Name() string       { return f.name }
func (f *BagFunctionDefinition) ReturnType() string { return f.returnType }

type FunctionDeclaration interface {
	Name() string
	ReturnType() string
}

type InterfaceDefinition struct {
	Name       string
	Functions  []FunctionDeclaration
	Interfaces []*InterfaceDefinition
}

package cli

import (
	"errors"
	"strings"
)

type MessageArgument struct {
	Name   string
	Type   *ArgumentType
	Format string
}

type MessageEntryParametrizedString struct {
	parent    *MessageEntryMessageBag
	key       string
	message   map[string]*ParametrizedString // language tag -> message
	arguments []*MessageArgument
}

func (MessageEntryParametrizedString) With(key string, message map[string]*ParametrizedString) *MessageEntryParametrizedString {
	return &MessageEntryParametrizedString{
		key:     key,
		message: message,
	}
}

func (*MessageEntryParametrizedString) Kind() MessageEntryKind                            { return MessageEntryParametrized }
func (p *MessageEntryParametrizedString) Key() string                                     { return p.key }
func (p *MessageEntryParametrizedString) AsParametrized() *MessageEntryParametrizedString { return p }
func (p *MessageEntryParametrizedString) Args() []*MessageArgument                        { return p.arguments }
func (p *MessageEntryParametrizedString) Lang(lang string) *ParametrizedString {
	return p.message[strings.ReplaceAll(lang, "_", "-")]
}

func (p *MessageEntryParametrizedString) FullPath() []string {
	return append(p.parent.FullPath(), p.key)
}
func (p *MessageEntryParametrizedString) FullPathAsStr() string {
	return strings.Join(p.FullPath(), ".")
}

func (p *MessageEntryParametrizedString) AssignParent(parent *MessageEntryMessageBag) {
	p.parent = parent
}

func (p *MessageEntryParametrizedString) Languages() *Set[string] {
	set := NewSet[string]()
	for k := range p.message {
		set.Add(k)
	}
	return set
}

func (p *MessageEntryParametrizedString) EnsureAllLanguagesPresent(defLang string, languages []string) bool {
	hadToFill := false
	for _, lang := range languages {
		if _, hasIt := p.message[lang]; !hasIt {
			p.message[lang] = p.message[defLang]
			hadToFill = true
		}
	}
	return hadToFill
}

func (p *MessageEntryParametrizedString) GetArgument(name string) (*MessageArgument, error) {
	for _, arg := range p.arguments {
		if arg.Name == name {
			return arg, nil
		}
	}
	return nil, ErrArgumentNotFound.WithArgs(name)
}

func (p *MessageEntryParametrizedString) Merge(other *MessageEntryParametrizedString) error {
	for lang, message := range other.message {
		if existingMsg, found := p.message[lang]; found {
			return ErrParametreizedMessageRedefinition.WithArgs(p.key, existingMsg, message)
		}
		p.message[lang] = message
	}
	for _, arg := range other.arguments {
		if err := p.AddArgument(arg.Name, arg.Type.Type, arg.Format); err != nil {
			return err
		}
	}
	return nil
}

func (p *MessageEntryParametrizedString) AddArgument(name, kind, format string) error {
	existing, err := p.GetArgument(name)
	if errors.Is(err, ErrArgumentNotFound) { // new argument
		argKind := FindArgumentType(kind)
		if argKind == nil { // the type is not specified yet
			argKind = AnyKind()
		}
		if format == "" {
			format = argKind.DefaultFormat
		}
		if !argKind.IsValidFormat(format) {
			return ErrInvalidArgumentFormat.WithArgs(kind, format)
		}
		p.arguments = append(p.arguments, &MessageArgument{
			Name:   name,
			Type:   argKind,
			Format: format,
		})
		return nil
	}

	newType := FindArgumentType(kind)
	if kind != "" && newType == nil {
		return ErrUnknwonArgumentType.WithArgs(kind)
	}
	if kind == "" || newType.IsAny { // already seen this argument but no type information is provided
		return nil
	}
	if existing.Type.IsAny { // we still don't know the type of this argument
		if !newType.IsAny { // a type is now specified
			// TODO this will ignore if the format was previously specified
			existing.Type = newType
			if format == "" {
				format = newType.DefaultFormat
			}
			if !newType.IsValidFormat(format) {
				return ErrInvalidArgumentFormat.WithArgs(kind, format)
			}
			existing.Format = format
		}
	} else { // we already know the type of this arg but it was specified again
		if existing.Type != newType || existing.Format != format {
			return ErrArgumentAlreadySpecified.WithArgs(name)
		}
	}
	return nil
}

func (p *MessageEntryParametrizedString) DefineFunction(namer MessageEntryNamer) *ParametrizedFunctionDefinition {
	return &ParametrizedFunctionDefinition{name: namer.FunctionName(p), Message: p, Args: p.arguments}
}

func (*MessageEntryParametrizedString) AsBag() *MessageEntryMessageBag {
	panic("called AsBag in a MessageEntryParametrizedString")
}

func (*MessageEntryParametrizedString) AsLiteral() *MessageEntryLiteralString {
	panic("called AsLiteral in a MessageEntryMessageBag")
}

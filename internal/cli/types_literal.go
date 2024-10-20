package cli

import "strings"

type MessageEntryLiteralString struct {
	parent  *MessageEntryMessageBag
	key     string
	message map[string]string // language tag -> message
}

func (MessageEntryLiteralString) With(key string, message map[string]string) *MessageEntryLiteralString {
	return &MessageEntryLiteralString{
		key:     key,
		message: message,
	}
}
func (*MessageEntryLiteralString) Kind() MessageEntryKind                        { return MessageEntryLiteral }
func (l *MessageEntryLiteralString) Key() string                                 { return l.key }
func (l *MessageEntryLiteralString) AsLiteral() *MessageEntryLiteralString       { return l }
func (l *MessageEntryLiteralString) FullPath() []string                          { return append(l.parent.FullPath(), l.key) }
func (l *MessageEntryLiteralString) AssignParent(parent *MessageEntryMessageBag) { l.parent = parent }
func (l *MessageEntryLiteralString) Message(lang string) string {
	return l.message[strings.ReplaceAll(lang, "_", "-")]
}
func (*MessageEntryLiteralString) AsBag() *MessageEntryMessageBag {
	panic("called AsBag in a MessageEntryLiteralString")
}
func (*MessageEntryLiteralString) AsParametrized() *MessageEntryParametrizedString {
	panic("called AsParametrized in a MessageEntryMessageBag")
}

func (l *MessageEntryLiteralString) Merge(other *MessageEntryLiteralString) error {
	for lang, message := range other.message {
		if existingMsg, found := l.message[lang]; found {
			return ErrLiteralMessageRedefinition.WithArgs(l.key, existingMsg, message)
		}
		l.message[lang] = message
	}
	return nil
}

func (l *MessageEntryLiteralString) Languages() *Set[string] {
	set := NewSet[string]()
	for k := range l.message {
		set.Add(k)
	}
	return set
}

func (l *MessageEntryLiteralString) EnsureAllLanguagesPresent(defLang string, languages []string) bool {
	hadToFill := false
	for _, lang := range languages {
		if _, hasIt := l.message[lang]; !hasIt {
			l.message[lang] = l.message[defLang]
			hadToFill = true
		}
	}
	return hadToFill
}

func (l *MessageEntryLiteralString) DefineFunction(namer MessageEntryNamer) *MessageFunctionDefinition {
	return &MessageFunctionDefinition{name: namer.FunctionName(l), Message: l}
}

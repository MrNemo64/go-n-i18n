package cli

type MessageEntryParametrizedString struct {
	parent  *MessageEntryMessageBag
	key     string
	message map[string]string // language tag -> message
}

func (MessageEntryParametrizedString) With(key string, message map[string]string) *MessageEntryParametrizedString {
	panic("todo")
}

func (*MessageEntryParametrizedString) Kind() MessageEntryKind                            { return MessageEntryParametrized }
func (l *MessageEntryParametrizedString) Key() string                                     { return l.key }
func (l *MessageEntryParametrizedString) AsParametrized() *MessageEntryParametrizedString { return l }
func (l *MessageEntryParametrizedString) FullPath() []string {
	return append(l.parent.FullPath(), l.key)
}
func (l *MessageEntryParametrizedString) AssignParent(parent *MessageEntryMessageBag) {
	l.parent = parent
}

func (*MessageEntryParametrizedString) AsBag() *MessageEntryMessageBag {
	panic("called AsBag in a MessageEntryParametrizedString")
}
func (*MessageEntryParametrizedString) AsLiteral() *MessageEntryLiteralString {
	panic("called AsLiteral in a MessageEntryMessageBag")
}
func (l *MessageEntryParametrizedString) Languages() *Set[string] {
	set := NewSet[string]()
	for k := range l.message {
		set.Add(k)
	}
	return set
}
func (l *MessageEntryParametrizedString) EnsureAllLanguagesPresent(defLang string, languages []string) bool {
	hadToFill := false
	for _, lang := range languages {
		if _, hasIt := l.message[lang]; !hasIt {
			l.message[lang] = l.message[defLang]
			hadToFill = true
		}
	}
	return hadToFill
}

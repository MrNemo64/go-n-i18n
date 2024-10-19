package cli

import (
	"errors"
	"sort"
	"strings"
)

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
	AsBag() *MessageEntryMessageBag
}

type MessageEntryMessageBag struct {
	parent  *MessageEntryMessageBag
	key     string
	entries []MessageEntry
}

func (MessageEntryMessageBag) With(key string, entries []MessageEntry) *MessageEntryMessageBag {
	return &MessageEntryMessageBag{key: key, entries: entries}
}
func (*MessageEntryMessageBag) Kind() MessageEntryKind                        { return MessageEntryBag }
func (b *MessageEntryMessageBag) Key() string                                 { return b.key }
func (b *MessageEntryMessageBag) AsBag() *MessageEntryMessageBag              { return b }
func (b *MessageEntryMessageBag) AssignParent(parent *MessageEntryMessageBag) { b.parent = parent }
func (b *MessageEntryMessageBag) IsRoot() bool                                { return b.key == "" }
func (b *MessageEntryMessageBag) FullPath() []string {
	if b.IsRoot() {
		return []string{}
	}
	return append(b.parent.FullPath(), b.key)
}
func (*MessageEntryMessageBag) AsLiteral() *MessageEntryLiteralString {
	panic("called AsLiteral in a MessageEntryMessageBag")
}

func (b *MessageEntryMessageBag) GetEntry(key string) (MessageEntry, error) {
	for _, e := range b.entries {
		if e.Key() == key {
			return e, nil
		}
	}
	return nil, ErrEntryNotFound.WithArgs(key)
}

func (b *MessageEntryMessageBag) FindEntry(path ...string) (MessageEntry, error) {
	if len(path) == 0 {
		return nil, ErrEntryNotFound.WithArgs("")
	}
	if len(path) == 1 {
		return b.GetEntry(path[0])
	}
	entry, err := b.GetEntry(path[0])
	if err != nil {
		return nil, err
	}
	if entry.Kind() != MessageEntryBag {
		return nil, ErrEntryParentIsNotBag.WithArgs(strings.Join(path, "."), path[0], entry.Kind)
	}
	for i := 1; i < len(path)-1; i++ {
		entry, err = entry.AsBag().GetEntry(path[i])
		if err != nil {
			return nil, ErrEntryNotFoundBecauseParentDoesNotExist.WithArgs(strings.Join(path, "."), strings.Join(path[:i+1], "."))
		}
		if entry.Kind() != MessageEntryBag {
			return nil, ErrEntryParentIsNotBag.WithArgs(strings.Join(path, "."), strings.Join(path[:i+1], "."), entry.Kind)
		}
	}
	entry, err = entry.AsBag().GetEntry(path[len(path)-1])
	if err != nil {
		return nil, ErrEntryNotFound.WithArgs(strings.Join(path, "."))
	}
	return entry, nil
}

func (b *MessageEntryMessageBag) Languages() *Set[string] {
	if len(b.entries) == 0 {
		return NewSet[string]()
	}
	set := b.entries[0].Languages()
	for i := 1; i < len(b.entries); i++ {
		set.AddAll(b.entries[i].Languages())
	}
	return set
}

func (b *MessageEntryMessageBag) FindOrCreateChildBag(path ...string) (*MessageEntryMessageBag, error) {
	actual := b
	for i := 0; i < len(path); i++ {
		found, err := actual.GetEntry(path[i])
		if errors.Is(err, ErrEntryNotFound) {
			new := &MessageEntryMessageBag{
				key:     path[i],
				entries: make([]MessageEntry, 0),
				parent:  actual,
			}
			actual.entries = append(actual.entries, new)
			actual = new
			continue
		}
		if err != nil {
			return nil, err
		}
		if found.Kind() != MessageEntryBag {
			return nil, ErrEntryParentIsNotBag.WithArgs(strings.Join(path, "."), strings.Join(path[:i+1], "."), found.Kind())
		}
		actual = found.AsBag()
	}
	return actual, nil
}

func (b *MessageEntryMessageBag) AddEntries(entries ...MessageEntry) error {
	for _, entry := range entries {
		existing, err := b.GetEntry(entry.Key())
		if errors.Is(err, ErrEntryNotFound) {
			b.entries = append(b.entries, entry)
			entry.AssignParent(b)
			sort.Slice(b.entries, func(i, j int) bool {
				return b.entries[i].Key() < b.entries[j].Key()
			})
			continue
		}
		if existing.Kind() != entry.Kind() {
			return ErrAddedEntryIsNotTheSameKind.WithArgs(entry.Kind(), existing.Kind())
		}
		switch existing.Kind() {
		case MessageEntryLiteral:
			if err := existing.AsLiteral().Merge(entry.AsLiteral()); err != nil {
				return err
			}
		case MessageEntryBag:
			existing.AsBag().AddEntries(entry.AsBag().entries...)
		}
	}
	return nil
}

func (b *MessageEntryMessageBag) RemoveEntriesWithoutLang(lang string) []MessageEntry {
	var removed []MessageEntry
	var remaining []MessageEntry
	for _, entry := range b.entries {
		if entry.Kind() == MessageEntryBag {
			if len(entry.AsBag().entries) > 0 {
				remaining = append(remaining, entry)
			} else {
				removed = append(removed, entry)
			}
		} else if entry.Kind() == MessageEntryLiteral {
			if entry.Languages().Contains(lang) {
				remaining = append(remaining, entry)
			} else {
				removed = append(removed, entry)
			}
		}
	}
	b.entries = remaining
	return removed
}

func (b *MessageEntryMessageBag) EnsureAllLanguagesPresent(defLang string, languages []string) bool {
	hadToFill := false
	for _, entry := range b.entries {
		if entry.EnsureAllLanguagesPresent(defLang, languages) {
			hadToFill = true
		}
	}
	return hadToFill
}

func (b *MessageEntryMessageBag) DefineInterface(namer MessageEntryNamer) *InterfaceDefinition {
	definition := &InterfaceDefinition{Name: namer.InterfaceName(b)}
	for _, entry := range b.entries {
		switch entry.Kind() {
		case MessageEntryLiteral:
			definition.Functions = append(definition.Functions, entry.AsLiteral().DefineFunction(namer))
		case MessageEntryBag:
			inner := entry.AsBag().DefineInterface(namer)
			definition.Functions = append(definition.Functions, &BagFunctionDefinition{
				name:       namer.FunctionName(entry),
				returnType: inner.Name,
			})
			definition.Interfaces = append(definition.Interfaces, inner)
		}
	}
	return definition
}

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

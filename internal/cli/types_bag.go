package cli

import (
	"errors"
	"sort"
	"strings"
)

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
func (*MessageEntryMessageBag) AsParametrized() *MessageEntryParametrizedString {
	panic("called AsParametrized in a MessageEntryMessageBag")
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
		case MessageEntryParametrized:
			if err := existing.AsParametrized().Merge(entry.AsParametrized()); err != nil {
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

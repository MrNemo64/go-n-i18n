package types

import (
	"fmt"

	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
)

type MessageBag struct {
	messageEntry
	children []MessageEntry
	Name     string
}

var (
	ErrParentIsNotBag          util.Error = util.MakeError("could not make or get bag entry %s because %s is not a bag")
	ErrAddedEntryIsNotSameType            = util.MakeError("the entry to add %s is of kind %d but there is already an entry of type %d")
	ErrInvalidName                        = util.MakeError("the name '%s' does not follow the allowed format (^[a-zA-Z][a-zA-Z0-9_-]*$)")
)

func IsValidName(name string) bool { return ValidKey.MatchString(name) }
func CheckName(name string) error {
	if !IsValidName(name) {
		return ErrInvalidName.WithArgs(name)
	}
	return nil
}

func NewMessageBag(key string) (*MessageBag, error) {
	if !IsValidKey(key) {
		return nil, ErrInvalidKey.WithArgs(key)
	}
	return &MessageBag{
		messageEntry: messageEntry{
			key: key,
		},
		children: make([]MessageEntry, 0),
	}, nil
}

func MakeRoot() *MessageBag {
	return &MessageBag{
		messageEntry: messageEntry{
			key: "",
		},
		children: make([]MessageEntry, 0),
	}
}
func (m *MessageBag) AsBag() *MessageBag         { return m }
func (*MessageBag) AsInstance() *MessageInstance { panic("called AsInstance on a bag") }
func (m *MessageBag) IsRoot() bool               { return m.key == "" }
func (*MessageBag) Type() MessageEntryType       { return MessageEntryBag }
func (*MessageBag) IsBag() bool                  { return true }
func (*MessageBag) IsInstance() bool             { return false }
func (m *MessageBag) Children() []MessageEntry   { return m.children }

func (b *MessageBag) GetEntry(key string) (MessageEntry, bool) {
	for _, e := range b.children {
		if e.Key() == key {
			return e, true
		}
	}
	return nil, false
}

func (m *MessageBag) FindOrCreateChildBag(path ...string) (*MessageBag, error) {
	if len(path) == 0 {
		return m, nil
	}
	actual := m
	for i := 0; i < len(path); i++ {
		found, ok := actual.GetEntry(path[i])
		if !ok {
			child, err := NewMessageBag(path[i])
			if err != nil {
				return nil, ErrCreateEntry.WithArgs(path[i], err)
			}
			actual.AddChildren(child)
			actual = child
			continue
		}
		if !found.IsBag() {
			return nil, ErrParentIsNotBag.WithArgs(PathAsStr(path), PathAsStr(path[:i+1]))
		}
		actual = found.AsBag()
	}
	return actual, nil
}

func (m *MessageBag) AddChildren(children ...MessageEntry) error {
	for _, child := range children {
		if child.IsBag() && child.AsBag().IsRoot() {
			if err := m.AddChildren(child.AsBag().children...); err != nil {
				return err
			}
			continue
		}

		existing, found := m.GetEntry(child.Key())
		if !found {
			child.AssignParent(m)
			m.children = append(m.children, child)
			continue
		}
		if existing.Type() != child.Type() {
			return ErrAddedEntryIsNotSameType.WithArgs(child.Key(), child.Type(), existing.Type())
		}
		switch existing.Type() {
		case MessageEntryBag:
			if err := existing.AsBag().AddChildren(child.AsBag().children...); err != nil {
				return err
			}
		case MessageEntryInstance:
			if err := existing.AsInstance().Merge(child.AsInstance()); err != nil {
				return err
			}
		default:
			panic(fmt.Errorf("unknown message entry type %d", existing.Type()))
		}
	}
	return nil
}

func (m *MessageBag) RemoveEntriesWithoutLang(lang string) []MessageEntry {
	var removed []MessageEntry
	var remaining []MessageEntry
	for _, child := range m.children {
		switch child.Type() {
		case MessageEntryBag:
			removed = append(removed, child.AsBag().RemoveEntriesWithoutLang(lang)...)
			if len(child.AsBag().children) > 0 {
				remaining = append(remaining, child)
			} else {
				removed = append(removed, child)
			}
		case MessageEntryInstance:
			if child.Languages().Contains(lang) {
				remaining = append(remaining, child)
			} else {
				removed = append(removed, child)
			}
		default:
			panic(fmt.Errorf("unknown message entry type %d", child.Type()))
		}
	}
	m.children = remaining
	return removed
}

func (m *MessageBag) MustHaveAllLangs(langs []string, defLang string) map[string][]string {
	ret := make(map[string][]string)
	for _, child := range m.children {
		util.MergeIntoA(ret, child.MustHaveAllLangs(langs, defLang), func(v1, v2 *[]string) []string { return append(*v1, *v2...) })
	}
	return ret
}

func (m *MessageBag) Languages() *util.Set[string] {
	set := util.NewSet[string]()
	for _, child := range m.children {
		set.AddAll(child.Languages())
	}
	return set
}

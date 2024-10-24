package types

import (
	"fmt"

	"github.com/MrNemo64/go-n-i18n/internal/cli/assert"
	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
)

var (
	ErrMessageRedefinition util.Error = util.MakeError("the message %s already is defined for %s but it got redefined")
)

type MessageInstance struct {
	messageEntry
	message map[string]MessageValue
}

func NewMessageInstance(key string) (*MessageInstance, error) {
	if !IsValidKey(key) {
		return nil, ErrInvalidKey.WithArgs(key)
	}
	return &MessageInstance{
		messageEntry: messageEntry{
			key: key,
		},
		message: make(map[string]MessageValue),
	}, nil
}

func (*MessageInstance) AsBag() *MessageBag             { panic("called AsBag on an instance") }
func (m *MessageInstance) AsInstance() *MessageInstance { return m }
func (*MessageInstance) Type() MessageEntryType         { return MessageEntryInstance }
func (*MessageInstance) IsBag() bool                    { return false }
func (*MessageInstance) IsInstance() bool               { return true }
func (m *MessageInstance) Message(lang string) (MessageValue, bool) {
	v, f := m.message[lang]
	return v, f
}

func (m *MessageInstance) AddLanguage(lang string, message MessageValue) error {
	assert.NonNil(message, "message")
	if _, found := m.message[lang]; found {
		return ErrMessageRedefinition.WithArgs(m.PathAsStr(), lang)
	}
	m.message[lang] = message
	return nil
}

func (m *MessageInstance) Merge(other *MessageInstance) error {
	panic("todo")
}

func (m *MessageInstance) Languages() *util.Set[string] {
	langs := util.NewSet[string]()
	for key := range m.message {
		langs.Add(key)
	}
	return langs
}

func (m *MessageInstance) MustHaveAllLangs(langs []string, defLang string) map[string][]string {
	defMsg, found := m.message[defLang]
	if !found {
		panic(fmt.Errorf("called MustHaveAllLangs with default lang %s but it is not present in the languages %+v", defLang, m.Languages().Get()))
	}
	missing := make(map[string][]string)
	path := m.PathAsStr()
	for _, lang := range langs {
		if _, hasIt := m.message[lang]; !hasIt {
			m.message[lang] = defMsg
			missing[lang] = []string{path}
		}
	}
	return missing
}

func (m *MessageInstance) DefineFunction(namer MessageEntryNamer) *MessageInstanceFunctionDefinition {
	return &MessageInstanceFunctionDefinition{
		source: m,
		name:   namer.FunctionName(m),
	}
}

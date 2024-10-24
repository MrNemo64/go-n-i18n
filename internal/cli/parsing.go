package cli

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/iancoleman/orderedmap"
)

var (
	ErrParsingNext                         Error = Error{msg: "could get next file to parse: %w"}
	ErrParsingIO                                 = Error{msg: "could not read contents of file %s: %w"}
	ErrParsingUnmarshal                          = Error{msg: "could not unmarshal contents of file %s: %w"}
	ErrInvalidKeyName                            = Error{msg: "invalid key '%s' in file %s. The key does not follow the allowed patter"}
	ErrKeyIsConditionalButValueIsNotObject       = Error{msg: "invalid key '%s': has the ? prefix so it's a conditional key but the value is not an object: %v"}
	ErrCouldNotAddEntry                          = Error{msg: "could not add %s entry %s: %w"}
	ErrCouldNotAddArg                            = Error{msg: "could not add argument {%s:%s:%s}: %w"}
	ErrUnknownEntryType                          = Error{msg: "could not identify the type of entry for %s: %+v in file %s"}
)

type MessagesParser struct {
	validKey          *regexp.Regexp
	argumentExtractor *regexp.Regexp
	wc                *WarningsCollector
}

func ParseJson(walker DirWalker, wc *WarningsCollector) (*MessageEntryMessageBag, error) {
	parser := &MessagesParser{
		validKey:          regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]*$"),
		argumentExtractor: regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*)(?::([a-zA-Z0-9_]+))?(?::([a-zA-Z0-9_%.]+))?\}`),
		wc:                wc,
	}
	return parser.scanMessagesInDir(walker)
}

func (m *MessagesParser) scanMessagesInDir(walker DirWalker) (*MessageEntryMessageBag, error) {
	root := &MessageEntryMessageBag{
		key:     "",
		entries: make([]MessageEntry, 0),
		parent:  nil,
	}
	for {
		file, err := walker.Next()
		if err == ErrNoMoreFiles {
			return root, nil
		}
		if err != nil {
			return nil, ErrParsingNext.WithArgs(err)
		}
		content, err := file.ReadContents()
		if err != nil {
			return nil, ErrParsingIO.WithArgs(file.FullPath(), err)
		}

		entries := orderedmap.New()
		if err := json.Unmarshal(content, entries); err != nil {
			return nil, ErrParsingUnmarshal.WithArgs(file.FullPath(), err)
		}

		dest, err := root.FindOrCreateChildBag(file.Path()...)
		if err != nil {
			return nil, err
		}
		// if any MessagesParser method returns an error, it should be instantly returned
		// errors returned by MessagesParser mean we cannot continue
		// errors returned by methods used by MessagesParser are added to the warnings collector
		// as these do not completly stop the process. The caller will decide what to do with these errors
		m.parseGroupOfMessages(entries, dest, file)
	}
}

func (m *MessagesParser) parseGroupOfMessages(entries *orderedmap.OrderedMap, dest *MessageEntryMessageBag, file FileEntry) {
	keys := entries.Keys()
	for _, key := range keys {
		value, found := entries.Get(key)
		if !found {
			panic(fmt.Sprintf("the ordered map is missing the key '%s', this is a bug in the github.com/iancoleman/orderedmap library. File: %s", key, file.FullPath()))
		}

		if strings.HasPrefix(key, "?") { // is conditional message?
			key = key[1:]
			if !m.validKey.MatchString(key) {
				m.wc.AddWarning(ErrInvalidKeyName.WithArgs(strings.Join(copySlice(dest.FullPath(), key), "."), file.FullPath()))
				continue
			}
			conditions, ok := value.(*orderedmap.OrderedMap)
			if !ok {
				m.wc.AddWarning(ErrKeyIsConditionalButValueIsNotObject.WithArgs(key, value))
				continue
			}
			newEntry := m.parseConditionalMessage(key, conditions)
			if err := dest.AddEntries(newEntry); err != nil {
				m.wc.AddWarning(ErrCouldNotAddEntry.WithArgs("conditional", strings.Join(copySlice(dest.FullPath(), newEntry.Key()), "."), err))
				continue
			}
		} else {
			if !m.validKey.MatchString(key) {
				m.wc.AddWarning(ErrInvalidKeyName.WithArgs(strings.Join(copySlice(dest.FullPath(), key), "."), file.FullPath()))
				continue
			}
			if innerEntries, ok := value.(orderedmap.OrderedMap); ok {
				if _, isParametrized := innerEntries.Get("_args"); isParametrized {
					newEntry := m.parseParametrizedWithSpecifiedArgsString(key, &innerEntries, file.Language())
					if err := dest.AddEntries(newEntry); err != nil {
						m.wc.AddWarning(ErrCouldNotAddEntry.WithArgs("parametrized", strings.Join(copySlice(dest.FullPath(), newEntry.Key()), "."), err))
						continue
					}
				} else {
					newEntries := m.parseNestedEntries(key, &innerEntries, file)
					if err := dest.AddEntries(newEntries); err != nil {
						m.wc.AddWarning(ErrCouldNotAddEntry.WithArgs("bag", strings.Join(copySlice(dest.FullPath(), newEntries.Key()), "."), err))
						continue
					}
				}
			} else if stringValue, ok := value.(string); ok {
				if m.argumentExtractor.MatchString(stringValue) {
					newEntry, ok := m.parseParametrizedLiteralString(key, stringValue, file.Language())
					if !ok {
						continue
					}
					if err := dest.AddEntries(newEntry); err != nil {
						m.wc.AddWarning(ErrCouldNotAddEntry.WithArgs("parametrized", strings.Join(copySlice(dest.FullPath(), newEntry.Key()), "."), err))
						continue
					}
				} else {
					newEntry := m.parseLiteralString(key, stringValue, file.Language())
					if err := dest.AddEntries(newEntry); err != nil {
						m.wc.AddWarning(ErrCouldNotAddEntry.WithArgs("literal", strings.Join(copySlice(dest.FullPath(), newEntry.Key()), "."), err))
						continue
					}
				}
			} else {
				m.wc.AddWarning(ErrUnknownEntryType.WithArgs(strings.Join(copySlice(dest.FullPath(), key), "."), value, file.FullPath()))
			}
		}
	}
}

func (m *MessagesParser) parseLiteralString(key, message, lang string) *MessageEntryLiteralString {
	return &MessageEntryLiteralString{
		key: key,
		message: map[string]string{
			lang: message,
		},
	}
}

func (m *MessagesParser) extractArgs(message string) {

}

func (m *MessagesParser) parseParametrizedLiteralString(key, message, lang string) (*MessageEntryParametrizedString, bool) {
	args := m.argumentExtractor.FindAllStringSubmatch(message, -1)
	p := &MessageEntryParametrizedString{
		key: key,
		message: map[string]string{
			lang: message,
		},
	}
	ok := true
	for _, arg := range args {
		if err := p.AddArgument(arg[1], arg[2], arg[3]); err != nil {
			m.wc.AddWarning(ErrCouldNotAddArg.WithArgs(arg[1], arg[2], arg[3], err))
			ok = false
		}
	}
	return p, ok
}

func (m *MessagesParser) parseParametrizedWithSpecifiedArgsString(key string, entries *orderedmap.OrderedMap, lang string) *MessageEntryParametrizedString {
	// TODO
	panic("not done")
}

func (m *MessagesParser) parseNestedEntries(key string, entries *orderedmap.OrderedMap, file FileEntry) *MessageEntryMessageBag {
	new := &MessageEntryMessageBag{
		key:     key,
		entries: make([]MessageEntry, 0),
	}
	m.parseGroupOfMessages(entries, new, file)
	return new
}

func (m *MessagesParser) parseConditionalMessage(key string, value *orderedmap.OrderedMap) MessageEntry {
	// TODO
	panic("not done")
}

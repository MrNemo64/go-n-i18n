package cli

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/iancoleman/orderedmap"
)

type MessagesParser struct {
	validKey     *regexp.Regexp
	hasArguments *regexp.Regexp
}

func ParseJson(walker DirWalker) (*MessageEntryMessageBag, error) {
	parser := MessagesParser{
		validKey:     regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]*$"),
		hasArguments: regexp.MustCompile("{.*?}"),
	}
	return parser.scanMessagesInDir(walker)
}

func (m MessagesParser) scanMessagesInDir(walker DirWalker) (*MessageEntryMessageBag, error) {
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
			return nil, fmt.Errorf("could get next file to parse: %w", err)
		}
		content, err := file.ReadContents()
		if err != nil {
			return nil, fmt.Errorf("could not read contents of file %s: %w", file.FullPath(), err)
		}

		entries := orderedmap.New()
		if err := json.Unmarshal(content, entries); err != nil {
			return nil, fmt.Errorf("could not unmarshal contents of file %s: %w", file.FullPath(), err)
		}

		if len(file.Path()) == 0 { // file is in root
			if err := m.parseGroupOfMessages(entries, root, file); err != nil {
				return nil, err
			}
		} else {
			dest, err := root.FindOrCreateChildBag(file.Path()...)
			if err != nil {
				return nil, err
			}
			if err := m.parseGroupOfMessages(entries, dest, file); err != nil {
				return nil, err
			}
		}
	}
}

func (m MessagesParser) parseGroupOfMessages(entries *orderedmap.OrderedMap, dest *MessageEntryMessageBag, file FileEntry) error {
	keys := entries.Keys()

	for _, key := range keys {
		value, found := entries.Get(key)
		if !found {
			panic(fmt.Sprintf("the ordered map is missing the key '%s', this is a bug in the github.com/iancoleman/orderedmap library. File: %s", key, file.FullPath()))
		}

		if strings.HasSuffix(key, "?") {
			key = key[:len(key)-1]
			if !m.validKey.MatchString(key) {
				return fmt.Errorf("invalid key '%s' in file %s. The key does not follow the allowed patter", key, file.FullPath())
			}
			conditions, ok := value.(*orderedmap.OrderedMap)
			if !ok {
				return fmt.Errorf("invalid key '%s': has the ? suffix so it's a conditional key but the value is not an object: %v", key, value)
			}
			if err := dest.AddEntries(m.parseConditionalMessage(key, conditions)); err != nil {
				return fmt.Errorf("could not add conditional entry: %w", err)
			}
		} else {
			if !m.validKey.MatchString(key) {
				return fmt.Errorf("invalid key '%s' in file %s. The key does not follow the allowed patter", key, file.FullPath())
			}
			if innerEntries, ok := value.(orderedmap.OrderedMap); ok {
				newEntries, err := m.parseNestedEntries(key, &innerEntries, file)
				if err != nil {
					return err
				}
				if err := dest.AddEntries(newEntries); err != nil {
					return fmt.Errorf("could not add bag entry: %w", err)
				}
			} else if stringValue, ok := value.(string); ok {
				if m.hasArguments.MatchString(stringValue) {
					// TODO
					//if err := dest.AddEntries(m.parseParametrizedString(stringValue)); err != nil {
					//	return fmt.Errorf("could not add parametrized entry: %w", err)
					//}
					if err := dest.AddEntries(m.parseLiteralString(key, stringValue, file.Language())); err != nil {
						return fmt.Errorf("could not add literal entry: %w", err)
					}
				} else {
					if err := dest.AddEntries(m.parseLiteralString(key, stringValue, file.Language())); err != nil {
						return fmt.Errorf("could not add literal entry: %w", err)
					}
				}
			} else {
				return fmt.Errorf("could not identify the type of entry for %s: %v in file %s", key, value, file.FullPath())
			}
		}
	}
	return nil
}

func (MessagesParser) parseLiteralString(key, message, lang string) *MessageEntryLiteralString {
	return &MessageEntryLiteralString{
		key: key,
		message: map[string]string{
			lang: message,
		},
	}
}

func (MessagesParser) parseParametrizedString(message string) MessageEntry {
	// TODO
	panic("not done")
}

func (m MessagesParser) parseNestedEntries(key string, entries *orderedmap.OrderedMap, file FileEntry) (*MessageEntryMessageBag, error) {
	new := &MessageEntryMessageBag{
		key:     key,
		entries: make([]MessageEntry, 0),
	}
	return new, m.parseGroupOfMessages(entries, new, file)
}

func (MessagesParser) parseConditionalMessage(key string, value *orderedmap.OrderedMap) MessageEntry {
	// TODO
	panic("not done")
}

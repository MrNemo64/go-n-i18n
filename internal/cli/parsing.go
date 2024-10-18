package cli

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/iancoleman/orderedmap"
)

type MessagesParser struct {
	log          *slog.Logger
	validKey     *regexp.Regexp
	hasArguments *regexp.Regexp
}

func ParseJson(walker DirWalker, log *slog.Logger) {
	parser := MessagesParser{
		log:          log,
		validKey:     regexp.MustCompile("^[a-zA-Z][a-zA-Z_-]*$"),
		hasArguments: regexp.MustCompile("{.*?}"),
	}
	parser.scanMessagesInDir(walker)
}

func (m MessagesParser) scanMessagesInDir(walker DirWalker) error {
	for {
		file, err := walker.Next()
		if err == ErrNoMoreFiles {
			return nil
		}
		if err != nil {
			return fmt.Errorf("could get next file to parse: %w", err)
		}
		content, err := file.ReadContents()
		if err != nil {
			return fmt.Errorf("could not read contents of file %s: %w", file.FullPath, err)
		}

		entries := orderedmap.New()
		if err := json.Unmarshal(content, entries); err != nil {
			return fmt.Errorf("could not unmarshal contents of file %s: %w", file.FullPath, err)
		}

		if err := m.parseGroupOfMessages(entries, file); err != nil {
			return err
		}

	}
}

func (m MessagesParser) parseGroupOfMessages(entries *orderedmap.OrderedMap, file *FileEntry) error {
	keys := entries.Keys()
	for _, key := range keys {
		value, found := entries.Get(key)
		if !found {
			panic(fmt.Sprintf("the ordered map is missing the key '%s', this is a bug in the github.com/iancoleman/orderedmap library. File: %s", key, file.FullPath))
		}

		if strings.HasSuffix(key, "?") {
			key = key[:len(key)-1]
			if !m.validKey.MatchString(key) {
				return fmt.Errorf("invalid key '%s' in file %s. The key does not follow the allowed patter", key, file.FullPath)
			}
			conditions, ok := value.(*orderedmap.OrderedMap)
			if !ok {
				return fmt.Errorf("invalid key '%s': has the ? suffix so it's a conditional key but the value is not an object: %v", key, value)
			}
			m.parseConditionalMessage(key, conditions)
		} else {
			if !m.validKey.MatchString(key) {
				return fmt.Errorf("invalid key '%s' in file %s. The key does not follow the allowed patter", key, file.FullPath)
			}
			if innerEntries, ok := value.(*orderedmap.OrderedMap); ok {
				m.parseNestedEntries(innerEntries)
			} else if stringValue, ok := value.(string); ok {
				if m.hasArguments.MatchString(stringValue) {
					m.parseParametrizedString(stringValue)
				} else {
					m.parseLiteralString(stringValue)
				}
			} else {
				return fmt.Errorf("could not identify the type of entry for %s: %v in file %s", key, value, file.FullPath)
			}
		}
	}
	return nil
}

func (MessagesParser) parseLiteralString(message string) {

}

func (MessagesParser) parseParametrizedString(message string) {

}

func (MessagesParser) parseNestedEntries(entries *orderedmap.OrderedMap) {

}

func (MessagesParser) parseConditionalMessage(key string, value *orderedmap.OrderedMap) {

}

func (MessagesParser) copy(arr []string, newElement ...string) []string {
	copied := make([]string, len(arr))
	copy(copied, arr)
	return append(copied, newElement...)
}

package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type CollectionError struct {
}

func (err CollectionError) Error() string {
	return ""
}

type MessageInstance struct {
	Message    string
	TimesFound int
}

type CollectedMessages struct {
	LanguageTag string
	Messages    map[string]*MessageInstance
}

func (cm *CollectedMessages) FindDuplicatedKeys() []string {
	var dups []string
	for key, v := range cm.Messages {
		if v.TimesFound > 1 {
			dups = append(dups, key)
		}
	}
	sort.Strings(dups)
	return dups
}

type MessageCollector interface {
	// Scans the given directory and sub directories for json files containing messages.
	// Returns a map `language tag -> message key -> message`
	FindAllMessagesInDir(dir string) (map[string]*CollectedMessages, error)
}

func copyAndAdd(arr []string, newElement string) []string {
	copied := make([]string, len(arr))
	copy(copied, arr)
	return append(copied, newElement)
}

func collectMessagesInto(messages map[string]any, dest *CollectedMessages, pres []string) {
	for key, value := range messages {
		switch val := value.(type) {
		case map[string]any:
			collectMessagesInto(val, dest, copyAndAdd(pres, key))
		case string:
			fullKey := strings.Join(append(pres, key), ".")
			if instance, found := dest.Messages[fullKey]; found {
				instance.TimesFound++
			} else {
				dest.Messages[fullKey] = &MessageInstance{
					Message:    val,
					TimesFound: 1,
				}
			}
		default:
			panic(fmt.Errorf("entry %s in file %s is not a string or an object", key, strings.Join(pres, string(filepath.Separator))))
		}
	}
}

type JsonMessageScanner struct {
}

func (s JsonMessageScanner) FindAllMessagesInDir(dir string) (map[string]*CollectedMessages, error) {
	result := make(map[string]*CollectedMessages)
	err := s.scanMessagesInDir(dir, result, []string{})
	return result, err
}

func (s JsonMessageScanner) scanMessagesInDir(dir string, dest map[string]*CollectedMessages, pres []string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("could not list files from directory '%s': %w", dir, err)
	}

	for _, file := range files {
		if file.IsDir() {
			s.scanMessagesInDir(filepath.Join(dir, file.Name()), dest, copyAndAdd(pres, file.Name()))
		} else {
			if filepath.Ext(file.Name()) != ".json" {
				continue
			}
			lang := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			collection, ok := dest[lang]
			if !ok {
				collection = &CollectedMessages{
					LanguageTag: lang,
					Messages:    make(map[string]*MessageInstance),
				}
				dest[lang] = collection
			}
			s.scanFileForMessages(filepath.Join(dir, file.Name()), collection, pres)
		}
	}
	return nil
}

func (s JsonMessageScanner) scanFileForMessages(file string, dest *CollectedMessages, pres []string) error {
	contents, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("could not read contents of file %s: %w", strings.Join(append(pres, file), string(filepath.Separator)), err)
	}
	messages := make(map[string]any)
	if err := json.Unmarshal(contents, &messages); err != nil {
		return fmt.Errorf("could not unmarshal contents of file %s: %w", strings.Join(append(pres, file), string(filepath.Separator)), err)
	}
	collectMessagesInto(messages, dest, pres)
	return nil
}

package messagecollector

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

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

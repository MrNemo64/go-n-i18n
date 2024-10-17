package messagecollector

import (
	"fmt"
	"path/filepath"
	"strings"
)

type MessageCollector interface {
	// Scans the given directory and sub directories for json files containing messages.
	// Returns a map `language tag -> message key -> message`
	FindAllMessagesInDir(dir string) (map[string]map[string]string, error)
}

func copyAndAdd(arr []string, newElement string) []string {
	copied := make([]string, len(arr))
	copy(copied, arr)
	return append(copied, newElement)
}

func collectMessagesInto(messages map[string]any, dest map[string]string, pres []string) {
	for key, value := range messages {
		switch val := value.(type) {
		case map[string]any:
			collectMessagesInto(val, dest, copyAndAdd(pres, key))
		case string:
			dest[strings.Join(append(pres, key), ".")] = val
		default:
			panic(fmt.Errorf("entry %s in file %s is not a string or an object", key, strings.Join(pres, string(filepath.Separator))))
		}
	}
}

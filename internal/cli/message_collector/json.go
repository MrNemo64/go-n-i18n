package messagecollector

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type JsonMessageScanner struct {
}

func (s JsonMessageScanner) FindAllMessagesInDir(dir string) (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)
	err := s.scanMessagesInDir(dir, result, []string{})
	return result, err
}

func (s JsonMessageScanner) scanMessagesInDir(dir string, dest map[string]map[string]string, pres []string) error {
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
				collection = make(map[string]string)
				dest[lang] = collection
			}
			s.scanFileForMessages(filepath.Join(dir, file.Name()), collection, pres)
		}
	}
	return nil
}

func (s JsonMessageScanner) scanFileForMessages(file string, dest map[string]string, pres []string) error {
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

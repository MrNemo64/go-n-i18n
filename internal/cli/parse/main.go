package parse

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"regexp"
	"strings"

	"github.com/MrNemo64/go-n-i18n/internal/cli/assert"
	"github.com/MrNemo64/go-n-i18n/internal/cli/types"
	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
	"github.com/iancoleman/orderedmap"
)

var (
	ErrNextFile         util.Error = util.MakeError("could get next file to parse: %w")
	ErrIO                          = util.MakeError("could not read contents of file %s: %w")
	ErrUnmarshal                   = util.MakeError("could not unmarshal contents of file %s: %w")
	ErrInvalidKeyName              = util.MakeError("invalid key in path %s: %w")
	ErrUnknownEntryType            = util.MakeError("could not identify the type of entry in the path %s: %+v")
	ErrAddChildren                 = util.MakeError("could not add child %s to %s: %w")

	ErrKeyIsConditionalButValueIsNotObject = util.MakeError("invalid key '%s': has the ? prefix so it's a conditional key but the value is not an object: %v")
	ErrCouldNotAddEntry                    = util.MakeError("could not add %s entry %s: %w")
	ErrCouldNotAddArg                      = util.MakeError("could not add argument {%s:%s:%s}: %w")
)

var ArgumentExtractor = regexp.MustCompile(`\{([a-zA-Z][a-zA-Z0-9_]*)(?::([a-zA-Z0-9_]+))?(?::([a-zA-Z0-9_%.]+))?\}`)

type JsonParser struct {
	*util.WarningsCollector
}

func ParseJson(walker DirWalker, wc *util.WarningsCollector, log *slog.Logger) (*types.MessageBag, error) {
	return (&JsonParser{WarningsCollector: wc}).ParseWalker(walker)
}

func (p *JsonParser) ParseWalker(walker DirWalker) (*types.MessageBag, error) {
	root := types.MakeRoot()
	for {
		file, err := walker.Next()
		if err == ErrNoMoreFiles {
			return root, nil
		}
		if err != nil {
			return nil, ErrNextFile.WithArgs(err)
		}
		content, err := file.ReadContents()
		if err != nil {
			return nil, ErrIO.WithArgs(file.FullPath(), err)
		}
		entries := orderedmap.New()
		if err := json.Unmarshal(content, entries); err != nil {
			return nil, ErrUnmarshal.WithArgs(file.FullPath(), err)
		}

		dest, err := root.FindOrCreateChildBag(file.Path()...)
		if err != nil {
			return nil, err
		}

		if err := p.ParseGroupOfMessagesInto(dest, entries, file.Language()); err != nil {
			return nil, err
		}
	}
}

func (p *JsonParser) ParseGroupOfMessagesInto(dest *types.MessageBag, entries *orderedmap.OrderedMap, lang string) error {
	keys := entries.Keys()
	for _, key := range keys {
		value, found := entries.Get(key)
		if !found {
			panic(fmt.Errorf("the ordered map is missing the key '%s', this is a bug in the github.com/iancoleman/orderedmap library. Dest: %s", key, dest.PathAsStr()))
		}

		if strings.HasPrefix(key, "?") { // is conditional?
			key = key[1:]
			if err := types.CheckKey(key); err != nil {
				p.AddWarning(ErrInvalidKeyName.WithArgs(types.PathAsStr(types.ResolveFullPath(dest, key)), err))
				continue
			}
			panic("todo: parsing conditionals")
		}

		if err := types.CheckKey(key); err != nil {
			p.AddWarning(ErrInvalidKeyName.WithArgs(types.PathAsStr(types.ResolveFullPath(dest, key)), err))
			continue
		}

		if inner, ok := value.(orderedmap.OrderedMap); ok { // is bag?
			newDest, err := dest.FindOrCreateChildBag(key)
			if err != nil {
				p.AddWarning(ErrAddChildren.WithArgs(key, dest.PathAsStr(), err))
				continue
			}
			if err := p.ParseGroupOfMessagesInto(newDest, &inner, lang); err != nil {
				return err
			}
			continue
		}

		parsed, ok := p.ParseMessageValue(types.PathAsStr(types.ResolveFullPath(dest, key)), value)
		if !ok {
			continue
		}
		newEntry, err := types.NewMessageInstance(key)
		assert.NoError(err) // key is valid, we checked it above
		err = newEntry.AddLanguage(lang, parsed)
		assert.NoError(err) // entry is empty, it must accept the new language
		if err := dest.AddChildren(newEntry); err != nil {
			p.AddWarning(ErrAddChildren.WithArgs(key, dest.PathAsStr(), err))
		}
	}
	return nil
}

func (p *JsonParser) ParseMessageValue(fullKey string, value any) (types.MessageValue, bool) {
	switch value.(type) {
	case string:
		str := value.(string)
		if p.HasArguments(str) {
			panic("tofo")
		} else {
			return types.NewStringLiteralValue(str), true
		}
	case []any:
		arr := value.([]any)
		if !p.IsStringSlice(arr) {
			p.AddWarning(ErrUnknownEntryType.WithArgs(fullKey, value))
			return nil, false
		}
		panic("tofo")
	default:
		p.AddWarning(ErrUnknownEntryType.WithArgs(fullKey, value))
		return nil, false
	}
}

func (*JsonParser) HasArguments(str string) bool { return ArgumentExtractor.MatchString(str) }
func (*JsonParser) IsStringSlice(arr []any) bool {
	for i := range arr {
		if _, ok := arr[i].(string); !ok {
			return false
		}
	}
	return true
}

package parse

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/MrNemo64/go-n-i18n/internal/cli/assert"
	"github.com/MrNemo64/go-n-i18n/internal/cli/types"
	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
	"github.com/iancoleman/orderedmap"
)

var (
	ErrNextFile                    util.Error = util.MakeError("could get next file to parse: %w")
	ErrIO                                     = util.MakeError("could not read contents of file %s: %w")
	ErrUnmarshal                              = util.MakeError("could not unmarshal contents of file %s: %w")
	ErrInvalidKeyName                         = util.MakeError("invalid key in path %s: %w")
	ErrInvalidBagName                         = util.MakeError("invalid bag name in path %s: %w")
	ErrBagNameReasignation                    = util.MakeError("the bag %s in the lang %s has the name %s but it got reasigned to %s")
	ErrUnknownEntryType                       = util.MakeError("could not identify the type of entry in the path %s: %+v")
	ErrAddChildren                            = util.MakeError("could not add child %s to %s: %w")
	ErrUnknwonArgumentType                    = util.MakeError("unknown argument type '%s' in path %s, using the unknown type")
	ErrInvalidConditionalEntry                = util.MakeError("the entry %s in the lang %s is marked as conditional but no conditions are provided")
	ErrInvalidConditionalCondition            = util.MakeError("the condition %s in the path %s in the lang %s is not a valid conditional value")
	ErrInvalidConditional                     = util.MakeError("the conditional in the path %s in the lang %s is not a valid: %w")

	ErrKeyIsConditionalButValueIsNotObject = util.MakeError("invalid key '%s': has the ? prefix so it's a conditional key but the value is not an object: %v")
	ErrCouldNotAddEntry                    = util.MakeError("could not add %s entry %s: %w")
	ErrCouldNotAddArg                      = util.MakeError("could not add argument {%s:%s:%s}: %w")
)

var ArgumentExtractor = regexp.MustCompile(`\{([a-zA-Z_]\w*):?(\w*)?:?([\w\.]*)?\}`)

type JsonParser struct {
	*util.WarningsCollector
	argProvider *types.ArgumentProvider
}

func ParseJson(walker DirWalker, wc *util.WarningsCollector, argProvider *types.ArgumentProvider) (*types.MessageBag, error) {
	return (&JsonParser{WarningsCollector: wc, argProvider: argProvider}).ParseWalker(walker)
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
			mapValue, ok := value.(orderedmap.OrderedMap)
			if !ok {
				p.WarningsCollector.AddWarning(ErrInvalidConditionalEntry.WithArgs(types.PathAsStr(types.ResolveFullPath(dest, key)), lang))
				continue
			}
			args := types.NewArgumentList()
			parsed, ok := p.ParseConditionalMessageValue(types.PathAsStr(types.ResolveFullPath(dest, key)), &mapValue, args, lang)
			if !ok {
				continue
			}
			newEntry, err := types.NewMessageInstance(key)
			assert.NoError(err)                                // key is valid, we checked it above
			assert.NoError(newEntry.AddArgs(args))             // entry is empty, it must accept the new args
			assert.NoError(newEntry.AddLanguage(lang, parsed)) // entry is empty, it must accept the new language
			if err := dest.AddChildren(newEntry); err != nil {
				p.AddWarning(ErrAddChildren.WithArgs(key, dest.PathAsStr(), err))
			}
			continue
		}

		if inner, ok := value.(orderedmap.OrderedMap); ok { // is bag or parametrized with `_args` to specify args
			if _, found := inner.Get("_args"); found { // parametrized with `_args``
				panic("todo")
			} else { // bag
				name := ""
				if strings.Contains(key, ":") {
					parts := strings.SplitN(key, ":", 2)
					if len(parts) == 1 || parts[1] == "" {
						name = parts[0]
					} else {
						name = parts[1]
					}
					key = key[:strings.Index(key, ":")]
				}

				if err := types.CheckKey(key); err != nil {
					p.AddWarning(ErrInvalidKeyName.WithArgs(types.PathAsStr(types.ResolveFullPath(dest, key)), err))
					continue
				}

				newDest, err := dest.FindOrCreateChildBag(key)
				if err != nil {
					p.AddWarning(ErrAddChildren.WithArgs(key, dest.PathAsStr(), err))
					continue
				}
				if newDest.Name == "" && name != "" {
					if err := types.CheckName(name); err != nil {
						p.AddWarning(ErrInvalidBagName.WithArgs(types.PathAsStr(types.ResolveFullPath(dest, key)), err))
						continue
					}
					newDest.Name = name
				} else if newDest.Name != name && name != "" {
					p.WarningsCollector.AddWarning(ErrBagNameReasignation.WithArgs(types.PathAsStr(types.ResolveFullPath(dest, key)), lang, newDest.Name, name))
					continue
				}
				if err := p.ParseGroupOfMessagesInto(newDest, &inner, lang); err != nil {
					return err
				}
				continue
			}
		}

		if err := types.CheckKey(key); err != nil {
			p.AddWarning(ErrInvalidKeyName.WithArgs(types.PathAsStr(types.ResolveFullPath(dest, key)), err))
			continue
		}

		args := types.NewArgumentList()
		parsed, ok := p.ParseMessageValue(types.PathAsStr(types.ResolveFullPath(dest, key)), value, args)
		if !ok {
			continue
		}
		newEntry, err := types.NewMessageInstance(key)
		assert.NoError(err)                                // key is valid, we checked it above
		assert.NoError(newEntry.AddArgs(args))             // entry is empty, it must accept the new args
		assert.NoError(newEntry.AddLanguage(lang, parsed)) // entry is empty, it must accept the new language
		if err := dest.AddChildren(newEntry); err != nil {
			p.AddWarning(ErrAddChildren.WithArgs(key, dest.PathAsStr(), err))
		}
	}
	return nil
}

func (p *JsonParser) ParseMessageValue(fullKey string, value any, argList *types.ArgumentList) (types.MessageValue, bool) {
	switch value.(type) {
	case string:
		str := value.(string)
		if !p.HasArguments(str) {
			return types.NewStringLiteralValue(str), true
		}
		return p.ParseParametrizedMessage(fullKey, str, argList)
	case []any:
		arr := value.([]any)
		if len(arr) == 0 || !p.IsStringSlice(arr) {
			p.AddWarning(ErrUnknownEntryType.WithArgs(fullKey, value))
			return nil, false
		}
		lines := make([]types.Multilineable, 0)
		for _, line := range arr {
			str := line.(string)
			if !p.HasArguments(str) {
				lines = append(lines, types.NewStringLiteralValue(str))
			} else {
				if parsed, ok := p.ParseParametrizedMessage(fullKey, str, argList); ok {
					lines = append(lines, parsed)
				} else {
					return nil, false
				}
			}
		}
		multi, err := types.NewMultilineValue(lines)
		assert.NoError(err) // err if len(lines) == 0 but we checked above
		return multi, true
	default:
		p.AddWarning(ErrUnknownEntryType.WithArgs(fullKey, value))
		return nil, false
	}
}

func (p *JsonParser) ParseConditionalMessageValue(fullKey string, value *orderedmap.OrderedMap, argList *types.ArgumentList, lang string) (*types.ValueConditional, bool) {
	finishOk := true
	var conditions []types.Condition
	var elseCondition types.Conditionable
	for _, condition := range value.Keys() {
		if condition == "_args" {
			panic(fmt.Errorf("specifying the args in a `_args` entry is not yet supported"))
		}
		value, found := value.Get(condition)
		if !found {
			panic(fmt.Errorf("the ordered map is missing the key '%s', this is a bug in the github.com/iancoleman/orderedmap library. Dest: %s", condition, fullKey))
		}
		parsed, ok := p.ParseMessageValue(fullKey+"."+condition, value, argList)
		if !ok {
			finishOk = false
			continue
		}
		conditionValue, ok := parsed.(types.Conditionable)
		if !ok {
			finishOk = false
			p.AddWarning(ErrInvalidConditionalCondition.WithArgs(condition, fullKey, lang))
			continue
		}
		if condition == "" {
			elseCondition = conditionValue
		} else {
			conditions = append(conditions, types.Condition{
				Condition: condition,
				Value:     conditionValue,
			})
		}
	}
	if !finishOk {
		return nil, false
	}
	cond, err := types.NewConditionalValue(conditions, elseCondition)
	if err != nil {
		p.WarningsCollector.AddWarning(ErrInvalidConditional.WithArgs(fullKey, lang, err))
		return nil, false
	}
	return cond, true
}

func (p *JsonParser) ParseParametrizedMessage(fullKey string, str string, argList *types.ArgumentList) (*types.ValueParametrized, bool) {
	textSegments, arguments := p.SeparateArgumentsFromText(str)
	if len(textSegments) != len(arguments)+1 {
		panic(fmt.Errorf("JsonParser.SeparateArgumentsFromText returned an unexpected amount of text segments (%d) and arguments (%d) for the path %s", len(textSegments), len(arguments), fullKey))
	}
	usedArgs := util.Map(arguments, func(index int, foundArg *foundArgument) *types.UsedArgument {
		argType, found := p.argProvider.FindArgument(foundArg.Type)
		if !found {
			if foundArg.Type != "" {
				p.WarningsCollector.AddWarning(ErrUnknwonArgumentType.WithArgs(foundArg.Type, fullKey))
			}
			argType = p.argProvider.UnknwonType()
		}
		arg, err := argList.AddArgument(&types.MessageArgument{
			Name: foundArg.Name,
			Type: argType,
		})
		if err != nil {
			p.WarningsCollector.AddWarning(err)
			return nil
		}
		return &types.UsedArgument{
			Argument: arg,
			Format:   foundArg.Format,
		}
	})
	if util.Has(usedArgs, nil) {
		return nil, false
	}
	parametrized, err := types.NewParametrizedStringValue(
		util.Map(textSegments, func(_ int, t *string) *types.ValueString { return types.NewStringLiteralValue(*t) }),
		usedArgs,
	)
	assert.NoError(err)
	return parametrized, true
}

type foundArgument struct {
	Name   string
	Type   string
	Format string
}

func (p *JsonParser) SeparateArgumentsFromText(message string) ([]string, []foundArgument) {
	var textSegments []string
	var arguments []foundArgument

	// Track the position as we move through the string
	lastIndex := 0
	matches := ArgumentExtractor.FindAllStringSubmatchIndex(message, -1)

	// If the first match starts at index 0, add an empty text segment at the beginning
	if len(matches) > 0 && matches[0][0] == 0 {
		textSegments = append(textSegments, "")
	}

	for i, match := range matches {
		start, end := match[0], match[1]

		// Capture the normal text before this argument
		if start > lastIndex {
			textSegments = append(textSegments, message[lastIndex:start])
		} else if i > 0 {
			// If two arguments are consecutive, insert an empty text segment
			textSegments = append(textSegments, "")
		}

		// Extract components based on regex capture groups
		name := message[match[2]:match[3]]
		argType := ""
		format := ""
		if match[4] != -1 {
			argType = message[match[4]:match[5]]
		}
		if match[6] != -1 {
			format = message[match[6]:match[7]]
		}

		// Create an Argument and add to the list
		arguments = append(arguments, foundArgument{Name: name, Type: argType, Format: format})

		// Update lastIndex to continue after this match
		lastIndex = end
	}

	// Append any remaining text after the last argument
	if lastIndex < len(message) {
		textSegments = append(textSegments, message[lastIndex:])
	} else if len(matches) > 0 && lastIndex == len(message) {
		// If the last match ends at the end of the input, add an empty text segment at the end
		textSegments = append(textSegments, "")
	}

	return textSegments, arguments
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

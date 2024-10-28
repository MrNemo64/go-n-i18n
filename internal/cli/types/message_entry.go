package types

import (
	"regexp"

	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
)

type MessageEntryType int

const (
	MessageEntryInstance MessageEntryType = iota
	MessageEntryBag
)

var (
	ErrInvalidKey  util.Error = util.MakeError("the key '%s' does not follow the allowed format (^[a-zA-Z][a-zA-Z0-9_-]*$)")
	ErrCreateEntry            = util.MakeError("could not make entry with key '%s': %w")
)

var ValidKey = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_-]*$")

func IsValidKey(key string) bool { return ValidKey.Match([]byte(key)) }
func CheckKey(key string) error {
	if !IsValidKey(key) {
		return ErrInvalidKey.WithArgs(key)
	}
	return nil
}

type MessageEntry interface {
	Key() string
	Parent() *MessageBag
	AssignParent(*MessageBag)
	Path() []string
	PathAsStr() string
	Type() MessageEntryType
	Languages() *util.Set[string]
	MustHaveAllLangs(langs []string, defLang string) map[string][]string

	IsBag() bool
	IsInstance() bool
	AsBag() *MessageBag
	AsInstance() *MessageInstance
}

type messageEntry struct {
	key    string
	parent *MessageBag
}

func (e *messageEntry) Key() string {
	return e.key
}

func (e *messageEntry) Parent() *MessageBag {
	return e.parent
}

func (e *messageEntry) AssignParent(parent *MessageBag) {
	e.parent = parent
}

func (e *messageEntry) Path() []string {
	return ResolveFullPath(e.Parent(), e.Key())
}

func (e *messageEntry) PathAsStr() string {
	return PathAsStr(e.Path())
}

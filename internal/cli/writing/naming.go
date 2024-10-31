package writing

import (
	"strings"

	"github.com/MrNemo64/go-n-i18n/internal/cli/types"
)

type MessageEntryNamer interface {
	FunctionName(me types.MessageEntry) string
	InterfaceName(me *types.MessageBag) string
	FunctionNameForLang(lang string, me types.MessageEntry) string
	InterfaceNameForLang(lang string, me *types.MessageBag) string
	TopLevelName() string
}

type goNamer struct{}

func GoNamer() *goNamer {
	return &goNamer{}
}

func (*goNamer) toGo(key string, private bool) string {
	newName := key[:1]
	if !private {
		newName = strings.ToUpper(key[:1])
	}

	for j := 1; j < len(key); j++ {
		if key[j] == '-' || key[j] == '_' {
			if j == len(key)-1 {
				break // last char is a _ or a - so just ignore it
			}
			j++
			newName += strings.ToUpper(key[j : j+1])
		} else {
			newName += key[j : j+1]
		}
	}
	return newName
}

func (m *goNamer) FunctionName(me types.MessageEntry) string {
	if me.Key() == "" {
		panic("tryed to get the function name of the root bag")
	}
	return m.toGo(me.Key(), false)
}

func (m goNamer) FunctionNameForLang(lang string, me types.MessageEntry) string {
	return strings.ReplaceAll(lang, "-", "_") + "_" + m.FunctionName(me)
}

func (m *goNamer) InterfaceName(me *types.MessageBag) string {
	if me.Key() == "" {
		return m.TopLevelName()
	}
	if me.Name != "" {
		return m.toGo(me.Name, false)
	}
	name := ""
	for _, part := range me.Path() {
		name += m.toGo(part, false)
	}
	return name
}

func (m *goNamer) InterfaceNameForLang(lang string, me *types.MessageBag) string {
	return strings.ReplaceAll(lang, "-", "_") + "_" + m.InterfaceName(me)
}

func (m *goNamer) TopLevelName() string {
	return "Messages"
}

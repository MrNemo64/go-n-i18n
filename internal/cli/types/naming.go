package types

import (
	"strings"
)

type MessageEntryNamer struct {
}

func (MessageEntryNamer) toGo(key string, private bool) string {
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

func (m MessageEntryNamer) FunctionName(me MessageEntry) string {
	if me.Key() == "" {
		panic("tryed to get the function name of the root bag")
	}
	return m.toGo(me.Key(), false)
}

func (m MessageEntryNamer) InterfaceName(me MessageEntry) string {
	if me.Key() == "" {
		return m.TopLevelName()
	}
	name := ""
	for _, part := range me.Path() {
		name += m.toGo(part, false)
	}
	return name
}

func (m MessageEntryNamer) TopLevelName() string {
	return "Messages"
}

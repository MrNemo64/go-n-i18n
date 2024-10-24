package writing

import (
	"fmt"
	"strings"

	"github.com/MrNemo64/go-n-i18n/internal/cli/assert"
	"github.com/MrNemo64/go-n-i18n/internal/cli/types"
	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
)

type GoCodeWriter struct {
	sb      *strings.Builder
	def     *types.InterfaceDefinition
	langs   []string
	defLang string
	pack    string
}

func GenerateGoCode(definition *types.InterfaceDefinition, langs []string, defLang, pack string) string {
	assert.Has(langs, defLang)
	cw := GoCodeWriter{
		sb:      &strings.Builder{},
		def:     definition,
		langs:   util.Map(langs, func(t *string) string { return strings.ReplaceAll(*t, "-", "_") }),
		defLang: strings.ReplaceAll(defLang, "-", "_"),
		pack:    pack,
	}
	cw.GenerateCode()
	return cw.sb.String()
}

func (w *GoCodeWriter) GenerateCode() {
	w.WriteHeader()
	w.WriteGetMethods()
}

func (w *GoCodeWriter) WriteHeader() {
	w.w("/** Code generated using https://github.com/MrNemo64/go-n-i18n \n")
	w.w(" * Any changes to this file will be lost on the next tool run */\n\n")
	w.w("package ")
	w.w(w.pack)
	w.w("\n\n")
	w.w("import (")
	w.w("    \"fmt\"")
	w.w("    \"strings\"")
	w.w(")")
	w.w("\n\n")
}

func (w *GoCodeWriter) WriteGetMethods() {
	w.w("func MessagesFor(tag string) (%s, bool) {\n", w.def.Name)
	w.w("    switch strings.ReplaceAll(tag, \"-\", \"_\") {")
	for _, lang := range w.langs {
		w.w("    case \"%s\":\n", lang)
		w.w("        return %s_%s{}, true\n", lang, w.def.Name)
	}
	w.w("    }\n")
	w.w("    return nil, false")
	w.w("}\n\n")

	w.w("func MessagesForMust(tag string) %s {\n", w.def.Name)
	w.w("    switch strings.ReplaceAll(tag, \"-\", \"_\") {")
	for _, lang := range w.langs {
		w.w("    case \"%s\":\n", lang)
		w.w("        return %s_%s{}\n", lang, w.def.Name)
	}
	w.w("    }\n")
	w.w("    panic(fmt.Errorf(\"unknwon language tag: \" + tag))")
	w.w("}\n\n")

	w.w("func MessagesForOrDefault(tag string) %s {\n", w.def.Name)
	w.w("    switch strings.ReplaceAll(tag, \"-\", \"_\") {")
	for _, lang := range w.langs {
		w.w("    case \"%s\":\n", lang)
		w.w("        return %s_%s{}\n", lang, w.def.Name)
	}
	w.w("    }\n")
	w.w("    return %s_%s{}\n", w.defLang, w.def.Name)
	w.w("}\n\n")
}

func (w *GoCodeWriter) WriteInterfaces() {
	w.writeInterface(w.def)
}

func (w *GoCodeWriter) writeInterface(i *types.InterfaceDefinition) {
	w.w("type %s interface{\n", i.Name)
	for _, f := range i.Functions {
		w.w("    %s() %s\n", f.Name(), f.ReturnType())
	}
	w.w("}\n")

	for _, ii := range i.Interfaces {
		w.writeInterface(ii)
	}
}

func (w *GoCodeWriter) w(str string, args ...any) {
	w.sb.WriteString(fmt.Sprintf(str, args...))
}

package writing

import (
	"fmt"
	"strings"

	"github.com/MrNemo64/go-n-i18n/internal/cli/assert"
	"github.com/MrNemo64/go-n-i18n/internal/cli/types"
	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
)

type GoCodeWriter struct {
	sb        *strings.Builder
	indent    int
	inNewLine bool
	msgs      *types.MessageBag
	namer     MessageEntryNamer
	langs     []string
	defLang   string
	pack      string
}

func GenerateGoCode(msgs *types.MessageBag, namer MessageEntryNamer, langs []string, defLang, pack string) string {
	assert.Has(langs, defLang)
	cw := GoCodeWriter{
		sb:        &strings.Builder{},
		msgs:      msgs,
		indent:    0,
		inNewLine: true,
		namer:     namer,
		langs:     langs,
		defLang:   defLang,
		pack:      pack,
	}
	cw.GenerateCode()
	return cw.sb.String()
}

func (w *GoCodeWriter) GenerateCode() {
	w.WriteHeader()
	w.WriteGetMethods()
	w.WriteInterfaces()
	w.WriteStructs()
}

func (w *GoCodeWriter) WriteHeader() {
	w.w("/** Code generated using https://github.com/MrNemo64/go-n-i18n \n")
	w.w(" * Any changes to this file will be lost on the next tool run */\n\n")
	w.w("package ")
	w.w(w.pack)
	w.w("\n\n")
	w.w("import (\n")
	w.w("    \"fmt\"\n")
	w.w("    \"strings\"\n")
	w.w(")\n\n")
}

func (w *GoCodeWriter) WriteGetMethods() {
	w.w("func MessagesFor(tag string) (%s, bool) {\n", w.namer.TopLevelName())
	w.w("    switch strings.ReplaceAll(tag, \"-\", \"_\") {\n")
	for _, lang := range w.langs {
		w.w("    case \"%s\":\n", lang)
		w.w("        return %s{}, true\n", w.namer.InterfaceNameForLang(lang, w.msgs))
	}
	w.w("    }\n")
	w.w("    return nil, false")
	w.w("}\n\n")

	w.w("func MessagesForMust(tag string) %s {\n", w.namer.TopLevelName())
	w.w("    switch strings.ReplaceAll(tag, \"-\", \"_\") {\n")
	for _, lang := range w.langs {
		w.w("    case \"%s\":\n", lang)
		w.w("        return %s{}\n", w.namer.InterfaceNameForLang(lang, w.msgs))
	}
	w.w("    }\n")
	w.w("    panic(fmt.Errorf(\"unknwon language tag: \" + tag))")
	w.w("}\n\n")

	w.w("func MessagesForOrDefault(tag string) %s {\n", w.namer.TopLevelName())
	w.w("    switch strings.ReplaceAll(tag, \"-\", \"_\") {\n")
	for _, lang := range w.langs {
		w.w("    case \"%s\":\n", lang)
		w.w("        return %s{}\n", w.namer.InterfaceNameForLang(lang, w.msgs))
	}
	w.w("    }\n")
	w.w("    return %s{}\n", w.namer.InterfaceNameForLang(w.defLang, w.msgs))
	w.w("}\n\n")
}

func (w *GoCodeWriter) WriteInterfaces() {
	w.writeInterface(w.msgs)
	w.w("\n")
}

func (w *GoCodeWriter) writeInterface(i *types.MessageBag) {
	w.w("type %s interface{\n", w.namer.InterfaceName(i))
	w.addIndent()
	for _, child := range i.Children() {
		w.w("%s(%s) ", w.namer.FunctionName(child), w.createArgList(child))
		switch child.Type() {
		case types.MessageEntryBag:
			w.w("%s\n", w.namer.InterfaceName(child.AsBag()))
		case types.MessageEntryInstance:
			w.w("string\n")
		default:
			panic(fmt.Errorf("unknown message entry type %d", child.Type()))
		}
	}
	w.removeIndent()
	w.w("}\n")

	for _, child := range i.Children() {
		if child.IsBag() {
			w.writeInterface(child.AsBag())
		}
	}
}

func (w *GoCodeWriter) WriteStructs() {
	for _, lang := range w.langs {
		w.writeStruct(lang, w.msgs)
		w.w("\n\n")
	}
}

func (w *GoCodeWriter) writeStruct(lang string, msgs *types.MessageBag) {
	w.w("type %s struct{}\n", w.namer.InterfaceNameForLang(lang, msgs))
	for _, child := range msgs.Children() {
		w.writeFunction(lang, child)
		if child.IsBag() {
			w.writeStruct(lang, child.AsBag())
		}
	}
}

func (w *GoCodeWriter) writeFunction(lang string, msg types.MessageEntry) {
	w.w("func (%s) %s(%s) ", w.namer.InterfaceNameForLang(lang, msg.Parent()), w.namer.FunctionName(msg), w.createArgList(msg))
	switch msg.Type() {
	case types.MessageEntryBag:
		w.w("%s {\n", w.namer.InterfaceName(msg.AsBag()))
		w.w("    return %s{}\n", w.namer.InterfaceNameForLang(lang, msg.AsBag()))
		w.w("}\n")
	case types.MessageEntryInstance:
		w.w("string {\n")
		w.addIndent()
		w.writeFunctionBody(lang, msg.AsInstance())
		w.removeIndent()
		w.w("}\n")
	default:
		panic(fmt.Errorf("unknown message entry type %d", msg.Type()))
	}
}

func (w *GoCodeWriter) createArgList(msg types.MessageEntry) string {
	switch msg.Type() {
	case types.MessageEntryBag:
		return ""
	case types.MessageEntryInstance:
		return strings.Join(
			util.Map(msg.AsInstance().Args().Args, func(_ int, t **types.MessageArgument) string { return (*t).Name + " " + (*t).Type.Type }),
			", ",
		)
	default:
		panic(fmt.Errorf("unknown message entry type %d", msg.Type()))
	}
}

func (w *GoCodeWriter) writeFunctionBody(lang string, msg *types.MessageInstance) {
	val := msg.MessageMust(lang)
	switch val.(type) {
	case *types.ValueString:
		w.w("return \"%s\"\n", val.AsValueString().Escaped("\""))
	case *types.ValueParametrized:
		p := val.AsValueParametrized()
		messagePart := strings.Builder{}
		for i, arg := range p.Args {
			messagePart.WriteString(p.TextSegments[i].Escaped("\""))
			messagePart.WriteString("%")
			if arg.Format == "" {
				messagePart.WriteString(p.Args[i].Argument.Type.DefaultFormat)
			} else {
				messagePart.WriteString(p.Args[i].Format)
			}
		}
		messagePart.WriteString(p.TextSegments[len(p.TextSegments)-1].Escaped("\""))
		argListPart := strings.Join(
			util.Map(val.AsValueParametrized().Args, func(_ int, t **types.UsedArgument) string { return (*t).Argument.Name }),
			", ",
		)
		w.w("return fmt.Sprintf(\"%s\", %s)\n", messagePart.String(), argListPart)
	}
}

func (w *GoCodeWriter) w(str string, args ...any) {
	if w.indent > 0 && w.inNewLine {
		w.sb.WriteString(strings.Repeat(" ", w.indent))
	}
	msg := fmt.Sprintf(str, args...)
	w.sb.WriteString(msg)
	w.inNewLine = strings.HasSuffix(msg, "\n")
}

func (w *GoCodeWriter) indentBy(amount int) {
	w.indent = max(0, w.indent+amount)
}

func (w *GoCodeWriter) addIndent() {
	w.indentBy(4)
}

func (w *GoCodeWriter) removeIndent() {
	w.indentBy(-4)
}

package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type codeWritter struct {
	*strings.Builder
	args     CliArgs
	messages *MessageEntryMessageBag
	namer    MessageEntryNamer
	def      *InterfaceDefinition
	langs    []string
}

func WriteCode(messages *MessageEntryMessageBag, args CliArgs) error {
	namer := MessageEntryNamer{}
	cw := &codeWritter{
		Builder:  &strings.Builder{},
		args:     args,
		messages: messages,
		namer:    namer,
		def:      messages.DefineInterface(namer),
		langs:    messages.Languages().Get(),
	}
	return cw.write()
}

func (w *codeWritter) write() error {
	w.WriteString("package ")
	w.WriteString(w.args.Package)
	w.WriteString("\n\n")
	w.WriteString(`import "fmt"`)
	w.WriteString("\n\n")

	w.writeGetMethods()
	w.writeInterface(w.def)
	w.writeStructs()

	file, err := os.Create(w.args.OutFile)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(w.String())
	return err

}

func (w *codeWritter) writeGetMethods() {
	langs := w.messages.Languages().Get()
	w.WriteString(fmt.Sprintf("func MessagesFor(tag string) (%s, error) {\n", w.namer.TopLevelName()))
	w.WriteString("    switch tag {\n")
	for _, lang := range langs {
		w.WriteString(`    case "` + lang + "\":\n")
		w.WriteString(fmt.Sprintf("        return %s_%s{}, nil\n", strings.ReplaceAll(lang, "-", "_"), w.namer.TopLevelName()))
	}
	w.WriteString("    }\n")
	w.WriteString("    return nil, fmt.Errorf(\"unknown language tag: %s\", tag)\n")
	w.WriteString("}")
	w.WriteString("\n\n")
	w.WriteString(fmt.Sprintf("func MessagesForMust(tag string) %s {\n", w.namer.TopLevelName()))
	w.WriteString("    switch tag {\n")
	for _, lang := range langs {
		w.WriteString(`    case "` + lang + "\":\n")
		w.WriteString(fmt.Sprintf("        return %s_%s{}\n", strings.ReplaceAll(lang, "-", "_"), w.namer.TopLevelName()))
	}
	w.WriteString("    }\n")
	w.WriteString("    panic(fmt.Errorf(\"unknown language tag: %s\", tag))\n")
	w.WriteString("}")
	w.WriteString("\n\n")
	w.WriteString(fmt.Sprintf("func MessagesForOrDefault(tag string) %s {\n", w.namer.TopLevelName()))
	w.WriteString("    switch tag {\n")
	for _, lang := range langs {
		w.WriteString(`    case "` + lang + "\":\n")
		w.WriteString(fmt.Sprintf("        return %s_%s{}\n", strings.ReplaceAll(lang, "-", "_"), w.namer.TopLevelName()))
	}
	w.WriteString("    }\n")
	w.WriteString(fmt.Sprintf("    return %s_%s{}\n", strings.ReplaceAll(w.args.DefaultLanguage, "-", "_"), w.namer.TopLevelName()))
	w.WriteString("}")
	w.WriteString("\n\n")
}

func (w *codeWritter) writeInterface(def *InterfaceDefinition) {
	w.WriteString("type ")
	w.WriteString(def.Name)
	w.WriteString(" interface {\n")
	for _, functionDefinition := range def.Functions {
		w.WriteString("    ")
		w.WriteString(functionDefinition.Name())
		w.writeFunctionArgs(functionDefinition.Arguments())
		w.WriteString(functionDefinition.ReturnType())
		w.WriteString("\n")
	}
	w.WriteString("}\n\n")
	for _, interfaceDefinition := range def.Interfaces {
		w.writeInterface(interfaceDefinition)
	}
}

func (w *codeWritter) writeFunctionArgs(args []*MessageArgument) {
	if len(args) == 0 {
		w.WriteString("() ")
		return
	}
	w.WriteString("(")
	w.WriteString(strings.Join(Map(args, func(t **MessageArgument) string {
		ma := *t
		return ma.Name + " " + ma.Type.Type
	}), ", "))
	w.WriteString(") ")
}

func (w *codeWritter) writeStructs() {
	for _, lang := range w.langs {
		w.writeLangStruct(strings.ReplaceAll(lang, "-", "_"), w.def)
		w.WriteString("\n")
	}
}

func (w *codeWritter) writeLangStruct(lang string, def *InterfaceDefinition) {
	structName := fmt.Sprintf("%s_%s", lang, def.Name)
	w.WriteString(fmt.Sprintf("type %s struct{}\n", structName))
	for _, f := range def.Functions {
		switch f.(type) {
		case *BagFunctionDefinition:
			w.writeBagFunctions(lang, structName, f.(*BagFunctionDefinition))
		case *MessageFunctionDefinition:
			w.writeMessageFunctions(lang, structName, f.(*MessageFunctionDefinition))
		case *ParametrizedFunctionDefinition:
			w.writeParametrizedFunctions(lang, structName, f.(*ParametrizedFunctionDefinition))
		}
	}
	for _, interfaceDefinition := range def.Interfaces {
		w.writeLangStruct(lang, interfaceDefinition)
	}
}

func (w *codeWritter) writeBagFunctions(lang string, structName string, f *BagFunctionDefinition) {
	w.WriteString(fmt.Sprintf("func (%s) %s() %s { return %s_%s{} }\n", structName, f.Name(), f.ReturnType(), lang, f.ReturnType()))
}

func (w *codeWritter) writeMessageFunctions(lang string, structName string, f *MessageFunctionDefinition) {
	w.WriteString(fmt.Sprintf("func (%s) %s() string { return \"%s\" }\n", structName, f.Name(), f.Message.Message(lang)))
}

func (w *codeWritter) writeParametrizedFunctions(lang string, structName string, f *ParametrizedFunctionDefinition) {
	w.WriteString(fmt.Sprintf("func (%s) %s", structName, f.Name()))
	w.writeFunctionArgs(f.Args)
	w.WriteString("string {\n")
	w.WriteString("    return fmt.Sprintf(\"")
	msg := f.Message.Message(lang)
	msg, args := w.extractPartsAndFormats(msg, f.Args)
	w.WriteString(msg)
	w.WriteString("\", ")
	w.WriteString(strings.Join(args, ", "))
	w.WriteString(")\n")
	w.WriteString("}\n")
}

func (*codeWritter) extractPartsAndFormats(msg string, args []*MessageArgument) (string, []string) {
	// TODO already needed to copy paste this pattern in 2 places
	// ALSO, the extraction of parts should probably be done on the parsing side not the writing but oh well...
	argumentExtractor := regexp.MustCompile(`\{([a-zA-Z_][a-zA-Z0-9_]*)(?::([a-zA-Z0-9_]+))?(?::([a-zA-Z0-9_%.]+))?\}`)
	var argNames []string

	argMap := make(map[string]string)
	for _, arg := range args {
		argMap[arg.Name] = "%" + arg.Format
	}

	formattedMessage := argumentExtractor.ReplaceAllStringFunc(msg, func(match string) string {
		matches := argumentExtractor.FindStringSubmatch(match)
		argName := matches[1]
		argNames = append(argNames, argName)
		if format, found := argMap[argName]; found {
			return format
		}
		return match
	})

	return formattedMessage, argNames
}

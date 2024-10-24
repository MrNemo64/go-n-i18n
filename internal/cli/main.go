package cli

import (
	"log/slog"
	"os"
	"strings"

	"github.com/MrNemo64/go-n-i18n/internal/cli/parse"
	"github.com/MrNemo64/go-n-i18n/internal/cli/types"
	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
	"github.com/MrNemo64/go-n-i18n/internal/cli/writing"
)

type CliArgs struct {
	MessagesDirectory string
	DefaultLanguage   string
	OutFile           string
	Package           string
	LogLevel          slog.Level
}

func Run(args CliArgs) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     args.LogLevel,
	}))
	wc := util.NewWarningsCollector()

	log.Info("Collecting files")
	walker, err := parse.IoDirWalker(args.MessagesDirectory)
	if err != nil {
		log.Error("Could not collect all files in the messages directory", "err", err)
		os.Exit(1)
	}

	log.Info("Parsing files")
	messages, err := parse.ParseJson(walker, wc, log)
	if err != nil {
		log.Error("Could not parse all files in the messages directory", "err", err)
		os.Exit(1)
	}

	if !wc.IsEmpty() {
		for _, warning := range wc.Warnings() {
			log.Warn(warning.Error())
		}
		os.Exit(1)
	}

	allLangs := messages.Languages()
	if !allLangs.Contains(args.DefaultLanguage) {
		log.Error("Could not find messages of the default language")
		os.Exit(1)
	}

	removed := messages.RemoveEntriesWithoutLang(args.DefaultLanguage)
	if len(removed) > 0 {
		log.Warn("Removed entries without the default language", "default-language", args.DefaultLanguage,
			"removed-entries", util.Map(removed, func(t *types.MessageEntry) string { return (*t).PathAsStr() }))
	}

	filled := messages.MustHaveAllLangs(allLangs.Get(), args.DefaultLanguage)
	if len(filled) > 0 {
		log.Warn("Some entries were missing in some languages. Using the message of the default language", "missing-entries", filled)
	}

	log.Info("Generating code")
	code := writing.GenerateGoCode(messages.DefineInterface(types.MessageEntryNamer{}), allLangs.Get(), args.DefaultLanguage, args.Package)

	file, err := os.Create(args.OutFile)
	if err != nil {
		log.Error("Could not open output file", "err", err)
		os.Exit(1)
	}
	defer file.Close()
	if _, err = file.WriteString(code); err != nil {
		log.Error("Could not write to output file", "err", err)
		os.Exit(1)
	}
}

func Main(args CliArgs) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     args.LogLevel,
	}))
	wc := &WarningsCollector{warnings: make([]error, 0)}

	log.Info("Collecting all language files")
	walker, err := IoDirWalker(args.MessagesDirectory)
	if err != nil {
		log.Error("Could not collect all files in the messages directory", "err", err)
		os.Exit(1)
	}

	log.Info("Parsing files")
	messages, err := ParseJson(walker, wc)
	if err != nil {
		log.Error("Could not parse all files in the messages directory", "err", err)
		os.Exit(1)
	}
	if !wc.IsEmpty() {
		for _, warning := range wc.Warnings() {
			log.Warn(warning.Error())
		}
		os.Exit(1)
	}

	log.Info("Verifying that all languages have all the keys and there is a default language")
	allLanguages := messages.Languages()
	if !allLanguages.Contains(args.DefaultLanguage) {
		log.Error("The default language does not exist", "default-language", args.DefaultLanguage, "found-languages", allLanguages)
		os.Exit(1)
	}

	if removedEntries := messages.RemoveEntriesWithoutLang(args.DefaultLanguage); len(removedEntries) > 0 {
		log.Warn("The following entries are not present in the default language. Ignoring them",
			"ignored-entries", Map(removedEntries, func(t *MessageEntry) string { return strings.Join((*t).FullPath(), ".") }))
	}

	if messages.EnsureAllLanguagesPresent(args.DefaultLanguage, allLanguages.Get()) {
		log.Warn("Some entries had not all the messages filled. Using the message from the default language")
	}

	log.Info("Generating code")
	if err := WriteCode(messages, args); err != nil {
		log.Error("Could not write code", "err", err)
		os.Exit(1)
	}
}

package cli

import (
	"log/slog"
	"os"
)

type CliArgs struct {
	MessagesDirectory string
	DefaultLanguage   string
	OutFile           string
	Package           string
	LogLevel          slog.Level
}

func Main(args CliArgs) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     args.LogLevel,
	}))

	log.Info("Collecting all language files")
	walker, err := IoDirWalker(args.MessagesDirectory)
	if err != nil {
		log.Error("Could not collect all files in the messages directory", "err", err)
		os.Exit(1)
	}

	log.Info("Parsing files")
	messages, err := ParseJson(walker)
	if err != nil {
		log.Error("Could not parse all files in the messages directory", "err", err)
		os.Exit(1)
	}

	log.Info("Verifying that all languages have all the keys and there is a default language")
	allLanguages := messages.Languages()
	if !allLanguages.Contains(args.DefaultLanguage) {
		log.Error("The default language does not exist", "default-language", args.DefaultLanguage, "found-languages", allLanguages)
		os.Exit(1)
	}

	if removedEntries := messages.RemoveEntriesWithoutLang(args.DefaultLanguage); len(removedEntries) > 0 {
		log.Warn("The following entries are not present in the default language. Ignoring them", "ignored-entries", removedEntries)
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

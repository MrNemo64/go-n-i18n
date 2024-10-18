package cli

import (
	"log/slog"
	"os"
)

type CliArgs struct {
	MessagesDirectory string
	DefaultLanguage   string
	LogLevel          slog.Level
}

func Main(args CliArgs) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
		Level:     args.LogLevel,
	}))

	walker, err := IoDirWalker(args.MessagesDirectory)
	if err != nil {
		log.Error("Could not collect all files in the messages directory", "err", err)
		os.Exit(1)
	}

	ParseJson(walker, log)
}

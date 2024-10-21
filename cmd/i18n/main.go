package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/MrNemo64/go-n-i18n/internal/cli"
)

func main() {
	defaultLanguage := flag.String("default-language", "", "Specifies the default language")
	messagesDir := flag.String("messages", "", "Specifies the directory with the files with the messages")
	outFile := flag.String("out-file", "generated_lang.go", "Specifies the output file with the messages")
	outPackage := flag.String("out-package", os.Getenv("GOPACKAGE"), "Specifies the output file with the messages")
	flag.Parse()

	if *defaultLanguage == "" || *messagesDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	cli.Main(cli.CliArgs{
		MessagesDirectory: *messagesDir,
		DefaultLanguage:   *defaultLanguage,
		OutFile:           *outFile,
		Package:           *outPackage,
		LogLevel:          slog.LevelDebug,
	})
}

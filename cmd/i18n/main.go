package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/MrNemo64/go-n-i18n/internal/cli"
)

func main() {
	defaultLanguage := flag.String("default-language", "", "Specifies the default language")
	messagesDir := flag.String("messages", "", "Specifies the directory with the files with the messages")

	flag.Parse()

	if *defaultLanguage == "" || *messagesDir == "" {
		fmt.Println("Usage: code-generator -default-language <lang> -messages <dir>")
		os.Exit(1)
	}

	cli.Main(cli.CliArgs{
		MessagesDirectory: *messagesDir,
		DefaultLanguage:   *defaultLanguage,
		LogLevel:          slog.LevelDebug,
	})
}

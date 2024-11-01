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
	outFile := flag.String("out-file", "generated_lang.go", "Specifies the output file with the messages")
	outPackage := flag.String("out-package", os.Getenv("GOPACKAGE"), "Specifies the output package name")
	topInterfaceName := flag.String("top-interface-name", "messages", "Specifies the name for the top level interface")
	publicNonNamedInterfaces := flag.Bool("public-non-named-interfaces", false, "Specifies that all generated interfaces should be public, even non named ones")
	flag.Parse()

	if *defaultLanguage == "" || *messagesDir == "" || *outFile == "" || *outPackage == "" || *topInterfaceName == "" {
		flag.Usage()
		fmt.Println("Version v0.0.3")
		os.Exit(1)
	}

	cli.Run(cli.CliArgs{
		MessagesDirectory:        *messagesDir,
		DefaultLanguage:          *defaultLanguage,
		OutFile:                  *outFile,
		Package:                  *outPackage,
		TopLevelInterfaceName:    *topInterfaceName,
		PublicNonNamedInterfaces: *publicNonNamedInterfaces,
		LogLevel:                 slog.LevelDebug,
	})
}

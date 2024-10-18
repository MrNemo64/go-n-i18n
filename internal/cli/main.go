package cli

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	messagecollector "github.com/MrNemo64/go-n-i18n/internal/cli/message_collector"
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

	allMessages, err := messagecollector.JsonMessageScanner{}.FindAllMessagesInDir(args.MessagesDirectory)
	if err != nil {
		log.Error(fmt.Sprintf("Could not collect all the messages in directory '%s': %s", args.MessagesDirectory, err.Error()))
		os.Exit(1)
	}

	for _, v := range allMessages {
		checkDuplicatedKeys(v, log)
	}
	for _, v := range allMessages {
		checkKeys(v, log)
	}

	defaultLanguage, foundDelfaultLanguage := allMessages[args.DefaultLanguage]
	if !foundDelfaultLanguage {
		log.Error(fmt.Sprintf("Could not find any message for the default language '%s'", args.DefaultLanguage))
		os.Exit(1)
	}

	for _, v := range allMessages {
		if v == defaultLanguage {
			continue
		}
		checkHasKeys(defaultLanguage, v, log)
	}

	for _, v := range allMessages {
		normalizeKeys(v)
	}
}

func checkHasKeys(reference *messagecollector.CollectedMessages, cm *messagecollector.CollectedMessages, log *slog.Logger) {
	for keyInReference, messageInReference := range reference.Messages {
		if _, cmHasKey := cm.Messages[keyInReference]; !cmHasKey {
			log.Warn(fmt.Sprintf("The language %s is missing the key '%s'. Using the key from %s", cm.LanguageTag, keyInReference, reference.LanguageTag))
			cm.Messages[keyInReference] = &messagecollector.MessageInstance{
				Message:    messageInReference.Message,
				TimesFound: 1,
			}
		}
	}
	var keysToDelete []string
	for keyInCm := range cm.Messages {
		if _, referenceHasKey := reference.Messages[keyInCm]; !referenceHasKey {
			log.Warn(fmt.Sprintf(fmt.Sprintf("The language %s has an extra key '%s' that %s does not have. Ignoring it", cm.LanguageTag, keyInCm, reference.LanguageTag)))
			keysToDelete = append(keysToDelete, keyInCm)
		}
	}
	for _, key := range keysToDelete {
		delete(cm.Messages, key)
	}
}

func normalizeKeys(cm *messagecollector.CollectedMessages) {
	normalizer := KeyNormalizer()
	newMap := make(map[string]*messagecollector.MessageInstance, len(cm.Messages))
	for k, v := range cm.Messages {
		newMap[normalizer.Normalize(k)] = v
	}
	cm.Messages = newMap
}

func checkDuplicatedKeys(cm *messagecollector.CollectedMessages, log *slog.Logger) {
	duplicates := cm.FindDuplicatedKeys()
	if len(duplicates) == 1 {
		log.Error(fmt.Sprintf("The language '%s' has a duplicated key: %s", cm.LanguageTag, duplicates[0]))
		os.Exit(1)
	} else if len(duplicates) > 1 {
		log.Error(fmt.Sprintf("The language '%s' has several duplicated keys: %s", cm.LanguageTag, strings.Join(duplicates, ", ")))
		os.Exit(1)
	}
}

func checkKeys(cm *messagecollector.CollectedMessages, log *slog.Logger) {
	var invalidKeys []string
	validator := KeyValidator()
	for key := range cm.Messages {
		if !validator.IsValidKey(key) {
			invalidKeys = append(invalidKeys, key)
		}
	}
	if len(invalidKeys) > 0 {
		log.Error(fmt.Sprintf("The language '%s' has invalid keys: %s", cm.LanguageTag, strings.Join(invalidKeys, ", ")))
		os.Exit(1)
	}
}

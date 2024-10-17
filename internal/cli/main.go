package cli

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	messagecollector "github.com/MrNemo64/go-n-i18n/internal/cli/message_collector"
	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
)

type CliArgs struct {
	MessagesDirectory string
	DefaultLanguage   string
	LogOptions        *slog.HandlerOptions
}

func Main(args CliArgs) {
	log := slog.New(slog.NewTextHandler(os.Stdout, args.LogOptions))

	allMessages, err := messagecollector.JsonMessageScanner{}.FindAllMessagesInDir(args.MessagesDirectory)
	if err != nil {
		util.Exit(1, fmt.Sprintf("Could not collect all the messages in directory '%s': %s", args.MessagesDirectory, err.Error()))
	}

	for _, v := range allMessages {
		checkDuplicatedKeys(v)
	}
	for _, v := range allMessages {
		checkKeys(v)
	}

	defaultLanguage, foundDelfaultLanguage := allMessages[args.DefaultLanguage]
	if !foundDelfaultLanguage {
		util.Exit(1, fmt.Sprintf("Could not find any message for the default language '%s'", args.DefaultLanguage))
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
			log.Warn(fmt.Sprintf("The language %s has an extra key '%s' that %s does not have. Ignoring it", cm.LanguageTag, keyInCm, reference.LanguageTag))
			keysToDelete = append(keysToDelete, keyInCm) // is it safe to delete as i iterate?
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

func checkDuplicatedKeys(cm *messagecollector.CollectedMessages) {
	duplicates := cm.FindDuplicatedKeys()
	if len(duplicates) == 1 {
		util.Exit(1, fmt.Sprintf("The language '%s' has a duplicated key: %s", cm.LanguageTag, duplicates[0]))
	} else if len(duplicates) > 1 {
		util.Exit(1, fmt.Sprintf("The language '%s' has several duplicated keys: %s", cm.LanguageTag, strings.Join(duplicates, ", ")))
	}
}

func checkKeys(cm *messagecollector.CollectedMessages) {
	var invalidKeys []string
	validator := KeyValidator()
	for key := range cm.Messages {
		if !validator.IsValidKey(key) {
			invalidKeys = append(invalidKeys, key)
		}
	}
	if len(invalidKeys) > 0 {
		util.Exit(1, fmt.Sprintf("The language '%s' has invalid keys: %s", cm.LanguageTag, strings.Join(invalidKeys, ", ")))
	}
}

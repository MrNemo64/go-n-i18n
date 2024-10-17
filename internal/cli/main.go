package cli

import (
	"fmt"
	"strings"

	messagecollector "github.com/MrNemo64/go-n-i18n/internal/cli/message_collector"
	"github.com/MrNemo64/go-n-i18n/internal/cli/util"
)

type CliArgs struct {
	MessagesDirectory string
	DefaultLanguage   string
}

func Main(args CliArgs) {
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

	if _, found := allMessages[args.DefaultLanguage]; !found {
		util.Exit(1, fmt.Sprintf("Could not find any message for the default language '%s'", args.DefaultLanguage))
	}

	for _, v := range allMessages {
		normalizeKeys(v)
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

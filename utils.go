package gero

import (
	"strings"
)

func processPrefix(mention, message string, prefixConfig *PrefixConfig) (string, bool) {
	message = strings.TrimSpace(message)

	if !prefixConfig.CheckForPrefix {
		return message, true
	}

	prefixes := prefixConfig.Prefixes

	if prefixConfig.ConsiderMentionPrefix {
		prefixes = append(prefixes, mention)
	}

	if len(prefixes) > 0 {
		for _, prefix := range prefixes {
			if strings.HasPrefix(message, prefix) {
				trimmedMessage := strings.TrimPrefix(message, prefix)
				trimmedMessage = strings.TrimSpace(trimmedMessage)
				return trimmedMessage, true
			}
		}
	}

	return "", false
}

package gero

import (
	"fmt"
	"strings"
)

func removePrefix(userId, message string, prefixConfig *PrefixConfig) (string, bool) {
	message = strings.TrimSpace(message)

	if !prefixConfig.CheckForPrefix {
		return message, true
	}

	if prefixConfig.ConsiderMentionPrefix {
		prefix := fmt.Sprintf("<@%s>", userId)
		if strings.HasPrefix(message, prefix) {
			newMessage := strings.TrimPrefix(message, prefix)
			return strings.TrimSpace(newMessage), true
		}
	}

	if len(prefixConfig.Prefixes) > 0 {
		for _, prefix := range prefixConfig.Prefixes {
			if strings.HasPrefix(message, prefix) {
				newMessage := strings.TrimPrefix(message, prefix)
				return newMessage, true
			}
		}
	}

	return message, false
}

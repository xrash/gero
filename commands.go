package gero

import (
	"github.com/bwmarrin/discordgo"
	"regexp"
)

type CommandHandler func(b *Bot, m *discordgo.MessageCreate, params []string) error

type PrefixConfig struct {
	CheckForPrefix        bool
	Prefixes              []string
	ConsiderMentionPrefix bool
}

func (b *Bot) HandleCommands(prefixConfig *PrefixConfig, commands map[string]CommandHandler) error {
	b.logger.Info("Configuring commands handler")

	regexps := make(map[string]*regexp.Regexp)
	for pattern, _ := range commands {
		r, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}

		regexps[pattern] = r
	}

	b.session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore messages by the author himself.
		if m.Author.ID == s.State.User.ID {
			return
		}

		command, ok := removePrefix(s.State.User.ID, m.Content, prefixConfig)
		if !ok {
			return
		}

		for pattern, handler := range commands {
			r := regexps[pattern]
			found := r.FindStringSubmatch(command)

			if len(found) < 1 {
				continue
			}

			err := handler(b, m, found[1:])
			if err != nil {
				b.logger.Error("Error in command handler [%s]: %s", pattern, err)
			}

			return
		}
	})

	return nil
}

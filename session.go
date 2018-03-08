package gero

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) CreateSession(oauth2token string) error {
	b.logger.Info("Creating discord session")

	s, err := discordgo.New(fmt.Sprintf("Bot %s", oauth2token))
	if err != nil {
		b.logger.Error("Error creating session: %s", err)
		return err
	}

	b.session = s

	return nil
}

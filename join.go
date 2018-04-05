package gero

import (
	"github.com/bwmarrin/discordgo"
)

type OnJoinCallback func(*discordgo.GuildCreate, *discordgo.Channel) error

func (b *Bot) OnJoin(callback OnJoinCallback) {
	b.session.AddHandler(func(s *discordgo.Session, gc *discordgo.GuildCreate) {
		var chosen *discordgo.Channel
		for _, c := range gc.Channels {
			if c.Type == discordgo.ChannelTypeGuildText {
				if chosen == nil {
					chosen = c
					continue
				}

				if c.Position < chosen.Position {
					chosen = c
				}
			}
		}

		if chosen != nil {
			if err := callback(gc, chosen); err != nil {
				b.logger.Error("Error running OnJoin callback: %s", err)
			}
		}
	})
}

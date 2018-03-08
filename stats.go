package gero

import (
	"github.com/bwmarrin/discordgo"
	"github.com/xrash/gero/stats"
	"time"
)

func (b *Bot) ActivateStatsCollector() error {
	b.logger.Info("Activating real-time stats collector")

	b.stats = stats.NewStats(100)
	go b.stats.ProcessActions()

	b.session.AddHandler(func(s *discordgo.Session, gc *discordgo.GuildCreate) {
		b.stats.PushAction(&stats.CreateGuildAction{
			GuildId:     gc.ID,
			GuildName:   gc.Name,
			GuildRegion: gc.Region,
			CreatedAt:   time.Now(),
		})
	})

	b.session.AddHandler(func(s *discordgo.Session, gd *discordgo.GuildDelete) {
		b.stats.PushAction(&stats.DeleteGuildAction{
			GuildId: gd.ID,
		})
	})

	return nil
}

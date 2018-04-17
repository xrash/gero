package gero

import (
	"github.com/bwmarrin/discordgo"
	"github.com/xrash/gero/msg"
	"github.com/xrash/gero/stats"
	"github.com/xrash/gol"
)

type Bot struct {
	// Session management.
	oauth2token string
	session     *discordgo.Session

	// Logger.
	logger *gol.Logger

	// Messenger.
	messenger *msg.Messenger

	// Current stats.
	stats *stats.Stats

	// Exit channel.
	exitCleaners []ExitCleaner
	exitChannel  chan string

	// Status Manager.
	statusManager *StatusManager
}

func NewBot() *Bot {
	return &Bot{
		exitCleaners:  make([]ExitCleaner, 0),
		exitChannel:   make(chan string, 1),
		logger:        gol.NewLogger(),
		statusManager: NewStatusManager(),
	}
}

func (b *Bot) Logger() Logger {
	return b.logger
}

func (b *Bot) Session() *discordgo.Session {
	return b.session
}

func (b *Bot) Stats() *stats.Stats {
	return b.stats
}

func (b *Bot) Msg() *msg.Messenger {
	return b.messenger
}

func (b *Bot) Status() *StatusManager {
	return b.statusManager.WithBot(b)
}

func (b *Bot) RegisterExitCleaner(ec ExitCleaner) {
	b.exitCleaners = append(b.exitCleaners, ec)
}

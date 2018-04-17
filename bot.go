package gero

import (
	"github.com/bwmarrin/discordgo"
	"github.com/xrash/gero/msg"
	"github.com/xrash/gero/stats"
	"github.com/xrash/gol"
)

type ExitCleaner func()

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
}

func NewBot() *Bot {
	return &Bot{
		exitCleaners: make([]ExitCleaner, 0),
		exitChannel:  make(chan string, 1),
		logger:       gol.NewLogger(),
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

func (b *Bot) SetStatus(s string) error {
	return b.session.UpdateStatus(0, s)
}

func (b *Bot) RegisterExitCleaner(ec ExitCleaner) {
	b.exitCleaners = append(b.exitCleaners, ec)
}

func (b *Bot) gracefullyExit() {
	for _, ec := range b.exitCleaners {
		ec()
	}
}

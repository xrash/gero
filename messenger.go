package gero

import (
	"github.com/xrash/gero/msg"
)

func (b *Bot) CreateMessenger(color int) error {
	b.logger.Info("Creating messenger")
	b.messenger = msg.NewMessenger(b.Session(), b.Logger(), color)
	return nil
}

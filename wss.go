package gero

func (b *Bot) OpenWssConnection() error {
	b.logger.Info("Opening wss connection")

	err := b.session.Open()
	if err != nil {
		b.logger.Error("Error opening wss connection: %s", err)
		return err
	}

	b.logger.Info("Wss connection successfully open")

	return nil
}

func (b *Bot) CloseWssConnection() error {
	b.logger.Info("Closing wss connection")

	err := b.session.Close()
	if err != nil {
		b.logger.Error("Error closing wss connection: %s", err)
	}

	return err
}

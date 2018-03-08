package gero

import (
	"os"
	"os/signal"
)

func (b *Bot) WaitProgramExit() {
	b.logger.Info("Waiting for program exit trigger")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)

	select {
	case s := <-signalChannel:
		b.logger.Notice("Program exiting because of signal: %s", s)
	case r := <-b.exitChannel:
		b.logger.Notice("Program purposely exiting: %s", r)
	}

	b.gracefullyExit()
}

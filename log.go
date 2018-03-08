package gero

import (
	"github.com/xrash/gol"
	"github.com/xrash/gol/formatters"
	"github.com/xrash/gol/handlers"
)

type Logger interface {
	Info(format string, params ...interface{}) []error
	Error(format string, params ...interface{}) []error
}

func (b *Bot) ConfigureLogger(verbose bool, logfile string) error {
	if verbose {
		b.ConfigureOutput()
	}

	if logfile != "" {
		if err := b.ConfigureLogfile(logfile); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) ConfigureOutput() {
	basicFormatter := formatters.NewBasicFormatter()
	stdoutHandler := handlers.NewStdoutHandler()
	b.logger.AddHandler(stdoutHandler, basicFormatter, gol.LEVEL_DEBUG)
	b.logger.Info("Setting verbose output")
}

func (b *Bot) ConfigureLogfile(logfile string) error {
	basicFormatter := formatters.NewBasicFormatter()
	fileHandler := handlers.NewFileHandler(logfile)

	if err := fileHandler.Open(); err != nil {
		return err
	}

	b.logger.AddHandler(fileHandler, basicFormatter, gol.LEVEL_DEBUG)

	b.logger.Info("Setting logfile %s", logfile)

	return nil
}

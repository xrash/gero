package gero

import (
	"github.com/xrash/gol"
	"github.com/xrash/gol/formatters"
	"github.com/xrash/gol/handlers"
)

type Logger interface {
	Info(format string, params ...interface{}) []error
	Error(format string, params ...interface{}) []error
	Debug(format string, params ...interface{}) []error
}

func (b *Bot) ConfigureLogger(verbose bool, logfile string, level string) error {
	if verbose {
		b.ConfigureOutput(level)
	}

	if logfile != "" {
		if err := b.ConfigureLogfile(logfile, level); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) ConfigureOutput(level string) {
	basicFormatter := formatters.NewBasicFormatter()
	stdoutHandler := handlers.NewStdoutHandler()
	b.logger.AddHandler(stdoutHandler, basicFormatter, gol.NewLogLevel(level))
	b.logger.Info("Setting verbose output")
}

func (b *Bot) ConfigureLogfile(logfile string, level string) error {
	basicFormatter := formatters.NewBasicFormatter()
	fileHandler := handlers.NewFileHandler(logfile)

	if err := fileHandler.Open(); err != nil {
		return err
	}

	b.logger.AddHandler(fileHandler, basicFormatter, gol.NewLogLevel(level))

	b.logger.Info("Setting logfile %s", logfile)

	return nil
}

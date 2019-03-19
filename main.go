package main

import (
	"xl/app"
	"xl/log"
	"xl/ui/termbox"

	"flag"

	"go.uber.org/zap"
)

func main() {
	logger := newLogger([]string{"xl.log"})
	defer func() {
		_ = logger.Sync()
	}()

	log.L = logger

	logger.Info("application starting")

	t := termbox.New()
	defer t.Close()

	a := app.New(&app.Config{
		Logger: logger,
		Input:  t.Input(),
		Output: t.Output(),
	})

	flag.Parse()
	args := flag.Args()

	if len(args) > 0 {
		err := a.OpenDocument(args[0])
		if err != nil {
			panic(err)
		}
	} else {
		a.ResetDocument()
	}

	a.Loop()

	logger.Info("application finishing")
}

func newLogger(logFiles []string) *zap.Logger {
	config := zap.Config{
		Level:         zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:      "console",
		EncoderConfig: zap.NewDevelopmentEncoderConfig(),
		OutputPaths:   logFiles,
	}
	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	return logger
}

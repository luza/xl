package app

import (
	"errors"
	"xl/document"
	"xl/fs"
	"xl/fs/bufcsv"
	"xl/ui"

	"go.uber.org/zap"
)

type App struct {
	ui.DataDelegateInterface

	logger *zap.Logger
	input  ui.InputInterface
	output ui.OutputInterface
	doc    *document.Document
	file   fs.FileInterface
}

type Config struct {
	Logger *zap.Logger
	Input  ui.InputInterface
	Output ui.OutputInterface
}

func New(config *Config) *App {
	a := &App{
		logger: config.Logger,
		input:  config.Input,
		output: config.Output,
	}
	a.output.SetDataDelegate(a)
	return a
}

func (a *App) ResetDocument() {
	a.doc = document.NewWithEmptySheet()
	a.output.SetDataDelegate(a)
	a.output.RefreshView()
}

func (a *App) OpenDocument(filename string) error {
	a.file = bufcsv.NewWithFilename(filename)
	var err error
	a.doc, err = a.file.Open()
	if err != nil {
		return err
	}
	if len(a.doc.Sheets) == 0 {
		return errors.New("no sheets at file open")
	}
	if a.doc.CurrentSheet == nil {
		a.doc.CurrentSheet = a.doc.Sheets[0]
		a.doc.CurrentSheetN = 0
	}
	a.output.RefreshView()
	return nil
}

func (a *App) Write() error {
	if a.file == nil {
		return errors.New("no file name")
	}
	return a.WriteAs("")
}

func (a *App) WriteAs(filename string) error {
	if filename != "" {
		a.file = bufcsv.NewWithFilename(filename)
	}
	return a.file.Write(a.doc)
}

func (a *App) Loop() {
	for {
		event, err := a.input.ReadKey()
		if err != nil {
			a.logger.Error("input read error: " + err.Error())
			return
		}
		if keyEvent, ok := event.(ui.KeyEvent); ok {
			stop := a.processKeyEvent(keyEvent)
			if stop {
				return
			}
		} else {
			a.logger.Error("unknown input event")
			return
		}
	}
}

func (a *App) ShowError(err error) {
	a.output.SetStatus(err.Error(), ui.StatusFlagError)
}

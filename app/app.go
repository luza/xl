package app

import (
	"xl/document"
	"xl/document/sheet"
	"xl/fs"
	"xl/fs/bufcsv"
	"xl/fs/bufxlsx"
	"xl/ui"

	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gdamore/tcell"
	"go.uber.org/zap"
)

const (
	rcFile = ".xlrc"
)

type App struct {
	ui.DataDelegateInterface
	screen  tcell.Screen
	logger  *zap.Logger
	input   ui.InputInterface
	output  ui.OutputInterface
	doc     *document.Document
	file    fs.FileInterface
	hotKeys map[Key]string

	// Keeps the cell for copy/cut/paste operations.
	cellBuffer *sheet.Cell
}

type Config struct {
	Screen tcell.Screen
	Logger *zap.Logger
	Input  ui.InputInterface
	Output ui.OutputInterface
}

func New(config *Config) *App {
	a := &App{
		screen:  config.Screen,
		logger:  config.Logger,
		input:   config.Input,
		output:  config.Output,
		hotKeys: make(map[Key]string),
	}
	a.output.SetDataDelegate(a)
	a.loadRC()
	return a
}

// ResetDocument creates a new empty document.
func (a *App) ResetDocument() {
	a.doc = document.NewWithEmptySheet()
	a.output.SetDataDelegate(a)
	a.output.RefreshView()
}

// OpenDocument reads document from file with given name.
func (a *App) OpenDocument(filename string) error {
	a.file = guessFileFormat(filename)
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

// Write writes document to the same file it was read from.
func (a *App) Write() error {
	if a.file == nil {
		return errors.New("no file name")
	}
	return a.WriteAs("")
}

// WriteAs writes document to file with given name.
func (a *App) WriteAs(filename string) error {
	if filename != "" {
		a.file = bufcsv.NewWithFilename(filename)
	}
	return a.file.Write(a.doc)
}

// Loop is the main loop, reads and processes key presses.
func (a *App) Loop() {
	for {
		ev := a.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			e := ui.KeyEvent{

				Mod: ev.Modifiers(),
				Key: ev.Key(),
				Ch:  ev.Rune(),
			}
			stop := a.processKeyEvent(e)
			if stop {
				return
			}
		case *tcell.EventResize:
			a.output.RefreshView()
		case *tcell.EventMouse:
			//handling event mouse
		case *tcell.EventError:
			a.logger.Error("unknown input event")
			return
		}
	}
}

// showErrors displays error message in status line.
func (a *App) showError(err error) {
	a.output.SetStatus(err.Error(), ui.StatusFlagError)
}

// loadRC reads rc file containing commands to be executed on launch.
func (a *App) loadRC() {
	rcLocation := os.Getenv("HOME") + "/" + rcFile
	data, err := ioutil.ReadFile(rcLocation)
	if err != nil {
		a.showError(err)
		return
	}
	for _, l := range strings.Split(string(data), "\n") {
		if l == "" {
			continue
		}
		a.processCommand(l)
	}
}

func guessFileFormat(filename string) fs.FileInterface {
	if strings.HasSuffix(filename, ".xlsx") {
		return bufxlsx.NewWithFilename(filename)
	} else {
		return bufcsv.NewWithFilename(filename)
	}
}

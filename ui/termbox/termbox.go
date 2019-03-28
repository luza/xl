package termbox

import (
	"xl/ui"

	"github.com/gdamore/tcell/termbox"
)

type Termbox struct {
	ui.InputInterface
	ui.OutputInterface

	dataDelegate ui.DataDelegateInterface

	// Value of termbox.Size()
	screenWidth  int
	screenHeight int

	// How many rows and columns are visible for last drawing iteration.
	calculatedViewportWidth  int
	calculatedViewportHeight int

	// Length in chars of vertical ruler for last drawing iteration.
	vRulerWidth int

	// Cursor position for last drawing iteration.
	lastCursorX int
	lastCursorY int

	// What need to redrawn on next draw iteration.
	dirty ui.DirtyFlag

	// Message displaying in status line and its decoration flags.
	statusMessage string
	statusFlags   int
}

func New() *Termbox {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	termbox.SetOutputMode(termbox.Output256)
	width, height := termbox.Size()
	return &Termbox{
		screenWidth:  width,
		screenHeight: height,
		dirty:        ui.DirtyHRuler | ui.DirtyVRuler | ui.DirtyGrid | ui.DirtyFormulaLine | ui.DirtyStatusLine,
	}
}

func (t *Termbox) Close() {
	termbox.Close()
}

func (t *Termbox) Input() ui.InputInterface {
	return t
}

func (t *Termbox) Output() ui.OutputInterface {
	return t
}

func (t *Termbox) ViewportHeight() int {
	return t.calculatedViewportHeight
}

func (t *Termbox) ViewportWidth() int {
	return t.calculatedViewportWidth
}

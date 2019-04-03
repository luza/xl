package termbox

import (
	"xl/ui"

	"github.com/gdamore/tcell"
)

type Termbox struct {
	ui.InputInterface
	ui.OutputInterface
	Screen       tcell.Screen
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
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	var s tcell.Screen
	var e error
	if s, e = tcell.NewScreen(); e != nil {
		panic(e)
	} else if e = s.Init(); e != nil {
		panic(e)
	}

	width, height := s.Size()
	return &Termbox{
		Screen:       s,
		screenWidth:  width,
		screenHeight: height,
		dirty:        ui.DirtyHRuler | ui.DirtyVRuler | ui.DirtyGrid | ui.DirtyFormulaLine | ui.DirtyStatusLine,
	}
}

func (t *Termbox) Close() {
	t.Screen.Fini()
}

func (t *Termbox) GetScreen() tcell.Screen {
	return t.Screen
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

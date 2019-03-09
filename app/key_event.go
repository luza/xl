package app

import (
	"fmt"
	"xl/document/sheet"
	"xl/ui"

	"github.com/nsf/termbox-go"
)

// processKeyEvent does the job associated with the key press.
func (a *App) processKeyEvent(event ui.KeyEvent) bool {
	switch event.Ch {
	case ':':
		stop := a.inputCommand()
		a.output.RefreshView()
		return stop
	}

	switch event.Key {
	case termbox.KeyCtrlC:
		return true
	case termbox.KeySpace:
		a.pageDown()
		a.output.RefreshView()
	case termbox.KeyArrowUp:
		a.moveCursorUp()
		a.output.RefreshView()
	case termbox.KeyArrowDown:
		a.moveCursorDown()
		a.output.RefreshView()
	case termbox.KeyArrowLeft:
		a.moveCursorLeft()
		a.output.RefreshView()
	case termbox.KeyArrowRight, termbox.KeyTab:
		a.moveCursorRight()
		a.output.RefreshView()
	case termbox.KeyEnter:
		a.editCell()
		a.output.RefreshView()
	default:
		a.output.SetStatus(fmt.Sprintf("ch: %v, key: %v", event.Ch, event.Key), 0)
		a.output.RefreshView()
	}

	return false
}

// moveCursorLeft moves cursor up on one cell.
func (a *App) moveCursorUp() bool {
	if a.doc.CurrentSheet.Cursor.Y <= 0 {
		return false
	}
	a.doc.CurrentSheet.Cursor.Y--
	if a.doc.CurrentSheet.Cursor.Y < a.doc.CurrentSheet.Viewport.Top {
		a.doc.CurrentSheet.Viewport.Top--
	}
	a.output.SetDirty(ui.DirtyVRuler | ui.DirtyGrid | ui.DirtyFormulaLine)
	return true
}

// moveCursorLeft moves cursor down on one cell.
func (a *App) moveCursorDown() bool {
	a.doc.CurrentSheet.Cursor.Y++
	if a.doc.CurrentSheet.Cursor.Y >= a.doc.CurrentSheet.Viewport.Top+a.output.ViewportHeight() {
		a.doc.CurrentSheet.Viewport.Top++
	}
	a.output.SetDirty(ui.DirtyVRuler | ui.DirtyGrid | ui.DirtyFormulaLine)
	return true
}

// moveCursorLeft moves cursor left on one cell.
func (a *App) moveCursorLeft() bool {
	if a.doc.CurrentSheet.Cursor.X <= 0 {
		return false
	}
	a.doc.CurrentSheet.Cursor.X--
	if a.doc.CurrentSheet.Cursor.X < a.doc.CurrentSheet.Viewport.Left {
		a.doc.CurrentSheet.Viewport.Left--
	}
	a.output.SetDirty(ui.DirtyHRuler | ui.DirtyGrid | ui.DirtyFormulaLine)
	return true
}

// moveCursorRight moves cursor right on one cell.
func (a *App) moveCursorRight() bool {
	a.doc.CurrentSheet.Cursor.X++
	if a.doc.CurrentSheet.Cursor.X >= a.doc.CurrentSheet.Viewport.Left+a.output.ViewportWidth() {
		a.doc.CurrentSheet.Viewport.Left++
	}
	a.output.SetDirty(ui.DirtyHRuler | ui.DirtyGrid | ui.DirtyFormulaLine)
	return true
}

// pageDown moves cursor down on number of lines equal to window height.
func (a *App) pageDown() bool {
	a.doc.CurrentSheet.Cursor.Y += a.output.ViewportHeight()
	a.doc.CurrentSheet.Viewport.Top += a.output.ViewportHeight()
	a.output.SetDirty(ui.DirtyHRuler | ui.DirtyVRuler | ui.DirtyGrid | ui.DirtyFormulaLine)
	return true
}

// inputCommand opens inline editor in status line, with ':' prompt.
// Once user finishes command input, processes the command.
func (a *App) inputCommand() bool {
	command, err := a.output.InputCommand()
	if err != nil {
		a.ShowError(err)
		return false
	}
	a.output.SetStatus("", 0)
	stop := a.processCommand(command)
	a.output.SetDirty(ui.DirtyStatusLine)
	return stop
}

// editCell opens an inline editor at formula line place with the current cell value.
// Once user exits editor (by Enter or Esc), writes new value to cell.
func (a *App) editCell() {
	cur := a.doc.CurrentSheet.Cursor
	cell := a.doc.CurrentSheet.Cell(cur.X, cur.Y)
	if cell == nil {
		cell = sheet.NewCellGeneral()
	}
	newValue, err := a.output.EditCellValue(cell.Value())
	if err != nil {
		a.logger.Error(err.Error())
		return
	}
	cell.SetValueText(newValue)
	a.doc.CurrentSheet.SetCell(cur.X, cur.Y, cell)
	a.output.SetDirty(ui.DirtyGrid | ui.DirtyFormulaLine)
}

package app

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strings"

	"xl/document/sheet"
	"xl/ui"
)

const colSizeIncrementStep = 6

// processCommand do the job associated with the command.
// If no such command found, shows the error in status line.
func (a *App) processCommand(c string) bool {
	c, args := parseArgs(c)
	switch c {
	case "q", "quit":
		return true
	case "w", "write":
		a.cmdWrite(arg1(args))
	case "wider":
		a.cmdResizeColumn(1)
	case "narrower":
		a.cmdResizeColumn(-1)
	case "appendSheet":
		a.cmdNewSheet(arg1(args))
	case "nextSheet":
		a.cmdNextSheet()
	case "bind":
		a.cmdBind(args)
	case "cutCell":
		a.cmdCutCell()
	case "pasteCell":
		a.cmdPasteCell()
	case "copyCell":
		a.cmdCopyCell()
	case "insertRow":
		a.cmdInsertRow(0)
	case "insertRowAfter":
		a.cmdInsertRow(1)
	case "insertCol":
		a.cmdInsertCol(0)
	case "insertColAfter":
		a.cmdInsertCol(1)
	case "deleteRow":
		a.cmdDeleteRow()
	case "deleteCol":
		a.cmdDeleteCol()
	case "mprof":
		a.cmdMemProf()
	default:
		a.output.SetStatus(fmt.Sprintf("unknown command %s", c), ui.StatusFlagError)
	}
	return false
}

// arg1 returns first argument or empty string.
func arg1(args []string) string {
	return argN(args, 1)
}

// argN returns Nth argument or empty string.
func argN(args []string, n int) string {
	if len(args) >= n {
		return args[n-1]
	}
	return ""
}

// parseArgs splits raw command line into command itself and list of command arguments.
// TODO: arguments can possibly be wrapped in quotes
func parseArgs(cmd string) (string, []string) {
	// FIXME: naive implementation
	c := strings.Split(cmd, " ")
	return c[0], c[1:]
}

// cmdResizeColumn resizes column under cursor so its width becomes given N pixels.
func (a *App) cmdResizeColumn(n int) {
	col := a.doc.CurrentSheet.Cursor.X
	size := a.doc.CurrentSheet.ColSize(col)
	a.doc.CurrentSheet.SetColSize(col, size+n*colSizeIncrementStep)
	a.output.SetDirty(ui.DirtyHRuler | ui.DirtyGrid)
}

// cmdWrite saves document to file.
func (a *App) cmdWrite(filename string) {
	var err error
	if filename != "" {
		err = a.WriteAs(filename)
	} else {
		err = a.Write()
	}
	if err != nil {
		a.showError(err)
	}
}

// cmdNewList creates a new sheet.
func (a *App) cmdNewSheet(title string) {
	_, err := a.doc.NewSheet(title)
	if err != nil {
		a.showError(err)
		return
	}
	a.output.SetDirty(ui.DirtyStatusLine)
}

// cmdNextSheet switches the current sheet to next one.
// If current sheet is the last one, it switches to first.
func (a *App) cmdNextSheet() {
	a.doc.CurrentSheetN++
	if a.doc.CurrentSheetN >= len(a.doc.Sheets) {
		a.doc.CurrentSheetN = 0
	}
	a.doc.CurrentSheet = a.doc.Sheets[a.doc.CurrentSheetN]
	a.output.SetDirty(ui.DirtyStatusLine | ui.DirtyGrid | ui.DirtyFormulaLine)
}

// cmdBind binds a command to a hot key.
func (a *App) cmdBind(args []string) {
	if len(args) < 2 {
		a.output.SetStatus("hot key and command must be specified", ui.StatusFlagError)
		return
	}
	k, ok := HotKeys[args[0]]
	if !ok {
		a.output.SetStatus(fmt.Sprintf("key %s is not a valid key", args[0]), ui.StatusFlagError)
		return
	}
	// TODO: escape with quotes
	a.hotKeys[k] = strings.Join(args[1:], " ")
}

// cmdCutCell erases the cell (but puts its value to buffer first).
func (a *App) cmdCutCell() {
	a.cmdCopyCell()
	a.doc.CurrentSheet.CellUnderCursor().EraseValue()
	a.output.SetDirty(ui.DirtyGrid | ui.DirtyFormulaLine)
}

// cmdCopyCell copies cell value to the buffer.
func (a *App) cmdCopyCell() {
	cell := a.doc.CurrentSheet.CellUnderCursor()
	if cell == nil {
		cell = sheet.NewCellEmpty()
	}
	cellCopy := *cell
	a.cellBuffer = &cellCopy
}

// cmdPasteCell replace cell under cursor with the value of previously copied or cut cell.
func (a *App) cmdPasteCell() {
	if a.cellBuffer == nil {
		a.output.SetStatus("buffer is empty", ui.StatusFlagError)
		return
	}
	cellCopy := *a.cellBuffer
	s := a.doc.CurrentSheet
	s.SetCell(s.Cursor.X, s.Cursor.Y, &cellCopy)
	a.output.SetDirty(ui.DirtyGrid | ui.DirtyFormulaLine)
}

func (a *App) cmdInsertRow(n int) {
	a.doc.CurrentSheet.Cursor.Y += n
	a.doc.CurrentSheet.InsertRow(a.doc.CurrentSheet.Cursor.Y)
	a.output.SetDirty(ui.DirtyGrid | ui.DirtyFormulaLine)
}

func (a *App) cmdInsertCol(n int) {
	a.doc.CurrentSheet.Cursor.X += n
	a.doc.CurrentSheet.InsertCol(a.doc.CurrentSheet.Cursor.X)
	a.output.SetDirty(ui.DirtyGrid | ui.DirtyFormulaLine)
}

func (a *App) cmdDeleteRow() {
	a.doc.CurrentSheet.DeleteRow(a.doc.CurrentSheet.Cursor.Y)
	a.output.SetDirty(ui.DirtyGrid | ui.DirtyFormulaLine)
}

func (a *App) cmdDeleteCol() {
	a.doc.CurrentSheet.DeleteCol(a.doc.CurrentSheet.Cursor.X)
	a.output.SetDirty(ui.DirtyGrid | ui.DirtyFormulaLine)
}

func (a *App) cmdMemProf() {
	f, err := os.Create("xl.mprof")
	if err != nil {
		return
	}
	_ = pprof.WriteHeapProfile(f)
	_ = f.Close()
}

package app

import (
	"xl/document/value"
	"xl/ui"

	"bytes"
	"strconv"
)

// Callbacks collection providing data to be displayed.

func (a *App) CellView(x, y int) *ui.CellView {
	c := a.doc.CurrentSheet.Cell(x, y)
	if c == nil {
		return &ui.CellView{
			Name: cellName(x, y),
		}
	}
	v, err := c.StringValue(value.NewEvalContext(a.doc))
	if err != nil {
		t := err.Error()
		return &ui.CellView{
			Name:  cellName(x, y),
			Error: &t,
		}
	}
	return &ui.CellView{
		Name:        cellName(x, y),
		DisplayText: v,
	}
}

func (a *App) RowView(n int) *ui.RowView {
	return &ui.RowView{
		Name:   rowName(n),
		Height: a.doc.CurrentSheet.RowSize(n),
	}
}

func (a *App) ColView(n int) *ui.ColView {
	return &ui.ColView{
		Name:  colName(n),
		Width: a.doc.CurrentSheet.ColSize(n),
	}
}

func (a *App) SheetView() *ui.SheetView {
	cell := a.doc.CurrentSheet.CellUnderCursor()
	var formulaText string
	if cell != nil {
		formulaText = cell.RawValue()
	}
	return &ui.SheetView{
		Name:     a.doc.CurrentSheet.Title,
		Cursor:   a.doc.CurrentSheet.Cursor,
		Viewport: a.doc.CurrentSheet.Viewport,
		FormulaLineView: ui.FormulaLineView{
			DisplayText: formulaText,
		},
	}
}

func (a *App) DocView() *ui.DocView {
	sheetNames := make([]string, len(a.doc.Sheets))
	currentSheetIdx := 0
	for i, s := range a.doc.Sheets {
		sheetNames[i] = s.Title
		if s == a.doc.CurrentSheet {
			currentSheetIdx = i
		}
	}
	return &ui.DocView{
		Sheets:          sheetNames,
		CurrentSheetIdx: currentSheetIdx,
	}
}

// rowName returns name of row for given index.
func rowName(n int) string {
	return strconv.Itoa(n + 1)
}

// colName returns name of column for given index.
func colName(n int) string {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var result bytes.Buffer
	result.WriteByte(alphabet[n%26])
	n /= 26
	for n > 0 {
		result.WriteByte(alphabet[(n-1)%26])
		n = (n - 1) / 26
	}
	// reverse bytes
	b := result.Bytes()
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

// cellName returns name for cell under given X and Y.
func cellName(x, y int) string {
	return colName(x) + rowName(y)
}

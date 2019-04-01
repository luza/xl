package app

import (
	"xl/document"
	"xl/document/eval"
	"xl/ui"
)

// Callbacks collection providing data to be displayed.

func (a *App) CellView(x, y int) *ui.CellView {
	c := a.doc.CurrentSheet.Cell(x, y)
	if c == nil {
		return &ui.CellView{
			Name: document.CellName(x, y),
		}
	}
	v, err := c.StringValue(eval.NewContext(a.doc, a.doc.CurrentSheet.Idx))
	if err != nil {
		t := err.Error()
		return &ui.CellView{
			Name:  document.CellName(x, y),
			Error: &t,
		}
	}
	return &ui.CellView{
		Name:        document.CellName(x, y),
		DisplayText: v,
		Expression:  c.Expression(eval.NewContext(a.doc, a.doc.CurrentSheet.Idx)),
	}
}

func (a *App) RowView(n int) *ui.RowView {
	return &ui.RowView{
		Name:   document.RowName(n),
		Height: a.doc.CurrentSheet.RowSize(n),
	}
}

func (a *App) ColView(n int) *ui.ColView {
	return &ui.ColView{
		Name:  document.ColName(n),
		Width: a.doc.CurrentSheet.ColSize(n),
	}
}

func (a *App) SheetView() *ui.SheetView {
	c := a.doc.CurrentSheet.CellUnderCursor()
	sv := &ui.SheetView{
		Name:     a.doc.CurrentSheet.Title,
		Cursor:   a.doc.CurrentSheet.Cursor,
		Viewport: a.doc.CurrentSheet.Viewport,
	}
	if c != nil {
		sv.FormulaLineView = ui.FormulaLineView{
			DisplayText: c.RawValue(),
			Expression:  c.Expression(eval.NewContext(a.doc, a.doc.CurrentSheet.Idx)),
		}
	}
	return sv
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

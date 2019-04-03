package ui

import (
	"xl/document/sheet"
	"xl/formula"

	"github.com/gdamore/tcell"
)

const (
	DirtyHRuler = 1 << iota
	DirtyVRuler
	DirtyGrid
	DirtyFormulaLine
	DirtyStatusLine
)

const (
	StatusFlagError = 1 << iota
)

type DirtyFlag int

type OutputInterface interface {
	SetDataDelegate(DataDelegateInterface)
	RefreshView()
	ViewportHeight() int
	ViewportWidth() int
	SetDirty(DirtyFlag)
	InputCommand() (string, error)
	EditCellValue(string) (string, error)
	SetStatus(string, int)
	Screen() tcell.Screen
}

type DataDelegateInterface interface {
	DocView() *DocView
	SheetView() *SheetView
	CellView(x, y int) *CellView
	RowView(n int) *RowView
	ColView(n int) *ColView
}

type CellView struct {
	Name        string
	DisplayText string
	Error       *string
	Expression  *formula.Expression
}

type RowView struct {
	Name   string
	Height int
}

type ColView struct {
	Name  string
	Width int
}

type SheetView struct {
	Name            string
	Cursor          sheet.Cursor
	Viewport        sheet.Viewport
	FormulaLineView FormulaLineView
}

type FormulaLineView struct {
	DisplayText string
	Expression  *formula.Expression
}

type DocView struct {
	Sheets          []string
	CurrentSheetIdx int
}

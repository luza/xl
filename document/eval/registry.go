package eval

import (
	"github.com/shopspring/decimal"
)

type RefRegistryInterface interface {
	NewCellRef(sheetTitle, cellName string) (*CellRef, error)
	NewRangeRef(sheetTitle, cellFromName, cellToName string) (*RangeRef, error)
	SheetTitle(sheetIdx int) (string, error)
	CellName(cell Cell) (string, error)
	Value(ec *Context, cell Cell) (Value, error)
	BoolValue(ec *Context, cell Cell) (bool, error)
	DecimalValue(ec *Context, cell Cell) (decimal.Decimal, error)
	StringValue(ec *Context, cell Cell) (string, error)
}

package eval

import (
	"github.com/shopspring/decimal"
)

type RefRegistryInterface interface {
	NewCellRef(cellName, sheetTitle string) (*CellRef, error)
	SheetTitle(r *CellRef) (string, error)
	CellName(r *CellRef) (string, error)
	Value(ec *Context, r *CellRef) (Value, error)
	BoolValue(ec *Context, r *CellRef) (bool, error)
	DecimalValue(ec *Context, r *CellRef) (decimal.Decimal, error)
	StringValue(ec *Context, r *CellRef) (string, error)
}

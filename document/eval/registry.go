package eval

import (
	"github.com/shopspring/decimal"
)

type CellAddress struct {
	SheetIdx int
	X        int
	Y        int
}

type CellReference struct {
	CellAddress
	AnchoredX bool
	AnchoredY bool
}

type RefRegistryInterface interface {
	AddRef(cell CellReference)
	FromAddress(cell CellReference) (string, string, error)
	ToAddress(sheetTitle, cellName string) (CellReference, error)
	Value(ec *Context, cell CellAddress) (Value, error)
	BoolValue(ec *Context, cell CellAddress) (bool, error)
	DecimalValue(ec *Context, cell CellAddress) (decimal.Decimal, error)
	StringValue(ec *Context, cell CellAddress) (string, error)
	IterateBoolValues(ec *Context, cell, cellTo CellAddress, f func(bool) error) error
	IterateDecimalValues(ec *Context, cell, cellTo CellAddress, f func(decimal.Decimal) error) error
	IterateStringValues(ec *Context, cell, cellTo CellAddress, f func(string) error) error
}

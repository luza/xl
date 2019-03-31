package eval

import (
	"github.com/shopspring/decimal"
)

type Axis struct {
	X int
	Y int
}

type CellRef struct {
	Value

	SheetIdx int
	Cell     Axis
}

func NewCellRef(sheetIdx int, cell Axis) *CellRef {
	return &CellRef{
		SheetIdx: sheetIdx,
		Cell:     cell,
	}
}

func (r *CellRef) Type(ec *Context) (int, error) {
	if ec.Visited(r) {
		return 0, NewError(ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(r)
	defer ec.ResetVisited(l)
	val, err := ec.DataProvider.Value(ec, r)
	if err != nil {
		return 0, err
	}
	return val.Type(ec)
}

func (r *CellRef) BoolValue(ec *Context) (bool, error) {
	if ec.Visited(r) {
		return false, NewError(ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(r)
	defer ec.ResetVisited(l)
	return ec.DataProvider.BoolValue(ec, r)
}

func (r *CellRef) DecimalValue(ec *Context) (decimal.Decimal, error) {
	if ec.Visited(r) {
		return decimal.Zero, NewError(ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(r)
	defer ec.ResetVisited(l)
	return ec.DataProvider.DecimalValue(ec, r)
}

func (r *CellRef) StringValue(ec *Context) (string, error) {
	if ec.Visited(r) {
		return "", NewError(ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(r)
	defer ec.ResetVisited(l)
	return ec.DataProvider.StringValue(ec, r)
}

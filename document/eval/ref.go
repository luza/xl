package eval

import (
	"github.com/shopspring/decimal"
)

type Cell struct {
	SheetIdx int
	X        int
	Y        int
}

type CellRef struct {
	Value

	Cell       Cell
	UsageCount int // TODO: garbage collecting
}

func NewCellRef(cell Cell) *CellRef {
	return &CellRef{
		Cell:       cell,
		UsageCount: 1,
	}
}

func (r *CellRef) Type(ec *Context) (int, error) {
	if ec.Visited(r.Cell) {
		return 0, NewError(ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(r.Cell)
	defer ec.ResetVisited(l)
	val, err := ec.DataProvider.Value(ec, r.Cell)
	if err != nil {
		return 0, err
	}
	return val.Type(ec)
}

func (r *CellRef) BoolValue(ec *Context) (bool, error) {
	if ec.Visited(r.Cell) {
		return false, NewError(ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(r.Cell)
	defer ec.ResetVisited(l)
	return ec.DataProvider.BoolValue(ec, r.Cell)
}

func (r *CellRef) DecimalValue(ec *Context) (decimal.Decimal, error) {
	if ec.Visited(r.Cell) {
		return decimal.Zero, NewError(ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(r.Cell)
	defer ec.ResetVisited(l)
	return ec.DataProvider.DecimalValue(ec, r.Cell)
}

func (r *CellRef) StringValue(ec *Context) (string, error) {
	if ec.Visited(r.Cell) {
		return "", NewError(ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(r.Cell)
	defer ec.ResetVisited(l)
	return ec.DataProvider.StringValue(ec, r.Cell)
}

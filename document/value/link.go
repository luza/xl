package value

import (
	"github.com/shopspring/decimal"
)

type LinkRegistryInterface interface {
	MakeLink(cellName string, sheetTitle string) (*Link, error)

	Value(ec *EvalContext, sheetIdx, x, y int) (Value, error)
	BoolValue(ec *EvalContext, sheetIdx, x, y int) (bool, error)
	DecimalValue(ec *EvalContext, sheetIdx, x, y int) (decimal.Decimal, error)
	StringValue(ec *EvalContext, sheetIdx, x, y int) (string, error)
}

type LinkCell struct {
	X int
	Y int
}

type Link struct {
	sheetIdx int
	cell     LinkCell
	cellTo   *LinkCell
	//broken   bool
	// TODO: linking context
}

func NewLink(sheetIdx int, cell LinkCell) *Link {
	return &Link{
		sheetIdx: sheetIdx,
		cell:     cell,
	}
}

func NewRangeLink(sheetIdx int, cell LinkCell, cellTo LinkCell) *Link {
	return &Link{
		sheetIdx: sheetIdx,
		cell:     cell,
		cellTo:   &cellTo,
	}
}

func (l *Link) Value(ec *EvalContext) (Value, error) {
	if l.cellTo != nil {
		return Value{}, NewError(ErrorKindCasting, "unable to use range as a single value")
	}
	if ec.Visited(l) {
		return Value{}, NewError(ErrorKindRef, "circular reference")
	}
	ec.AddVisited(l)
	return ec.LinkRegistry.Value(ec, l.sheetIdx, l.cell.X, l.cell.Y)
}

func (l *Link) BoolValue(ec *EvalContext) (bool, error) {
	if l.cellTo != nil {
		return false, NewError(ErrorKindCasting, "unable to cast range to bool")
	}
	if ec.Visited(l) {
		return false, NewError(ErrorKindRef, "circular reference")
	}
	ec.AddVisited(l)
	return ec.LinkRegistry.BoolValue(ec, l.sheetIdx, l.cell.X, l.cell.Y)
}

func (l *Link) DecimalValue(ec *EvalContext) (decimal.Decimal, error) {
	if l.cellTo != nil {
		return decimal.Zero, NewError(ErrorKindCasting, "unable to cast range to decimal")
	}
	if ec.Visited(l) {
		return decimal.Zero, NewError(ErrorKindRef, "circular reference")
	}
	ec.AddVisited(l)
	return ec.LinkRegistry.DecimalValue(ec, l.sheetIdx, l.cell.X, l.cell.Y)
}

func (l *Link) StringValue(ec *EvalContext) (string, error) {
	if l.cellTo != nil {
		return "", NewError(ErrorKindCasting, "unable to cast range to string")
	}
	if ec.Visited(l) {
		return "", NewError(ErrorKindRef, "circular reference")
	}
	ec.AddVisited(l)
	return ec.LinkRegistry.StringValue(ec, l.sheetIdx, l.cell.X, l.cell.Y)
}

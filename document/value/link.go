package value

import (
	"errors"

	"github.com/shopspring/decimal"
)

type LinkRegistryInterface interface {
	MakeLink(cellName string, sheetTitle string) (*Link, error)

	Value(sheetIdx, x, y int) (Value, error)
	BoolValue(sheetIdx, x, y int) (bool, error)
	DecimalValue(sheetIdx, x, y int) (decimal.Decimal, error)
	StringValue(sheetIdx, x, y int) (string, error)
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
	dd LinkRegistryInterface
	// TODO: linking context
}

func NewLink(sheetIdx int, cell LinkCell, dd LinkRegistryInterface) *Link {
	return &Link{
		sheetIdx: sheetIdx,
		cell:     cell,
		dd:       dd,
	}
}

func NewRangeLink(sheetIdx int, cell LinkCell, cellTo LinkCell, dd LinkRegistryInterface) *Link {
	return &Link{
		sheetIdx: sheetIdx,
		cell:     cell,
		cellTo:   &cellTo,
		dd:       dd,
	}
}

func (l *Link) Value() (Value, error) {
	if l.cellTo != nil {
		return Value{}, errors.New("unable to use range as a single value")
	}
	return l.dd.Value(l.sheetIdx, l.cell.X, l.cell.Y)
}

func (l *Link) BoolValue() (bool, error) {
	if l.cellTo != nil {
		return false, errors.New("unable to cast range to bool")
	}
	return l.dd.BoolValue(l.sheetIdx, l.cell.X, l.cell.Y)
}

func (l *Link) DecimalValue() (decimal.Decimal, error) {
	if l.cellTo != nil {
		return decimal.Zero, errors.New("unable to cast range to decimal")
	}
	return l.dd.DecimalValue(l.sheetIdx, l.cell.X, l.cell.Y)
}

func (l *Link) StringValue() (string, error) {
	if l.cellTo != nil {
		return "", errors.New("unable to cast range to string")
	}
	return l.dd.StringValue(l.sheetIdx, l.cell.X, l.cell.Y)
}

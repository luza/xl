package eval

import "github.com/shopspring/decimal"

type RangeRef struct {
	Value

	CellFromRef *CellRef
	CellToRef   *CellRef
}

func (r *RangeRef) Type(*Context) (int, error) {
	return 0, NewError(ErrorKindCasting, "unable to get type for a range")
}

func (r *RangeRef) BoolValue(ec *Context) (bool, error) {
	return false, NewError(ErrorKindCasting, "unable to cast range to bool")
}

func (r *RangeRef) DecimalValue(ec *Context) (decimal.Decimal, error) {
	return decimal.Zero, NewError(ErrorKindCasting, "unable to cast range to decimal")
}

func (r *RangeRef) StringValue(ec *Context) (string, error) {
	return "", NewError(ErrorKindCasting, "unable to cast range to string")
}

func (r *RangeRef) iterate(ec *Context, f func(Cell) error) error {
	x1, y1, x2, y2 := r.CellFromRef.Cell.X, r.CellFromRef.Cell.Y, r.CellToRef.Cell.X, r.CellToRef.Cell.Y
	if x1 > x2 || y1 > y2 {
		return NewError(ErrorKindRef, "invalid range")
	}
	for x := x1; x <= x2; x++ {
		for y := y1; y <= y2; y++ {
			cell := Cell{SheetIdx: r.CellFromRef.Cell.SheetIdx, X: x, Y: y}
			if ec.Visited(cell) {
				return NewError(ErrorKindRef, "circular reference")
			}
			l := ec.AddVisited(cell)
			err := f(cell)
			ec.ResetVisited(l)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *RangeRef) IterateBoolValues(ec *Context, f func(bool) error) error {
	return r.iterate(ec, func(cell Cell) error {
		v, err := ec.DataProvider.BoolValue(ec, cell)
		if err != nil {
			return err
		}
		return f(v)
	})
}

func (r *RangeRef) IterateDecimalValues(ec *Context, f func(decimal.Decimal) error) error {
	return r.iterate(ec, func(cell Cell) error {
		v, err := ec.DataProvider.DecimalValue(ec, cell)
		if err != nil {
			return err
		}
		return f(v)
	})
}

func (r *RangeRef) IterateStringValues(ec *Context, f func(string) error) error {
	return r.iterate(ec, func(cell Cell) error {
		v, err := ec.DataProvider.StringValue(ec, cell)
		if err != nil {
			return err
		}
		return f(v)
	})
}

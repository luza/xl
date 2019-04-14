package document

import (
	"xl/document/eval"

	"bytes"

	"github.com/shopspring/decimal"
)

func (d *Document) AddRef(cell eval.CellReference) {
	// existing link?
	for _, r := range d.refRegistry {
		if r == cell {
			return
		}
	}
	d.refRegistry = append(d.refRegistry, cell)
}

func (d *Document) FromAddress(cell eval.CellReference) (string, string, error) {
	s := d.sheetByIdx(cell.SheetIdx)
	if s == nil {
		return "", "", eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	var sheetTitle string
	if s != d.CurrentSheet {
		sheetTitle = s.Title
	}
	var buf bytes.Buffer
	if cell.AnchoredX {
		buf.WriteString("$")
	}
	buf.WriteString(ColName(cell.X))
	if cell.AnchoredY {
		buf.WriteString("$")
	}
	buf.WriteString(RowName(cell.Y))
	return sheetTitle, buf.String(), nil
}

func (d *Document) ToAddress(sheetTitle, cellName string) (eval.CellReference, error) {
	sheetIdx := d.CurrentSheet.Idx
	if sheetTitle != "" {
		found := false
		for i := range d.Sheets {
			if d.Sheets[i].Title == sheetTitle {
				sheetIdx = d.Sheets[i].Idx
				found = true
				break
			}
		}
		if !found {
			return eval.CellReference{}, eval.NewError(eval.ErrorKindRef, "sheet not found")
		}
	}
	x, y, anchoredX, anchoredY, err := CellAxis(cellName)
	if err != nil {
		return eval.CellReference{}, err
	}
	ca := eval.CellReference{
		CellAddress: eval.CellAddress{
			SheetIdx: sheetIdx,
			X:        x,
			Y:        y,
		},
		AnchoredX: anchoredX,
		AnchoredY: anchoredY,
	}
	return ca, nil
}

// TODO: do we really need this method?
func (d *Document) Value(ec *eval.Context, cell eval.CellAddress) (eval.Value, error) {
	s := d.sheetByIdx(cell.SheetIdx)
	if s == nil {
		return eval.NewEmptyValue(), eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(cell.X, cell.Y)
	if c == nil {
		return eval.NewEmptyValue(), nil
	}
	if ec.Visited(cell) {
		return eval.NewEmptyValue(), eval.NewError(eval.ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(cell)
	defer ec.ResetVisited(l)
	return c.Value(ec)
}

func (d *Document) BoolValue(ec *eval.Context, cell eval.CellAddress) (bool, error) {
	s := d.sheetByIdx(cell.SheetIdx)
	if s == nil {
		return false, eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(cell.X, cell.Y)
	if c == nil {
		return false, nil
	}
	if ec.Visited(cell) {
		return false, eval.NewError(eval.ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(cell)
	defer ec.ResetVisited(l)
	return c.BoolValue(ec)
}

func (d *Document) DecimalValue(ec *eval.Context, cell eval.CellAddress) (decimal.Decimal, error) {
	s := d.sheetByIdx(cell.SheetIdx)
	if s == nil {
		return decimal.Zero, eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(cell.X, cell.Y)
	if c == nil {
		return decimal.Zero, nil
	}
	if ec.Visited(cell) {
		return decimal.Zero, eval.NewError(eval.ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(cell)
	defer ec.ResetVisited(l)
	return c.DecimalValue(ec)
}

func (d *Document) StringValue(ec *eval.Context, cell eval.CellAddress) (string, error) {
	s := d.sheetByIdx(cell.SheetIdx)
	if s == nil {
		return "", eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(cell.X, cell.Y)
	if c == nil {
		return "", nil
	}
	if ec.Visited(cell) {
		return "", eval.NewError(eval.ErrorKindRef, "circular reference")
	}
	l := ec.AddVisited(cell)
	defer ec.ResetVisited(l)
	return c.StringValue(ec)
}

func (d *Document) iterate(ec *eval.Context, cell, cellTo eval.CellAddress, f func(eval.CellAddress) error) error {
	if cell.SheetIdx != cellTo.SheetIdx {
		// cross-sheets ranges are not allowed
		return eval.NewError(eval.ErrorKindRef, "cross-sheets ranges are not allowed")
	}
	if cell.X > cellTo.X || cell.Y > cellTo.Y {
		return eval.NewError(eval.ErrorKindRef, "invalid range bounds")
	}
	for x := cell.X; x <= cellTo.X; x++ {
		for y := cell.Y; y <= cellTo.Y; y++ {
			err := f(eval.CellAddress{SheetIdx: cell.SheetIdx, X: x, Y: y})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *Document) IterateBoolValues(ec *eval.Context, cell, cellTo eval.CellAddress, f func(bool) error) error {
	return d.iterate(ec, cell, cellTo, func(cell eval.CellAddress) error {
		v, err := ec.DataProvider.BoolValue(ec, cell)
		if err != nil {
			return err
		}
		return f(v)
	})
}

func (d *Document) IterateDecimalValues(ec *eval.Context, cell, cellTo eval.CellAddress, f func(decimal.Decimal) error) error {
	return d.iterate(ec, cell, cellTo, func(cell eval.CellAddress) error {
		v, err := ec.DataProvider.DecimalValue(ec, cell)
		if err != nil {
			return err
		}
		return f(v)
	})
}

func (d *Document) IterateStringValues(ec *eval.Context, cell, cellTo eval.CellAddress, f func(string) error) error {
	return d.iterate(ec, cell, cellTo, func(cell eval.CellAddress) error {
		v, err := ec.DataProvider.StringValue(ec, cell)
		if err != nil {
			return err
		}
		return f(v)
	})
}

func (d *Document) moveRefsRight(n int) {
	for _, r := range d.refRegistry {
		if r.SheetIdx == d.CurrentSheet.Idx && r.X >= n {
			r.X++
		}
	}
}

func (d *Document) moveRefsLeft(n int) {
	for _, r := range d.refRegistry {
		if r.SheetIdx == d.CurrentSheet.Idx && r.X > n {
			r.X--
		}
	}
}

func (d *Document) moveRefsDown(n int) {
	for _, r := range d.refRegistry {
		if r.SheetIdx == d.CurrentSheet.Idx && r.Y >= n {
			r.Y++
		}
	}
}

func (d *Document) moveRefsUp(n int) {
	for _, r := range d.refRegistry {
		if r.SheetIdx == d.CurrentSheet.Idx && r.Y > n {
			r.Y--
		}
	}
}

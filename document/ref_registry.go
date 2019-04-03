package document

import (
	"xl/document/eval"
	"xl/document/sheet"

	"github.com/shopspring/decimal"
)

func (d *Document) NewCellRef(sheetTitle, cellName string) (*eval.CellRef, error) {
	var s *sheet.Sheet
	if sheetTitle != "" {
		for i := range d.Sheets {
			if d.Sheets[i].Title == sheetTitle {
				s = d.Sheets[i]
				break
			}
		}
		// sheet not found
		if s == nil {
			return nil, eval.NewError(eval.ErrorKindName, "sheet does not exist")
		}
	} else {
		s = d.CurrentSheet
	}
	x, y, err := CellAxis(cellName)
	if err != nil {
		return nil, err
	}
	// existing link?
	for _, r := range d.refRegistry {
		if r.Cell.SheetIdx == s.Idx && r.Cell.X == x && r.Cell.Y == y {
			r.UsageCount++
			return r, nil
		}
	}
	// not found? create new one
	r := eval.NewCellRef(eval.Cell{SheetIdx: s.Idx, X: x, Y: y})
	d.refRegistry = append(d.refRegistry, r)
	return r, nil
}

func (d *Document) NewRangeRef(sheetFromTitle, cellFromName, sheetToTitle, cellToName string) (*eval.RangeRef, error) {
	fromRef, err := d.NewCellRef(sheetFromTitle, cellFromName)
	if err != nil {
		return nil, err
	}
	toRef, err := d.NewCellRef(sheetToTitle, cellToName)
	if err != nil {
		return nil, err
	}
	rr := &eval.RangeRef{
		CellFromRef: fromRef,
		CellToRef:   toRef,
	}
	return rr, nil
}

func (d *Document) SheetTitle(sheetIdx int) (string, error) {
	if s := d.sheetByIdx(sheetIdx); s != nil {
		return s.Title, nil
	}
	return "", eval.NewError(eval.ErrorKindRef, "sheet does not exist")
}

func (d *Document) CellName(cell eval.Cell) (string, error) {
	//  FIXME: accept sheet name?
	return CellName(cell.X, cell.Y), nil
}

func (d *Document) Value(ec *eval.Context, cell eval.Cell) (eval.Value, error) {
	s := d.sheetByIdx(cell.SheetIdx)
	if s == nil {
		return eval.NewEmptyValue(), eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(cell.X, cell.Y)
	if c == nil {
		return eval.NewEmptyValue(), nil
	}
	return c.Value(ec)
}

func (d *Document) BoolValue(ec *eval.Context, cell eval.Cell) (bool, error) {
	s := d.sheetByIdx(cell.SheetIdx)
	if s == nil {
		return false, eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(cell.X, cell.Y)
	if c == nil {
		return false, nil
	}
	return c.BoolValue(ec)
}

func (d *Document) DecimalValue(ec *eval.Context, cell eval.Cell) (decimal.Decimal, error) {
	s := d.sheetByIdx(cell.SheetIdx)
	if s == nil {
		return decimal.Zero, eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(cell.X, cell.Y)
	if c == nil {
		return decimal.Zero, nil
	}
	return c.DecimalValue(ec)
}

func (d *Document) StringValue(ec *eval.Context, cell eval.Cell) (string, error) {
	s := d.sheetByIdx(cell.SheetIdx)
	if s == nil {
		return "", eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(cell.X, cell.Y)
	if c == nil {
		return "", nil
	}
	return c.StringValue(ec)
}

func (d *Document) moveRefsRight(n int) {
	for _, r := range d.refRegistry {
		if r.Cell.SheetIdx == d.CurrentSheet.Idx && r.Cell.X >= n {
			r.Cell.X++
		}
	}
}

func (d *Document) moveRefsLeft(n int) {
	for _, r := range d.refRegistry {
		if r.Cell.SheetIdx == d.CurrentSheet.Idx && r.Cell.X > n {
			r.Cell.X--
		}
	}
}

func (d *Document) moveRefsDown(n int) {
	for _, r := range d.refRegistry {
		if r.Cell.SheetIdx == d.CurrentSheet.Idx && r.Cell.Y >= n {
			r.Cell.Y++
		}
	}
}

func (d *Document) moveRefsUp(n int) {
	for _, r := range d.refRegistry {
		if r.Cell.SheetIdx == d.CurrentSheet.Idx && r.Cell.Y > n {
			r.Cell.Y--
		}
	}
}

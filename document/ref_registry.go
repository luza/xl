package document

import (
	"fmt"

	"xl/document/eval"
	"xl/document/sheet"
	"xl/log"

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
	if r, ok := d.refRegistry[s.Idx][x][y]; ok {
		log.L.Error(fmt.Sprintf("reused link sheet %d x %d y %d\n", s.Idx, x, y))
		return r, nil
	}
	// not found? create new one
	r := eval.NewCellRef(s.Idx, eval.Axis{X: x, Y: y})
	if _, ok := d.refRegistry[s.Idx]; !ok {
		d.refRegistry[s.Idx] = make(map[int]map[int]*eval.CellRef)
	}
	if _, ok := d.refRegistry[s.Idx][x]; !ok {
		d.refRegistry[s.Idx][x] = make(map[int]*eval.CellRef)
	}
	d.refRegistry[s.Idx][x][y] = r
	return r, nil
}

func (d *Document) SheetTitle(r *eval.CellRef) (string, error) {
	s := d.sheetByIdx(r.SheetIdx)
	if s != nil {
		return "", eval.NewError(eval.ErrorKindRef, "sheet does not exist")
	}
	return s.Title, nil
}

func (d *Document) CellName(r *eval.CellRef) (string, error) {
	return CellName(r.Cell.X, r.Cell.Y), nil
}

func (d *Document) Value(ec *eval.Context, r *eval.CellRef) (eval.Value, error) {
	s := d.sheetByIdx(r.SheetIdx)
	if s == nil {
		return eval.NullValue(), eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(r.Cell.X, r.Cell.Y)
	if c == nil {
		return eval.NullValue(), nil
	}
	return c.Value(ec)
}

func (d *Document) BoolValue(ec *eval.Context, r *eval.CellRef) (bool, error) {
	s := d.sheetByIdx(r.SheetIdx)
	if s == nil {
		return false, eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(r.Cell.X, r.Cell.Y)
	if c == nil {
		return false, nil
	}
	return c.BoolValue(ec)
}

func (d *Document) DecimalValue(ec *eval.Context, r *eval.CellRef) (decimal.Decimal, error) {
	s := d.sheetByIdx(r.SheetIdx)
	if s == nil {
		return decimal.Zero, eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(r.Cell.X, r.Cell.Y)
	if c == nil {
		return decimal.Zero, nil
	}
	return c.DecimalValue(ec)
}

func (d *Document) StringValue(ec *eval.Context, r *eval.CellRef) (string, error) {
	s := d.sheetByIdx(r.SheetIdx)
	if s == nil {
		return "", eval.NewError(eval.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(r.Cell.X, r.Cell.Y)
	if c == nil {
		return "", nil
	}
	return c.StringValue(ec)
}

package document

import (
	"fmt"

	"xl/document/sheet"
	"xl/document/value"
	"xl/log"

	"github.com/shopspring/decimal"
)

func (d *Document) MakeLink(cellName string, sheetTitle string) (*value.Link, error) {
	log.L.Error(fmt.Sprintf("linking cell %s sheet %v\n", cellName, sheetTitle))
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
			return nil, value.NewError(value.ErrorKindName, "sheet does not exist")
		}
	} else {
		s = d.CurrentSheet
	}
	x, y, err := cellNameToXY(cellName)
	if err != nil {
		return nil, err
	}
	// existing link?
	if l, ok := d.linksRegistry[s.Idx][x][y]; ok {
		log.L.Error(fmt.Sprintf("reused link sheet %d x %d y %d\n", s.Idx, x, y))
		return l, nil
	}
	// not found? create new one
	l := value.NewLink(s.Idx, value.LinkCell{X: x, Y: y})
	log.L.Error(fmt.Sprintf("created link sheet %d x %d y %d\n", s.Idx, x, y))
	if _, ok := d.linksRegistry[s.Idx]; !ok {
		d.linksRegistry[s.Idx] = make(map[int]map[int]*value.Link)
	}
	if _, ok := d.linksRegistry[s.Idx][x]; !ok {
		d.linksRegistry[s.Idx][x] = make(map[int]*value.Link)
	}
	d.linksRegistry[s.Idx][x][y] = l
	return l, nil
}

func (d *Document) Value(ec *value.EvalContext, sheetIdx, x, y int) (value.Value, error) {
	s := d.sheetByIdx(sheetIdx)
	if s == nil {
		return value.Value{}, value.NewError(value.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(x, y)
	if c == nil {
		return value.Value{}, nil
	}
	return c.Value(ec)
}

func (d *Document) BoolValue(ec *value.EvalContext, sheetIdx, x, y int) (bool, error) {
	s := d.sheetByIdx(sheetIdx)
	if s == nil {
		return false, value.NewError(value.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(x, y)
	if c == nil {
		return false, nil
	}
	return c.BoolValue(ec)
}

func (d *Document) DecimalValue(ec *value.EvalContext, sheetIdx, x, y int) (decimal.Decimal, error) {
	s := d.sheetByIdx(sheetIdx)
	if s == nil {
		return decimal.Zero, value.NewError(value.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(x, y)
	if c == nil {
		return decimal.Zero, nil
	}
	return c.DecimalValue(ec)
}

func (d *Document) StringValue(ec *value.EvalContext, sheetIdx, x, y int) (string, error) {
	s := d.sheetByIdx(sheetIdx)
	if s == nil {
		return "", value.NewError(value.ErrorKindName, "sheet does not exist")
	}
	c := s.Cell(x, y)
	if c == nil {
		return "", nil
	}
	return c.StringValue(ec)
}

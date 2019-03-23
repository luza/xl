package formula

import (
	"strings"
)

type VarCell struct {
	Sheet string
	Cell  string
}

// FIXME: is this possible to get rid of this struct
type Var struct {
	Cell   VarCell
	CellTo *VarCell
}

func newVar(c *CellRange) Var {
	var s string
	if c.Cell.Sheet != nil {
		s = string(*c.Cell.Sheet)
	}
	v := Var{
		Cell: VarCell{
			Sheet: s,
			Cell:  strings.ToUpper(c.Cell.Cell),
		},
	}
	if c.CellTo != nil {
		var s string
		if c.CellTo.Sheet != nil {
			s = string(*c.CellTo.Sheet)
		}
		v.CellTo = &VarCell{
			Sheet: s,
			Cell:  strings.ToUpper(c.CellTo.Cell),
		}
	}
	return v
}

type VarBin struct {
	Vars []Var
}

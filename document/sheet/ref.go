package sheet

import (
	"xl/document/eval"
	"xl/formula"
)

type ref struct {
	Cell   eval.CellReference
	CellTo *eval.CellReference
}

func toAddress(ec *eval.Context, c *formula.Cell) (eval.CellReference, error) {
	var sheetTitle string
	if c.Sheet != nil {
		sheetTitle = string(*c.Sheet)
	}
	return ec.DataProvider.ToAddress(sheetTitle, c.CellName)
}

func fromAddress(ec *eval.Context, ca eval.CellReference) (*formula.Cell, error) {
	sheetTitle, cellName, err := ec.DataProvider.FromAddress(ca)
	if err != nil {
		return nil, err
	}
	var s *formula.Sheet
	if sheetTitle != "" {
		fs := formula.Sheet(sheetTitle)
		s = &fs
	}
	c := &formula.Cell{
		Sheet:    s,
		CellName: cellName,
	}
	return c, nil
}

func makeRefs(ec *eval.Context, vars []*formula.Variable) ([]ref, error) {
	refs := make([]ref, len(vars))
	for i, v := range vars {
		ca, err := toAddress(ec, v.Cell)
		if err != nil {
			return nil, err
		}
		ec.DataProvider.AddRef(ca)
		ref := ref{
			Cell: ca,
		}
		if v.CellTo != nil {
			ca, err = toAddress(ec, v.CellTo)
			if err != nil {
				return nil, err
			}
			ec.DataProvider.AddRef(ca)
			ref.CellTo = &ca
		}
		refs[i] = ref
	}
	return refs, nil
}

func refsToValues(refs []ref) []eval.Value {
	values := make([]eval.Value, len(refs))
	for i, r := range refs {
		values[i] = eval.NewRefValue(r.Cell, r.CellTo)
	}
	return values
}

func updateVars(ec *eval.Context, x *formula.Expression, refs []ref) error {
	for i, v := range x.Variables() {
		c, err := fromAddress(ec, refs[i].Cell)
		if err != nil {
			return err
		}
		v.Cell = c
		if refs[i].CellTo != nil {
			c, err = fromAddress(ec, *refs[i].CellTo)
			if err != nil {
				return err
			}
			v.CellTo = c
		}
	}
	return nil
}

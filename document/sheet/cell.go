package sheet

import (
	"strconv"

	"xl/document/eval"
	"xl/formula"

	"github.com/shopspring/decimal"
)

const (
	CellValueUntyped = iota
	CellValueTypeEmpty
	CellValueTypeText
	CellValueTypeInteger
	CellValueTypeDecimal
	CellValueTypeBool
	CellValueTypeFormula
)

type Cell struct {
	valueType int

	// values union
	rawValue     string
	intValue     int
	decimalValue *decimal.Decimal
	boolValue    bool
	formulaValue formula.Function

	// formula params
	expression *formula.Expression
	refs       []eval.Value
}

func NewCellEmpty() *Cell {
	return &Cell{
		valueType: CellValueTypeEmpty,
	}
}

func NewCellUntyped(v string) *Cell {
	c := &Cell{}
	c.SetValueUntyped(v)
	return c
}

func (c *Cell) Free() {
	for _, r := range c.refs {
		switch r := r.(type) {
		case *eval.CellRef:
			r.UsageCount--
		}
	}
}

// EraseValue resets cell value to initial.
func (c *Cell) EraseValue() {
	c.rawValue = ""
	c.boolValue = false
	c.intValue = 0
	c.decimalValue = nil
	c.formulaValue = nil
	c.valueType = CellValueTypeEmpty
}

// RawValue returns raw cell value as string. No evaluation performed.
func (c *Cell) RawValue() string {
	return c.rawValue
}

func (c *Cell) Expression(ec *eval.Context) *formula.Expression {
	if c.valueType == CellValueUntyped {
		if err := c.evaluateType(ec); err != nil {
			return nil
		}
	}
	if c.valueType != CellValueTypeFormula {
		return nil
	}
	if err := updateVars(ec, c.expression, c.refs); err != nil {
		return nil
	}
	return c.expression
}

func (c *Cell) BoolValue(ec *eval.Context) (bool, error) {
	if c.valueType == CellValueUntyped {
		if err := c.evaluateType(ec); err != nil {
			return false, err
		}
	}
	switch c.valueType {
	case CellValueTypeEmpty:
		return false, nil
	case CellValueTypeText:
		return false, eval.NewError(eval.ErrorKindCasting, "unable to cast text to bool")
	case CellValueTypeInteger:
		return c.intValue != 0, nil
	case CellValueTypeDecimal:
		return !c.decimalValue.Equal(decimal.Zero), nil
	case CellValueTypeBool:
		return c.boolValue, nil
	case CellValueTypeFormula:
		val, err := c.formulaValue(ec, c.refs)
		if err != nil {
			return false, err
		}
		return val.BoolValue(ec)
	}
	return false, nil
}

// DecimalValue returns evaluated cell value as decimal.
func (c *Cell) DecimalValue(ec *eval.Context) (decimal.Decimal, error) {
	if c.valueType == CellValueUntyped {
		if err := c.evaluateType(ec); err != nil {
			return decimal.Zero, err
		}
	}
	switch c.valueType {
	case CellValueTypeEmpty:
		return decimal.Zero, nil
	case CellValueTypeText:
		return decimal.Zero, eval.NewError(eval.ErrorKindCasting, "unable to cast text to decimal")
	case CellValueTypeInteger:
		return decimal.New(int64(c.intValue), 0), nil
	case CellValueTypeDecimal:
		return *c.decimalValue, nil
	case CellValueTypeBool:
		return decimal.Zero, eval.NewError(eval.ErrorKindCasting, "unable to cast bool to decimal")
	case CellValueTypeFormula:
		val, err := c.formulaValue(ec, c.refs)
		if err != nil {
			return decimal.Zero, err
		}
		return val.DecimalValue(ec)
	}
	return decimal.Zero, nil
}

// StringValue returns evaluated cell rawValue as string.
func (c *Cell) StringValue(ec *eval.Context) (string, error) {
	if c.valueType == CellValueUntyped {
		if err := c.evaluateType(ec); err != nil {
			return "", err
		}
	}
	if c.valueType == CellValueTypeFormula {
		val, err := c.formulaValue(ec, c.refs)
		if err != nil {
			return "", err
		}
		return val.StringValue(ec)
	}
	return c.rawValue, nil
}

func (c *Cell) Value(ec *eval.Context) (eval.Value, error) {
	if c.valueType == CellValueUntyped {
		if err := c.evaluateType(ec); err != nil {
			return eval.NullValue(), err
		}
	}
	switch c.valueType {
	case CellValueTypeEmpty:
		return eval.NewStringValue(""), nil
	case CellValueTypeText:
		return eval.NewStringValue(c.rawValue), nil
	case CellValueTypeInteger:
		return eval.NewDecimalValue(decimal.New(int64(c.intValue), 0)), nil
	case CellValueTypeDecimal:
		return eval.NewDecimalValue(*c.decimalValue), nil
	case CellValueTypeBool:
		return eval.NewBoolValue(c.boolValue), nil
	case CellValueTypeFormula:
		return c.formulaValue(ec, c.refs)
	}
	panic("unsupported type")
}

// SetValueUntyped fill new cell value with no any type associated with it.
// Type will be determined later on demand.
func (c *Cell) SetValueUntyped(v string) {
	c.EraseValue()
	c.valueType = CellValueUntyped
	c.rawValue = v
}

func (c *Cell) evaluateType(ec *eval.Context) error {
	t, castedV := guessCellType(c.rawValue)
	switch t {
	case CellValueTypeInteger:
		c.intValue = castedV.(int)
	case CellValueTypeDecimal:
		d, _ := decimal.NewFromString(c.rawValue)
		c.decimalValue = &d
	case CellValueTypeBool:
		c.boolValue = castedV.(bool)
	case CellValueTypeFormula:
		c.formulaValue = nil
		c.refs = nil
		expr, err := formula.Parse(c.rawValue)
		if err != nil {
			return err
		}
		c.formulaValue, _ = expr.BuildFunc()
		c.expression = expr
		c.rawValue = expr.String() // need this?
		c.refs, err = makeRefs(expr.Variables(), ec)
		if err != nil {
			return err
		}
	}
	c.valueType = t
	return nil
}

// guessCellType detects cell type based on its rawValue.
// Returns detected type and either casted rawValue or nil if casting wasn't done.
func guessCellType(v string) (int, interface{}) {
	if len(v) == 0 {
		return CellValueTypeEmpty, nil
	} else if v[0] == '=' && len(v) > 1 {
		return CellValueTypeFormula, nil
	} else {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return CellValueTypeInteger, int(i)
		}
		if _, err := strconv.ParseFloat(v, 64); err == nil {
			return CellValueTypeDecimal, nil
		}
		if b, err := strconv.ParseBool(v); err == nil {
			return CellValueTypeBool, b
		}
	}
	return CellValueTypeText, v
}

func makeRefs(vars []*formula.Variable, ec *eval.Context) ([]eval.Value, error) {
	values := make([]eval.Value, len(vars))
	for i := range vars {
		if vars[i].CellTo != nil {
			// range
			//links[i] = dd.LinkRange(c.cell, c.CellTo, c.Sheet)
			//values[i] = eval.NewLinkValue(l)
		} else {
			c := vars[i].Cell
			var s string
			if c.Sheet != nil {
				s = string(*c.Sheet)
			}
			ref, err := ec.DataProvider.NewCellRef(s, c.Cell)
			if err != nil {
				return nil, err
			}
			values[i] = ref
		}
	}
	return values, nil
}

func updateVars(ec *eval.Context, x *formula.Expression, refs []eval.Value) error {
	for i, v := range x.Variables() {
		switch r := refs[i].(type) {
		case *eval.CellRef:
			sheetTitle, err := ec.DataProvider.SheetTitle(r)
			if err != nil {
				return err
			}
			s := formula.Sheet(sheetTitle)
			v.Cell.Sheet = &s
			cellName, err := ec.DataProvider.CellName(r)
			if err != nil {
				return err
			}
			v.Cell.Cell = cellName
		default:
			panic("unexpected value type")
		}
	}
	return nil
}

package sheet

import (
	"strconv"

	"xl/document/value"
	"xl/formula"
	"xl/log"

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

	// formula arguments
	args []value.Value
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

func (c *Cell) BoolValue(ec *value.EvalContext) (bool, error) {
	if c.valueType == CellValueUntyped {
		if err := c.evaluateType(ec); err != nil {
			return false, err
		}
	}
	switch c.valueType {
	case CellValueTypeEmpty:
		return false, nil
	case CellValueTypeText:
		return false, value.NewError(value.ErrorKindCasting, "unable to cast text to bool")
	case CellValueTypeInteger:
		return c.intValue != 0, nil
	case CellValueTypeDecimal:
		return !c.decimalValue.Equal(decimal.Zero), nil
	case CellValueTypeBool:
		return c.boolValue, nil
	case CellValueTypeFormula:
		val, err := c.formulaValue(ec, c.args)
		if err != nil {
			return false, err
		}
		return val.BoolValue(ec)
	}
	return false, nil
}

// DecimalValue returns evaluated cell value as decimal.
func (c *Cell) DecimalValue(ec *value.EvalContext) (decimal.Decimal, error) {
	if c.valueType == CellValueUntyped {
		if err := c.evaluateType(ec); err != nil {
			return decimal.Zero, err
		}
	}
	switch c.valueType {
	case CellValueTypeEmpty:
		return decimal.Zero, nil
	case CellValueTypeText:
		return decimal.Zero, value.NewError(value.ErrorKindCasting, "unable to cast text to decimal")
	case CellValueTypeInteger:
		return decimal.New(int64(c.intValue), 0), nil
	case CellValueTypeDecimal:
		return *c.decimalValue, nil
	case CellValueTypeBool:
		return decimal.Zero, value.NewError(value.ErrorKindCasting, "unable to cast bool to decimal")
	case CellValueTypeFormula:
		val, err := c.formulaValue(ec, c.args)
		if err != nil {
			return decimal.Zero, err
		}
		return val.DecimalValue(ec)
	}
	return decimal.Zero, nil
}

// StringValue returns evaluated cell rawValue as string.
func (c *Cell) StringValue(ec *value.EvalContext) (string, error) {
	if c.valueType == CellValueUntyped {
		if err := c.evaluateType(ec); err != nil {
			return "", err
		}
	}
	if c.valueType == CellValueTypeFormula {
		val, err := c.formulaValue(ec, c.args)
		if err != nil {
			return "", err
		}
		return val.StringValue(ec)
	}
	return c.rawValue, nil
}

func (c *Cell) Value(ec *value.EvalContext) (value.Value, error) {
	if c.valueType == CellValueUntyped {
		if err := c.evaluateType(ec); err != nil {
			return value.Value{}, err
		}
	}
	switch c.valueType {
	case CellValueTypeEmpty:
		return value.NewStringValue(""), nil
	case CellValueTypeText:
		return value.NewStringValue(c.rawValue), nil
	case CellValueTypeInteger:
		return value.NewDecimalValue(decimal.New(int64(c.intValue), 0)), nil
	case CellValueTypeDecimal:
		return value.NewDecimalValue(*c.decimalValue), nil
	case CellValueTypeBool:
		return value.NewBoolValue(c.boolValue), nil
	case CellValueTypeFormula:
		return c.formulaValue(ec, c.args)
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

func (c *Cell) evaluateType(ec *value.EvalContext) error {
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
		c.args = nil
		formulaValue, vars, err := formula.Parse(c.rawValue)
		if err != nil {
			return err
		}
		c.formulaValue = formulaValue
		c.args, err = makeLinks(vars, ec)
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
		if b, err := strconv.ParseBool(v); err == nil {
			return CellValueTypeBool, b
		}
		if _, err := strconv.ParseFloat(v, 64); err == nil {
			return CellValueTypeDecimal, nil
		}
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return CellValueTypeInteger, i
		}
	}
	return CellValueTypeText, v
}

func makeLinks(vb *formula.VarBin, ec *value.EvalContext) ([]value.Value, error) {
	values := make([]value.Value, len(vb.Vars))
	for i := range vb.Vars {
		log.L.Error("converting var to link")
		if vb.Vars[i].CellTo != nil {
			// range
			//links[i] = dd.LinkRange(c.Cell, c.CellTo, c.Sheet)
			//values[i] = value.NewLinkValue(l)
		} else {
			c := vb.Vars[i].Cell
			l, err := ec.LinkRegistry.MakeLink(c.Cell, c.Sheet)
			if err != nil {
				return nil, err
			}
			values[i] = value.NewLinkValue(l)
		}
	}
	return values, nil
}

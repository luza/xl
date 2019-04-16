package eval

import (
	"github.com/shopspring/decimal"
)

const (
	TypeEmpty = iota
	TypeBool
	TypeDecimal
	TypeString
	TypeRef
	TypeRangeRef
)

// Значение - это единица информация, над которой производятся вычисления в формулах.
// Значение может быть пустым, быть константным заначением одного из трех типов, или хранить в себе ссылку
// на ячейку или диапазон ячеек.

type Value interface {
	Type() int
	BoolValue(*Context) (bool, error)
	DecimalValue(*Context) (decimal.Decimal, error)
	StringValue(*Context) (string, error)
	Cell() CellReference
	CellTo() CellReference
}

type staticValue struct {
	Value

	valueType    int
	boolValue    bool
	decimalValue decimal.Decimal
	stringValue  string

	// ref
	cell   *CellReference
	cellTo *CellReference
}

func NewEmptyValue() Value {
	return staticValue{
		valueType: TypeEmpty,
	}
}

func NewBoolValue(v bool) Value {
	return staticValue{
		valueType: TypeBool,
		boolValue: v,
	}
}

func NewDecimalValue(v decimal.Decimal) Value {
	return staticValue{
		valueType:    TypeDecimal,
		decimalValue: v,
	}
}

func NewStringValue(v string) Value {
	return staticValue{
		valueType:   TypeString,
		stringValue: v,
	}
}

// TODO(low): accept address instead of reference?
func NewRefValue(cell CellReference, cellTo *CellReference) Value {
	t := TypeRef
	if cellTo != nil {
		t = TypeRangeRef
	}
	return staticValue{
		valueType: t,
		cell:      &cell,
		cellTo:    cellTo,
	}
}

func (v staticValue) Type() int {
	return v.valueType
}

// Возвращает щначение, приведенное к булевому типу.
func (v staticValue) BoolValue(ec *Context) (bool, error) {
	switch v.valueType {
	case TypeEmpty:
		return false, nil
	case TypeBool:
		return v.boolValue, nil
	case TypeDecimal:
		return !v.decimalValue.Equal(decimal.Zero), nil
	case TypeString:
		if len(v.stringValue) == 0 {
			return false, nil
		}
		return false, NewError(ErrorKindCasting, "unable to cast string value %s to bool", v.stringValue)
	case TypeRef:
		return ec.DataProvider.BoolValue(ec, v.cell.CellAddress)
	case TypeRangeRef:
		return false, NewError(ErrorKindCasting, "unable to use range as bool value")
	default:
		panic("invalid type")
	}
}

// Возвращает значение, приведенное к числовому типу.
func (v staticValue) DecimalValue(ec *Context) (decimal.Decimal, error) {
	switch v.valueType {
	case TypeEmpty:
		return decimal.Zero, nil
	case TypeBool:
		if v.boolValue {
			return decimal.NewFromFloat(1.0), nil
		} else {
			return decimal.Zero, nil
		}
	case TypeDecimal:
		return v.decimalValue, nil
	case TypeString:
		if len(v.stringValue) == 0 {
			return decimal.Zero, nil
		}
		return decimal.Zero, NewError(ErrorKindCasting, "unable to cast string value %s to decimal", v.stringValue)
	case TypeRef:
		return ec.DataProvider.DecimalValue(ec, v.cell.CellAddress)
	case TypeRangeRef:
		return decimal.Zero, NewError(ErrorKindCasting, "unable to use range as decimal value")
	default:
		panic("invalid type")
	}
}

// Возвращает значение, приведенное к строке.
func (v staticValue) StringValue(ec *Context) (string, error) {
	switch v.valueType {
	case TypeEmpty:
		return "", nil
	case TypeBool:
		if v.boolValue {
			return "TRUE", nil
		} else {
			return "FALSE", nil
		}
	case TypeDecimal:
		return v.decimalValue.String(), nil
	case TypeString:
		return v.stringValue, nil
	case TypeRef:
		return ec.DataProvider.StringValue(ec, v.cell.CellAddress)
	case TypeRangeRef:
		return "", NewError(ErrorKindCasting, "unable to use range as string value")
	default:
		panic("invalid type")
	}
}

func (v staticValue) Cell() CellReference {
	if v.valueType != TypeRef && v.valueType != TypeRangeRef {
		panic("type is not TypeRef")
	}
	return *v.cell
}

func (v staticValue) CellTo() CellReference {
	if v.valueType != TypeRangeRef {
		panic("type is not TypeRangeRef")
	}
	return *v.cellTo
}

package eval

import "github.com/shopspring/decimal"

type staticValue struct {
	Value

	valueType    int
	boolValue    bool
	decimalValue decimal.Decimal
	stringValue  string
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

func (v staticValue) Type(*Context) (int, error) {
	return v.valueType, nil
}

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
	default:
		panic("invalid type")
	}
}

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
	default:
		panic("invalid type")
	}
}

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
	default:
		panic("invalid type")
	}
}

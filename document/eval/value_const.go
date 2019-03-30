package eval

import "github.com/shopspring/decimal"

type staticValue struct {
	Value

	valueType    int
	boolValue    bool
	decimalValue decimal.Decimal
	stringValue  string
}

func NullValue() Value {
	return staticValue{}
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
	case TypeBool:
		return v.boolValue, nil
	case TypeDecimal:
		return !v.decimalValue.Equal(decimal.Zero), nil
	case TypeString:
		return false, NewError(ErrorKindCasting, "unable to cast string value %s to bool", v.stringValue)
	default:
		panic("invalid type")
	}
}

func (v staticValue) DecimalValue(ec *Context) (decimal.Decimal, error) {
	switch v.valueType {
	case TypeBool:
		if v.boolValue {
			return decimal.NewFromFloat(1.0), nil
		} else {
			return decimal.Zero, nil
		}
	case TypeDecimal:
		return v.decimalValue, nil
	case TypeString:
		return decimal.Zero, NewError(ErrorKindCasting, "unable to cast string value %s to decimal", v.stringValue)
	default:
		panic("invalid type")
	}
}

func (v staticValue) StringValue(ec *Context) (string, error) {
	switch v.valueType {
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

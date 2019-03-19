package value

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

const (
	TypeBool = iota
	TypeDecimal
	TypeString
	TypeCell
	TypeRange
)

type Value struct {
	Type int

	boolValue    bool
	decimalValue decimal.Decimal
	stringValue  string

	link *Link
}

func NewBoolValue(v bool) Value {
	return Value{
		Type:      TypeBool,
		boolValue: v,
	}
}

func NewDecimalValue(v decimal.Decimal) Value {
	return Value{
		Type:         TypeDecimal,
		decimalValue: v,
	}
}

func NewStringValue(v string) Value {
	return Value{
		Type:        TypeString,
		stringValue: v,
	}
}

func NewLinkValue(v *Link) Value {
	// TODO: range?
	return Value{
		Type: TypeCell,
		link: v,
	}
}

func (v Value) BoolValue() (bool, error) {
	switch v.Type {
	case TypeBool:
		return v.boolValue, nil
	case TypeDecimal:
		return !v.decimalValue.Equal(decimal.Zero), nil
	case TypeString:
		return false, fmt.Errorf("unable to cast string value %s to bool", v.stringValue)
	case TypeCell:
		return v.link.BoolValue()
	case TypeRange:
		return false, fmt.Errorf("unable to cast range to bool")
	default:
		panic("invalid type")
	}
}

func (v Value) DecimalValue() (decimal.Decimal, error) {
	switch v.Type {
	case TypeBool:
		if v.boolValue {
			return decimal.NewFromFloat(1.0), nil
		} else {
			return decimal.Zero, nil
		}
	case TypeDecimal:
		return v.decimalValue, nil
	case TypeString:
		return decimal.Zero, fmt.Errorf("unable to cast string value %s to decimal", v.stringValue)
	case TypeCell:
		return v.link.DecimalValue()
	case TypeRange:
		return decimal.Zero, fmt.Errorf("unable to cast range to decimal")
	default:
		panic("invalid type")
	}
}

func (v Value) StringValue() (string, error) {
	switch v.Type {
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
	case TypeCell:
		return v.link.StringValue()
	case TypeRange:
		return "", fmt.Errorf("unable to cast range to string")
	default:
		panic("invalid type")
	}
}

func (v Value) Link() (*Link, error) {
	switch v.Type {
	case TypeBool, TypeDecimal, TypeString:
		return nil, errors.New("unable to use static value as a link")
	case TypeCell:
		return v.link, nil
	//case TypeRange:
	//	return v.linkedRange, nil
	default:
		panic("invalid type")
	}
}

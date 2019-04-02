package eval

import (
	"github.com/shopspring/decimal"
)

const (
	TypeEmpty = iota
	TypeBool
	TypeDecimal
	TypeString
)

type Value interface {
	Type(*Context) (int, error)
	BoolValue(*Context) (bool, error)
	DecimalValue(*Context) (decimal.Decimal, error)
	StringValue(*Context) (string, error)
}

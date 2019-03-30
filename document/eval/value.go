package eval

import (
	"github.com/shopspring/decimal"
)

const (
	TypeBool = iota
	TypeDecimal
	TypeString
)

type Value interface {
	Type(*Context) (int, error)
	BoolValue(*Context) (bool, error)
	DecimalValue(*Context) (decimal.Decimal, error)
	StringValue(*Context) (string, error)
}

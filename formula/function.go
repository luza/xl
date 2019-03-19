package formula

import (
	"strings"

	"xl/document/value"

	"github.com/shopspring/decimal"
)

const maxArguments = 1000

type Function func([]value.Value) (value.Value, error)

type functionDef struct {
	F       Function
	MinArgs int
	MaxArgs int
}

var functions = map[string]functionDef{
	"TRIM": {trim, 1, 1},
	"SUM":  {sum, 1, maxArguments},
}

func trim(args []value.Value) (value.Value, error) {
	s, err := args[0].StringValue()
	if err != nil {
		return value.Value{}, err
	}
	return value.NewStringValue(strings.Trim(s, "\n\r\t ")), nil
}

func sum(args []value.Value) (value.Value, error) {
	s := decimal.Zero
	for i := range args {
		d, err := args[i].DecimalValue()
		if err != nil {
			return value.Value{}, err
		}
		s = s.Add(d)
	}
	return value.NewDecimalValue(s), nil
}

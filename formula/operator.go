package formula

import (
	"strings"

	"xl/document/value"

	"github.com/shopspring/decimal"
)

func evalOperator(ec *value.EvalContext, op string, args ...value.Value) (value.Value, error) {
	var err error
	v := value.Value{}
	// all operands is being casted to first operand type
	switch args[0].Type {
	case value.TypeBool:
		argsBool := make([]bool, len(args))
		for i := range args {
			if argsBool[i], err = args[i].BoolValue(ec); err != nil {
				return v, err
			}
		}
		if v, err = evalBoolOperator(op, argsBool); err != nil {
			return v, err
		}
	case value.TypeDecimal:
		argsDecimal := make([]decimal.Decimal, len(args))
		for i := range args {
			if argsDecimal[i], err = args[i].DecimalValue(ec); err != nil {
				return v, err
			}
		}
		if v, err = evalDecimalOperator(op, argsDecimal); err != nil {
			return v, err
		}
	case value.TypeString:
		argsString := make([]string, len(args))
		for i := range args {
			if argsString[i], err = args[i].StringValue(ec); err != nil {
				return v, err
			}
		}
		if v, err = evalStringOperator(op, argsString); err != nil {
			return v, err
		}
	case value.TypeCell:
		l, err := args[0].Link()
		if err != nil {
			return v, err
		}
		args[0], err = l.Value(ec)
		if err != nil {
			return v, err
		}
		return evalOperator(ec, op, args...)
	default:
		panic("unsupported type")
	}
	return v, nil
}

func evalBoolOperator(op string, args []bool) (value.Value, error) {
	switch op {
	case "=":
		return value.NewBoolValue(args[0] == args[1]), nil
	case "<>":
		return value.NewBoolValue(args[0] != args[1]), nil
	// TRUE > FALSE
	case "<":
		return value.NewBoolValue(!args[0] && args[1]), nil
	case "<=":
		return value.NewBoolValue((!args[0] && args[1]) || args[0] == args[1]), nil
	case ">":
		return value.NewBoolValue(args[0] && !args[1]), nil
	case ">=":
		return value.NewBoolValue((args[0] && !args[1]) || args[0] == args[1]), nil
	case "+":
		if len(args) == 1 {
			// unary
			return value.NewBoolValue(args[0]), nil
		}
		if args[0] && args[1] {
			// TRUE + TRUE = 2
			return value.NewDecimalValue(decimal.NewFromFloat(2)), nil
		} else if args[0] || args[1] {
			// TRUE + FALSE = 1
			return value.NewDecimalValue(decimal.NewFromFloat(1)), nil
		} else {
			// FALSE + FALSE = 0
			return value.NewDecimalValue(decimal.Zero), nil
		}
	case "-":
		if len(args) == 1 {
			// unary
			if args[0] {
				// -TRUE = -1
				return value.NewDecimalValue(decimal.NewFromFloat(-1)), nil
			} else {
				// -FALSE = 0
				return value.NewDecimalValue(decimal.Zero), nil
			}
		} else {
			if args[0] && args[1] {
				// TRUE - TRUE = 0
				return value.NewDecimalValue(decimal.Zero), nil
			} else if args[0] && !args[1] {
				// TRUE - FALSE = 1
				return value.NewDecimalValue(decimal.NewFromFloat(1)), nil
			} else if !args[0] && args[1] {
				// FALSE - TRUE = -1
				return value.NewDecimalValue(decimal.NewFromFloat(-1)), nil
			} else {
				// FALSE - FALSE = 0
				return value.NewDecimalValue(decimal.Zero), nil
			}
		}
	case "*":
		if args[0] && args[1] {
			// TRUE + TRUE = 2
			return value.NewDecimalValue(decimal.NewFromFloat(2)), nil
		} else if args[0] || args[1] {
			// TRUE + FALSE = 1
			return value.NewDecimalValue(decimal.NewFromFloat(1)), nil
		} else {
			// FALSE + FALSE = 0
			return value.NewDecimalValue(decimal.Zero), nil
		}
	case "/":
		if args[0] && args[1] {
			// TRUE - TRUE = 0
			return value.NewDecimalValue(decimal.Zero), nil
		} else if args[0] && !args[1] {
			// TRUE - FALSE = 1
			return value.NewDecimalValue(decimal.NewFromFloat(1)), nil
		} else if !args[0] && args[1] {
			// FALSE - TRUE = -1
			return value.NewDecimalValue(decimal.NewFromFloat(-1)), nil
		} else {
			// FALSE - FALSE = 0
			return value.NewDecimalValue(decimal.Zero), nil
		}
	default:
		panic("unsupported operator")
	}
}

func evalDecimalOperator(op string, args []decimal.Decimal) (value.Value, error) {
	switch op {
	case "=":
		return value.NewBoolValue(args[0].Equal(args[1])), nil
	case "<>":
		return value.NewBoolValue(!args[0].Equal(args[1])), nil
	case "<":
		return value.NewBoolValue(args[0].LessThan(args[1])), nil
	case "<=":
		return value.NewBoolValue(args[0].LessThanOrEqual(args[1])), nil
	case ">":
		return value.NewBoolValue(args[0].GreaterThan(args[1])), nil
	case ">=":
		return value.NewBoolValue(args[0].GreaterThanOrEqual(args[1])), nil
	case "+":
		if len(args) == 1 {
			// unary
			return value.NewDecimalValue(args[0]), nil
		} else {
			return value.NewDecimalValue(args[0].Add(args[1])), nil
		}
	case "-":
		if len(args) == 1 {
			// unary
			return value.NewDecimalValue(args[0].Neg()), nil
		} else {
			return value.NewDecimalValue(args[0].Sub(args[1])), nil
		}
	case "*":
		return value.NewDecimalValue(args[0].Mul(args[1])), nil
	case "/":
		if args[1].Equal(decimal.Zero) {
			return value.Value{}, value.NewError(value.ErrorKindDiv0, "division by zero")
		}
		return value.NewDecimalValue(args[0].Div(args[1])), nil
	default:
		panic("unsupported operator")
	}
}

func evalStringOperator(op string, args []string) (value.Value, error) {
	if len(args) == 1 {
		// unary neg
		return value.Value{}, value.NewError(value.ErrorKindFormula, "arithmetic (%s) on string operand", op)
	}
	res := strings.Compare(args[0], args[1])
	switch op {
	case "=":
		return value.NewBoolValue(args[0] == args[1]), nil
	case "<>":
		return value.NewBoolValue(args[0] != args[1]), nil
	case "<":
		return value.NewBoolValue(res < 0), nil
	case "<=":
		return value.NewBoolValue(res <= 0), nil
	case ">":
		return value.NewBoolValue(res > 0), nil
	case ">=":
		return value.NewBoolValue(res >= 0), nil
	case "+", "-", "*", "/":
		return value.Value{}, value.NewError(value.ErrorKindFormula, "arithmetic (%s) on string operand", op)
	default:
		panic("unsupported operator")
	}
}

func evalFunc(ec *value.EvalContext, name string, args []value.Value) (value.Value, error) {
	if f, ok := functions[name]; ok {
		if len(args) < f.MinArgs || len(args) > f.MaxArgs {
			return value.Value{}, value.NewError(value.ErrorKindFormula, "function %s accepts from %d to %d arguments, %d provided",
				name, f.MinArgs, f.MaxArgs, len(args))
		}
		return f.F(ec, args)
	} else {
		return value.Value{}, value.NewError(value.ErrorKindFormula, "function %s does not exist", name)
	}
}

package formula

import (
	"strings"

	"xl/document/eval"

	"github.com/shopspring/decimal"
)

func evalOperator(ec *eval.Context, op string, args ...eval.Value) (eval.Value, error) {
	v := eval.NewEmptyValue()
	t, err := args[0].Type(ec)
	if err != nil {
		return v, err
	}
	// all operands is being casted to first operand type
	switch t {
	case eval.TypeBool:
		argsBool := make([]bool, len(args))
		for i := range args {
			if argsBool[i], err = args[i].BoolValue(ec); err != nil {
				return v, err
			}
		}
		if v, err = evalBoolOperator(op, argsBool); err != nil {
			return v, err
		}
	case eval.TypeEmpty, eval.TypeDecimal:
		argsDecimal := make([]decimal.Decimal, len(args))
		for i := range args {
			if argsDecimal[i], err = args[i].DecimalValue(ec); err != nil {
				return v, err
			}
		}
		if v, err = evalDecimalOperator(op, argsDecimal); err != nil {
			return v, err
		}
	case eval.TypeString:
		argsString := make([]string, len(args))
		for i := range args {
			if argsString[i], err = args[i].StringValue(ec); err != nil {
				return v, err
			}
		}
		if v, err = evalStringOperator(op, argsString); err != nil {
			return v, err
		}
	default:
		panic("unsupported type")
	}
	return v, nil
}

func evalBoolOperator(op string, args []bool) (eval.Value, error) {
	switch op {
	case "=":
		return eval.NewBoolValue(args[0] == args[1]), nil
	case "<>":
		return eval.NewBoolValue(args[0] != args[1]), nil
	// TRUE > FALSE
	case "<":
		return eval.NewBoolValue(!args[0] && args[1]), nil
	case "<=":
		return eval.NewBoolValue((!args[0] && args[1]) || args[0] == args[1]), nil
	case ">":
		return eval.NewBoolValue(args[0] && !args[1]), nil
	case ">=":
		return eval.NewBoolValue((args[0] && !args[1]) || args[0] == args[1]), nil
	case "+":
		if len(args) == 1 {
			// unary
			return eval.NewBoolValue(args[0]), nil
		}
		if args[0] && args[1] {
			// TRUE + TRUE = 2
			return eval.NewDecimalValue(decimal.NewFromFloat(2)), nil
		} else if args[0] || args[1] {
			// TRUE + FALSE = 1
			return eval.NewDecimalValue(decimal.NewFromFloat(1)), nil
		} else {
			// FALSE + FALSE = 0
			return eval.NewDecimalValue(decimal.Zero), nil
		}
	case "-":
		if len(args) == 1 {
			// unary
			if args[0] {
				// -TRUE = -1
				return eval.NewDecimalValue(decimal.NewFromFloat(-1)), nil
			} else {
				// -FALSE = 0
				return eval.NewDecimalValue(decimal.Zero), nil
			}
		} else {
			if args[0] && args[1] {
				// TRUE - TRUE = 0
				return eval.NewDecimalValue(decimal.Zero), nil
			} else if args[0] && !args[1] {
				// TRUE - FALSE = 1
				return eval.NewDecimalValue(decimal.NewFromFloat(1)), nil
			} else if !args[0] && args[1] {
				// FALSE - TRUE = -1
				return eval.NewDecimalValue(decimal.NewFromFloat(-1)), nil
			} else {
				// FALSE - FALSE = 0
				return eval.NewDecimalValue(decimal.Zero), nil
			}
		}
	case "*":
		if args[0] && args[1] {
			// TRUE + TRUE = 2
			return eval.NewDecimalValue(decimal.NewFromFloat(2)), nil
		} else if args[0] || args[1] {
			// TRUE + FALSE = 1
			return eval.NewDecimalValue(decimal.NewFromFloat(1)), nil
		} else {
			// FALSE + FALSE = 0
			return eval.NewDecimalValue(decimal.Zero), nil
		}
	case "/":
		if args[0] && args[1] {
			// TRUE - TRUE = 0
			return eval.NewDecimalValue(decimal.Zero), nil
		} else if args[0] && !args[1] {
			// TRUE - FALSE = 1
			return eval.NewDecimalValue(decimal.NewFromFloat(1)), nil
		} else if !args[0] && args[1] {
			// FALSE - TRUE = -1
			return eval.NewDecimalValue(decimal.NewFromFloat(-1)), nil
		} else {
			// FALSE - FALSE = 0
			return eval.NewDecimalValue(decimal.Zero), nil
		}
	case "^":
		if args[0] || !args[1] {
			// TRUE^FALSE = 1, TRUE^TRUE = 1, FALSE^FALSE = 1 (in Excel FALSE^FALSE = error)
			return eval.NewDecimalValue(decimal.NewFromFloat(1)), nil
		} else {
			// FALSE^TRUE = 0
			return eval.NewDecimalValue(decimal.Zero), nil
		}
	default:
		panic("unsupported operator")
	}
}

func evalDecimalOperator(op string, args []decimal.Decimal) (eval.Value, error) {
	switch op {
	case "=":
		return eval.NewBoolValue(args[0].Equal(args[1])), nil
	case "<>":
		return eval.NewBoolValue(!args[0].Equal(args[1])), nil
	case "<":
		return eval.NewBoolValue(args[0].LessThan(args[1])), nil
	case "<=":
		return eval.NewBoolValue(args[0].LessThanOrEqual(args[1])), nil
	case ">":
		return eval.NewBoolValue(args[0].GreaterThan(args[1])), nil
	case ">=":
		return eval.NewBoolValue(args[0].GreaterThanOrEqual(args[1])), nil
	case "+":
		if len(args) == 1 {
			// unary
			return eval.NewDecimalValue(args[0]), nil
		} else {
			return eval.NewDecimalValue(args[0].Add(args[1])), nil
		}
	case "-":
		if len(args) == 1 {
			// unary
			return eval.NewDecimalValue(args[0].Neg()), nil
		} else {
			return eval.NewDecimalValue(args[0].Sub(args[1])), nil
		}
	case "*":
		return eval.NewDecimalValue(args[0].Mul(args[1])), nil
	case "/":
		if args[1].Equal(decimal.Zero) {
			return eval.NewEmptyValue(), eval.NewError(eval.ErrorKindDiv0, "division by zero")
		}
		return eval.NewDecimalValue(args[0].Div(args[1])), nil
	case "^":
		return eval.NewDecimalValue(args[0].Pow(args[1])), nil
	default:
		panic("unsupported operator")
	}
}

func evalStringOperator(op string, args []string) (eval.Value, error) {
	if len(args) == 1 {
		// unary neg
		return eval.NewEmptyValue(), eval.NewError(eval.ErrorKindFormula, "arithmetic (%s) on string operand", op)
	}
	res := strings.Compare(args[0], args[1])
	switch op {
	case "=":
		return eval.NewBoolValue(args[0] == args[1]), nil
	case "<>":
		return eval.NewBoolValue(args[0] != args[1]), nil
	case "<":
		return eval.NewBoolValue(res < 0), nil
	case "<=":
		return eval.NewBoolValue(res <= 0), nil
	case ">":
		return eval.NewBoolValue(res > 0), nil
	case ">=":
		return eval.NewBoolValue(res >= 0), nil
	case "+", "-", "*", "/", "^":
		return eval.NewEmptyValue(), eval.NewError(eval.ErrorKindFormula, "arithmetic (%s) on string operand", op)
	default:
		panic("unsupported operator")
	}
}

func evalFunc(ec *eval.Context, name string, args []eval.Value) (eval.Value, error) {
	if f, ok := functions[name]; ok {
		if len(args) < f.MinArgs || len(args) > f.MaxArgs {
			return eval.NewEmptyValue(), eval.NewError(eval.ErrorKindFormula, "function %s accepts from %d to %d arguments, %d provided",
				name, f.MinArgs, f.MaxArgs, len(args))
		}
		return f.F(ec, args)
	} else {
		return eval.NewEmptyValue(), eval.NewError(eval.ErrorKindFormula, "function %s does not exist", name)
	}
}

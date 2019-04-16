package formula

import (
	"xl/document/eval"

	"github.com/shopspring/decimal"
)

// Методые для создания функции, которая производит вычисления, описанные формулой.

func (e *Expression) BuildFunc() (Function, int) {
	return e.Equality.BuildFunc()
}

func (e *Equality) BuildFunc() (Function, int) {
	if e.Next == nil {
		return e.Comparison.BuildFunc()
	}
	subFunc1, consumedArgs1 := e.Comparison.BuildFunc()
	subFunc2, consumedArgs2 := e.Next.BuildFunc()
	return evalBinaryOperator(
		e.Op,
		subFunc1, consumedArgs1,
		subFunc2, consumedArgs2,
	)
}

func (e *Comparison) BuildFunc() (Function, int) {
	if e.Next == nil {
		return e.Addition.BuildFunc()
	}
	subFunc1, consumedArgs1 := e.Addition.BuildFunc()
	subFunc2, consumedArgs2 := e.Next.BuildFunc()
	return evalBinaryOperator(
		e.Op,
		subFunc1, consumedArgs1,
		subFunc2, consumedArgs2,
	)
}

func (e *Addition) BuildFunc() (Function, int) {
	if e.Next == nil {
		return e.Multiplication.BuildFunc()
	}
	subFunc1, consumedArgs1 := e.Multiplication.BuildFunc()
	subFunc2, consumedArgs2 := e.Next.BuildFunc()
	return evalBinaryOperator(
		e.Op,
		subFunc1, consumedArgs1,
		subFunc2, consumedArgs2,
	)
}

func (e *Multiplication) BuildFunc() (Function, int) {
	if e.Next == nil {
		return e.Power.BuildFunc()
	}
	subFunc1, consumedArgs1 := e.Power.BuildFunc()
	subFunc2, consumedArgs2 := e.Next.BuildFunc()
	return evalBinaryOperator(
		e.Op,
		subFunc1, consumedArgs1,
		subFunc2, consumedArgs2,
	)
}

func (e *Power) BuildFunc() (Function, int) {
	if len(e.Exponent) == 0 {
		return e.Base.BuildFunc()
	}
	subFunc1, consumedArgs1 := e.Base.BuildFunc()
	for _, x := range e.Exponent {
		subFunc2, consumedArgs2 := x.BuildFunc()
		subFunc1, consumedArgs1 = evalBinaryOperator(
			"^",
			subFunc1, consumedArgs1,
			subFunc2, consumedArgs2,
		)
	}
	return subFunc1, consumedArgs1
}

func (e *Unary) BuildFunc() (Function, int) {
	if e.Primary != nil {
		return e.Primary.BuildFunc()
	} else {
		subFunc, consumedArgs := e.Unary.BuildFunc()
		return evalUnaryOperator(
			e.Op,
			subFunc, consumedArgs,
		)
	}
}

func (e *Primary) BuildFunc() (Function, int) {
	if e.SubExpression != nil {
		return e.SubExpression.BuildFunc()
	} else if e.Func != nil {
		return e.Func.BuildFunc()
	} else if e.Boolean != nil {
		f := func(*eval.Context, []eval.Value) (eval.Value, error) {
			return eval.NewBoolValue(bool(*e.Boolean)), nil
		}
		return f, 0
	} else if e.Number != nil {
		f := func(*eval.Context, []eval.Value) (eval.Value, error) {
			return eval.NewDecimalValue(decimal.NewFromFloat(*e.Number)), nil
		}
		return f, 0
	} else if e.String != nil {
		f := func(*eval.Context, []eval.Value) (eval.Value, error) {
			return eval.NewStringValue(string(*e.String)), nil
		}
		return f, 0
	} else {
		f := func(ec *eval.Context, args []eval.Value) (eval.Value, error) {
			if len(args) == 0 {
				panic("too few arguments")
			}
			return args[0], nil
		}
		return f, 1
	}
}

func (e *Func) BuildFunc() (Function, int) {
	totalConsumedArgs := 0
	subFunc := make([]Function, len(e.Arguments))
	consumedArgs := make([]int, len(e.Arguments))
	for i, a := range e.Arguments {
		subFunc[i], consumedArgs[i] = a.BuildFunc()
		totalConsumedArgs += consumedArgs[i]
	}
	f := func(ec *eval.Context, args []eval.Value) (eval.Value, error) {
		var err error
		values := make([]eval.Value, len(e.Arguments))
		ca := 0
		for i := range e.Arguments {
			values[i], err = subFunc[i](ec, args[ca:])
			if err != nil {
				return eval.NewEmptyValue(), err
			}
			ca += consumedArgs[i]
		}
		return evalFunc(ec, string(e.Name), values)
	}
	return f, totalConsumedArgs
}

func evalBinaryOperator(op string, f1 Function, consumedArgs1 int, f2 Function, consumedArgs2 int) (Function, int) {
	f := func(ec *eval.Context, args []eval.Value) (eval.Value, error) {
		var v1, v2 eval.Value
		var err error
		if v1, err = f1(ec, args); err != nil {
			return eval.NewEmptyValue(), err
		}
		if v2, err = f2(ec, args[consumedArgs1:]); err != nil {
			return eval.NewEmptyValue(), err
		}
		return evalOperator(ec, op, v1, v2)
	}
	return f, consumedArgs1 + consumedArgs2
}

func evalUnaryOperator(op string, f1 Function, consumedArgs1 int) (Function, int) {
	f := func(ec *eval.Context, args []eval.Value) (eval.Value, error) {
		v, err := f1(ec, args)
		if err != nil {
			return eval.NewEmptyValue(), err
		}
		return evalOperator(ec, op, v)
	}
	return f, consumedArgs1
}

package formula

import (
	"strconv"
)

type OutputFunc func(string, int)

const (
	OutputTypeSymbol = iota
	OutputTypeWhitespace
	OutputTypeOperator
	OutputTypeNumber
	OutputTypeBoolean
	OutputTypeString
	OutputTypeFunction
	OutputTypeSheet
	OutputTypeCell
)

func (e *Expression) Output(of OutputFunc) {
	of("=", OutputTypeSymbol)
	e.Equality.Output(of)
}

func (e *Equality) Output(of OutputFunc) {
	e.Comparison.Output(of)
	if e.Op != "" && e.Next != nil {
		of(e.Op, OutputTypeOperator)
		e.Next.Output(of)
	}
}

func (e *Comparison) Output(of OutputFunc) {
	e.Addition.Output(of)
	if e.Op != "" && e.Next != nil {
		of(e.Op, OutputTypeOperator)
		e.Next.Output(of)
	}
}

func (e *Addition) Output(of OutputFunc) {
	e.Multiplication.Output(of)
	if e.Op != "" && e.Next != nil {
		of(e.Op, OutputTypeOperator)
		e.Next.Output(of)
	}
}

func (e *Multiplication) Output(of OutputFunc) {
	e.Power.Output(of)
	if e.Op != "" && e.Next != nil {
		of(e.Op, OutputTypeOperator)
		e.Next.Output(of)
	}
}

func (e *Power) Output(of OutputFunc) {
	e.Base.Output(of)
	for _, x := range e.Exponent {
		of("^", OutputTypeOperator)
		x.Output(of)
	}
}

func (e *Unary) Output(of OutputFunc) {
	if e.Primary != nil {
		e.Primary.Output(of)
	} else {
		of(e.Op, OutputTypeOperator)
		e.Unary.Output(of)
	}
}

func (e *Primary) Output(of OutputFunc) {
	if e.SubExpression != nil {
		of("(", OutputTypeSymbol)
		e.SubExpression.Output(of)
		of(")", OutputTypeSymbol)
	}
	if e.Number != nil {
		of(strconv.FormatFloat(*e.Number, 'f', -1, 64), OutputTypeNumber)
	} else if e.String != nil {
		of("\"", OutputTypeSymbol)
		of(string(*e.String), OutputTypeString)
		of("\"", OutputTypeSymbol)
	} else if e.Boolean != nil {
		if *e.Boolean {
			of("TRUE", OutputTypeBoolean)
		} else {
			of("FALSE", OutputTypeBoolean)
		}
	} else if e.Func != nil {
		e.Func.Output(of)
	} else if e.Variable != nil {
		e.Variable.Output(of)
	}
}

func (e *Func) Output(of OutputFunc) {
	of(string(e.Name), OutputTypeFunction)
	of("(", OutputTypeSymbol)
	for i, a := range e.Arguments {
		a.Output(of)
		if i < len(e.Arguments)-1 {
			of(";", OutputTypeSymbol)
			of(" ", OutputTypeWhitespace)
		}
	}
	of(")", OutputTypeSymbol)
}

func (e *Variable) Output(of OutputFunc) {
	e.Cell.Output(of)
	if e.CellTo != nil {
		of(":", OutputTypeSymbol)
		e.CellTo.Output(of)
	}
}

func (e *Cell) Output(of OutputFunc) {
	if e.Sheet != nil {
		// FIXME: use '' only if necessary
		of("'", OutputTypeSymbol)
		of(string(*e.Sheet), OutputTypeSheet)
		of("'", OutputTypeSymbol)
		of("!", OutputTypeSymbol)
	}
	of(e.CellName, OutputTypeCell)
}

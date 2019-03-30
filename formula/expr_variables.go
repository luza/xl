package formula

func (e *Expression) Variables() []*Variable {
	return e.Equality.Variables()
}

func (e *Equality) Variables() []*Variable {
	vars := e.Comparison.Variables()
	if e.Next != nil {
		vars = append(vars, e.Next.Variables()...)
	}
	return vars
}

func (e *Comparison) Variables() []*Variable {
	vars := e.Addition.Variables()
	if e.Next != nil {
		vars = append(vars, e.Next.Variables()...)
	}
	return vars
}

func (e *Addition) Variables() []*Variable {
	vars := e.Multiplication.Variables()
	if e.Next != nil {
		vars = append(vars, e.Next.Variables()...)
	}
	return vars
}

func (e *Multiplication) Variables() []*Variable {
	vars := e.Power.Variables()
	if e.Next != nil {
		vars = append(vars, e.Next.Variables()...)
	}
	return vars
}

func (e *Power) Variables() []*Variable {
	vars := e.Unary.Variables()
	if e.Next != nil {
		vars = append(vars, e.Next.Variables()...)
	}
	return vars
}

func (e *Unary) Variables() []*Variable {
	if e.Primary != nil {
		return e.Primary.Variables()
	} else {
		return e.Unary.Variables()
	}
}

func (e *Primary) Variables() []*Variable {
	if e.Variable != nil {
		return []*Variable{e.Variable}
	} else if e.SubExpression != nil {
		return e.SubExpression.Variables()
	} else if e.Func != nil {
		return e.Func.Variables()
	} else {
		return nil
	}
}

func (e *Func) Variables() []*Variable {
	var vars []*Variable
	for _, a := range e.Arguments {
		vars = append(vars, a.Variables()...)
	}
	return vars
}

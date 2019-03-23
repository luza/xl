package value

type EvalContext struct {
	LinkRegistry LinkRegistryInterface

	visitedCells map[*Link]bool
}

func NewEvalContext(lr LinkRegistryInterface) *EvalContext {
	return &EvalContext{
		LinkRegistry: lr,
		visitedCells: make(map[*Link]bool),
	}
}

func (ec *EvalContext) Reset() {
	ec.visitedCells = make(map[*Link]bool)
}

func (ec *EvalContext) AddVisited(l *Link) {
	ec.visitedCells[l] = true
}

func (ec *EvalContext) Visited(l *Link) bool {
	_, ok := ec.visitedCells[l]
	return ok
}

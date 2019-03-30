package eval

type Context struct {
	DataProvider RefRegistryInterface
	visitedCells map[*CellRef]bool
}

func NewContext(dp RefRegistryInterface) *Context {
	ec := &Context{
		DataProvider: dp,
	}
	ec.Reset()
	return ec
}

func (ec *Context) Reset() {
	ec.visitedCells = make(map[*CellRef]bool)
}

func (ec *Context) AddVisited(r *CellRef) {
	ec.visitedCells[r] = true
}

func (ec *Context) Visited(r *CellRef) bool {
	_, ok := ec.visitedCells[r]
	return ok
}

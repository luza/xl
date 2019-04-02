package eval

type Context struct {
	DataProvider    RefRegistryInterface
	CurrentSheetIdx int
	visitedCells    []*CellRef
}

func NewContext(dp RefRegistryInterface, currentSheetIdx int) *Context {
	ec := &Context{
		DataProvider:    dp,
		CurrentSheetIdx: currentSheetIdx,
	}
	return ec
}

func (ec *Context) AddVisited(r *CellRef) int {
	oldLen := len(ec.visitedCells)
	ec.visitedCells = append(ec.visitedCells, r)
	return oldLen
}

func (ec *Context) ResetVisited(i int) {
	ec.visitedCells = ec.visitedCells[:i]
}

func (ec *Context) Visited(r *CellRef) bool {
	for i := range ec.visitedCells {
		if ec.visitedCells[i] == r {
			return true
		}
	}
	return false
}

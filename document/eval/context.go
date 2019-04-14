package eval

type Context struct {
	DataProvider    RefRegistryInterface
	CurrentSheetIdx int
	visitedCells    []CellAddress
}

func NewContext(dp RefRegistryInterface, currentSheetIdx int) *Context {
	ec := &Context{
		DataProvider:    dp,
		CurrentSheetIdx: currentSheetIdx,
	}
	return ec
}

func (ec *Context) AddVisited(cell CellAddress) int {
	oldLen := len(ec.visitedCells)
	ec.visitedCells = append(ec.visitedCells, cell)
	return oldLen
}

func (ec *Context) Len() int {
	return len(ec.visitedCells)
}

func (ec *Context) ResetVisited(i int) {
	ec.visitedCells = ec.visitedCells[:i]
}

func (ec *Context) Visited(cell CellAddress) bool {
	for i := range ec.visitedCells {
		if ec.visitedCells[i] == cell {
			return true
		}
	}
	return false
}

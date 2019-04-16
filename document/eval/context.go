package eval

// Контекст вычисления. Экземпляр этой структуры програсывается через все вычисления
// в рамках одной формулы и служит двум целям:
// - считать посещенные ячейки при переходе по Ссылкам, чтобы пресекать циклические ссылки
// - предоставлять доступ к документу при разрешении Ссылок

type Context struct {
	// Это делегат, предоставляющий методы разрешения ссылок.
	DataProvider RefRegistryInterface

	// Лист, в контексте которого делается вычисление формулы.
	CurrentSheetIdx int

	// Слайс для хранения посещенных ячеек.
	visitedCells []CellAddress
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

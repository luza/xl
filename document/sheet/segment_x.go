package sheet

// Extrapolation segment.

// Когда пользователь растягивает ячейки вниз, вправо, влево или вверх, на месте
// выделенных ячеек возникает экстраполяционный сегмент (xSegment). Одна из принадлежащих
// ему ячеек объявляется ключевой, и для нее задается значение. Остальные ячейки сегмента
// вычисляются в момент запроса; для вычисления произвольной клетки сегмента берется
// значение его ключевой ячейки, ссылки в которой сдвинуты соразмерно смещению
// запрошенной ячейки относительно ключевой.

type xSegment struct {
	baseSegment

	// Координаты ключевой ячейки внутри сегмента.
	keyX int
	keyY int

	// Значение ключевой ячейки.
	keyCell Cell
}

func newXSegment(x, y, width, height, keyX, keyY int, keyCell Cell) Segment {
	return &xSegment{
		baseSegment: baseSegment{
			size: Rect{
				X:      x,
				Y:      y,
				Width:  width,
				Height: height,
			},
		},
		keyX:    keyX,
		keyY:    keyY,
		keyCell: keyCell,
	}
}

// Cell returns cell under the given X and Y.
func (s *xSegment) Cell(x, y int) *Cell {
	localX, localY := x-s.size.X, y-s.size.Y
	if localX == s.keyX && localY == s.keyY {
		return &s.keyCell
	}
	return NewCellAsCopyWithOffset(&s.keyCell, localX-s.keyX, localY-s.keyY)
}

// SetCell fills new key cell.
func (s *xSegment) SetCell(x, y int, cell *Cell) {
	localX, localY := x-s.size.X, y-s.size.Y
	if localX == s.keyX && localY == s.keyY {
		s.keyCell = *cell
	}
	panic("writing cell is possible only for key cell")
}

func (s *xSegment) InsertEmptyRow(y int) {
	panic("not supported")
}

func (s *xSegment) InsertEmptyCol(x int) {
	panic("not supported")
}

func (s *xSegment) DeleteRow(y int) {
	panic("not supported")
}

func (s *xSegment) DeleteCol(x int) {
	panic("not supported")
}

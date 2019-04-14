package sheet

// Extrapolation segment.

type xSegment struct {
	baseSegment

	keyX    int
	keyY    int
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
	//FIXME
	//return NewCellAsCopyWithOffset(&s.keyCell, localX-s.keyX, localY-s.keyY)
	return nil
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

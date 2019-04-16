package sheet

// Static segment.

// Статичный сегмент состоит из клеток с фиксированным значением.

// staticSegment is a segment of static values.
type staticSegment struct {
	baseSegment
	Cells [][]Cell
}

func newStaticSegment(x, y, w, h int, cells [][]Cell) Segment {
	return &staticSegment{
		baseSegment: baseSegment{
			size: Rect{
				X:      x,
				Y:      y,
				Width:  w,
				Height: h,
			},
		},
		Cells: cells,
	}
}

// Cell returns cell under the given X and Y.
func (s *staticSegment) Cell(x, y int) *Cell {
	return &s.Cells[x-s.size.X][y-s.size.Y]
}

// SetCell fills new cell on position of given X and Y.
func (s *staticSegment) SetCell(x, y int, cell *Cell) {
	s.Cells[x-s.size.X][y-s.size.Y] = *cell
}

func (s *staticSegment) InsertEmptyRow(y int) {
	for x := 0; x < s.size.Width; x++ {
		s.Cells[x] = append(s.Cells[x], Cell{})
		copy(s.Cells[x][y+1:], s.Cells[x][y:])
		s.Cells[x][y] = *NewCellEmpty()
	}
	s.size.Height++
}

func (s *staticSegment) InsertEmptyCol(x int) {
	s.Cells = append(s.Cells, []Cell{})
	copy(s.Cells[x+1:], s.Cells[x:])
	col := make([]Cell, s.size.Height)
	for y := 0; y < s.size.Height; y++ {
		col[y] = *NewCellEmpty()
	}
	s.Cells[x] = col
	s.size.Width++
}

func (s *staticSegment) DeleteRow(y int) {
	for x := 0; x < s.size.Width; x++ {
		copy(s.Cells[x][y:], s.Cells[x][y+1:])
		s.Cells[x][len(s.Cells[x])-1] = Cell{}
		s.Cells[x] = s.Cells[x][:len(s.Cells[x])-1]
	}
	s.size.Height--
}

func (s *staticSegment) DeleteCol(x int) {
	copy(s.Cells[x:], s.Cells[x+1:])
	s.Cells[len(s.Cells)-1] = nil
	s.Cells = s.Cells[:len(s.Cells)-1]
	s.size.Width--
}

package sheet

type Segment interface {
	Contains(x, y int) bool
	ContainsX(x int) bool
	ContainsY(x int) bool
	Cell(x, y int) *Cell
	Size() Rect
	Move(newX, newY int)
	InsertRow(y int)
	InsertCol(x int)
	DeleteRow(y int)
	DeleteCol(x int)
}

// Base segment.

// baseSegment contains properties and methods common for all segment types.
type baseSegment struct {
	Segment
	size Rect
}

// Contains checks if given X and Y belong to the segment.
func (s *baseSegment) Contains(x, y int) bool {
	return s.ContainsX(x) && s.ContainsY(y)
}

func (s *baseSegment) ContainsX(x int) bool {
	return x >= s.size.X && x < s.size.X+s.size.Width
}

func (s *baseSegment) ContainsY(y int) bool {
	return y >= s.size.Y && y < s.size.Y+s.size.Height
}

// Size returns rect of segment bounds.
func (s *baseSegment) Size() Rect {
	return s.size
}

func (s *baseSegment) Move(newX, newY int) {
	s.size.X = newX
	s.size.Y = newY
}

// Static segment.

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

func (s *staticSegment) InsertRow(y int) {
	for x := 0; x < s.size.Width; x++ {
		s.Cells[x] = append(s.Cells[x], Cell{})
		copy(s.Cells[x][y+1:], s.Cells[x][y:])
		s.Cells[x][y] = *NewCellEmpty()
	}
	s.size.Height++
}

func (s *staticSegment) InsertCol(x int) {
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

// Extrapolation segment.

//type xpSegment struct {
//	baseSegment
//
//	Xps [][]xp1
//}

// xp1 is one-dimensional extrapolation rule.
//type xp1 struct {
//	direction int
//	rule      string
//}

//func newXpSegment(x, y, w, h int) Segment {
//	return &xpSegment{
//		baseSegment: baseSegment{
//			X:      x,
//			Y:      y,
//			Width:  w,
//			Height: h,
//		},
//	}
//}

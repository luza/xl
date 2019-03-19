package sheet

type Segment interface {
	Contains(x, y int) bool
	Cell(x, y int) *Cell
	Size() Rect
}

// Base segment.

// baseSegment contains properties and methods common for all segment types.
type baseSegment struct {
	Segment
	size Rect
}

// Contains checks if given X and Y belong to the segment.
func (s *baseSegment) Contains(x, y int) bool {
	return x >= s.size.X && x < s.size.X+s.size.Width && y >= s.size.Y && y < s.size.Y+s.size.Height
}

// Size returns rect of segment bounds.
func (s *baseSegment) Size() Rect {
	return s.size
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

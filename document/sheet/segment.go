package sheet

type Segment interface {
	Contains(x, y int) bool
	ContainsX(x int) bool
	ContainsY(x int) bool
	Cell(x, y int) *Cell
	SetCell(x, y int, cell *Cell)
	Size() Rect
	Move(newX, newY int)
	InsertEmptyRow(y int)
	InsertEmptyCol(x int)
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

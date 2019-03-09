package sheet

//type xp1 struct {
//	direction int
//	rule      string
//}

type Segment interface {
	Contains(x, y int) bool
	Cell(x, y int) *Cell
	Size() Rect
}

type baseSegment struct {
	Segment
	size Rect
}

func (s *baseSegment) Contains(x, y int) bool {
	return x >= s.size.X && x < s.size.X+s.size.Width && y >= s.size.Y && y < s.size.Y+s.size.Height
}

func (s *baseSegment) Size() Rect {
	return s.size
}

type staticSegment struct {
	baseSegment
	Cells [][]Cell
}

func (s *staticSegment) Cell(x, y int) *Cell {
	return &s.Cells[x-s.size.X][y-s.size.Y]
}

func (s *staticSegment) SetCell(x, y int, cell *Cell) {
	s.Cells[x-s.size.X][y-s.size.Y] = *cell
}

//type xpSegment struct {
//	baseSegment
//
//	Xps [][]xp1
//}

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

//func newExtrapolationSegment(x, y, w, h int) Segment {
//	return &xpSegment{
//		baseSegment: baseSegment{
//			X:      x,
//			Y:      y,
//			Width:  w,
//			Height: h,
//		},
//	}
//}

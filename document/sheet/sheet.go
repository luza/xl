package sheet

const (
	CellDefaultWidth  = 80
	CellDefaultHeight = 10
	CellMaxWidth      = CellDefaultWidth * 10
	CellMaxHeight     = CellDefaultHeight * 10
)

type Cursor struct {
	X int
	Y int
}

type Viewport struct {
	Left int
	Top  int
}

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r *Rect) Right() int {
	return r.X + r.Width
}

func (r *Rect) Bottom() int {
	return r.Y + r.Height
}

type Sheet struct {
	Name     string
	Cursor   Cursor
	Viewport Viewport
	Size     Rect
	Segments []Segment

	colSizes map[int]int
	rowSizes map[int]int
}

func New(name string) *Sheet {
	return &Sheet{
		Name:     name,
		Cursor:   Cursor{0, 0},
		Viewport: Viewport{0, 0},
		Size:     Rect{0, 0, 0, 0},

		colSizes: make(map[int]int),
		rowSizes: make(map[int]int),
	}
}

func (s *Sheet) ColSize(n int) int {
	if val, ok := s.colSizes[n]; ok {
		return val
	}
	return CellDefaultWidth
}

func (s *Sheet) SetColSize(n, size int) {
	if size < 1 || size > CellMaxWidth {
		return
	}
	s.colSizes[n] = size
}

func (s *Sheet) RowSize(n int) int {
	if val, ok := s.rowSizes[n]; ok {
		return val
	}
	return CellDefaultHeight
}

func (s *Sheet) AddStaticSegment(x, y, width, height int, cells [][]Cell) Segment {
	if len(cells) == 0 || len(cells[0]) == 0 {
		return nil
	}
	// TODO: check intersections?
	// TODO: check merge possible
	segment := newStaticSegment(x, y, width, height, cells)
	s.Segments = append(s.Segments, segment)

	// adjust sheet size
	if x < s.Size.X {
		s.Size.X = x
	}
	if y < s.Size.Y {
		s.Size.Y = y
	}
	if x+width > s.Size.Width {
		s.Size.Width = x + width
	}
	if y+height > s.Size.Height {
		s.Size.Height = y + height
	}

	return segment
}

func (s *Sheet) Cell(x, y int) *Cell {
	for _, segment := range s.Segments {
		if segment.Contains(x, y) {
			return segment.Cell(x, y)
		}
	}
	return nil
}

func (s *Sheet) CellUnderCursor() *Cell {
	return s.Cell(s.Cursor.X, s.Cursor.Y)
}

func (s *Sheet) SetCell(x, y int, cell *Cell) {
	segment := s.FindSegment(x, y)
	if segment == nil {
		// create new Segment
		segment = s.AddStaticSegment(x, y, 1, 1, [][]Cell{{*cell}})
	}
	if st, ok := segment.(*staticSegment); ok {
		st.SetCell(x, y, cell)
	} else {
		// TODO: other types of Segments
	}
}

func (s *Sheet) FindSegment(x, y int) Segment {
	for _, segment := range s.Segments {
		if segment.Contains(x, y) {
			return segment
		}
	}
	return nil
}

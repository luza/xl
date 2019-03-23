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

// MaxX returns maximum X belonging to rect.
func (r *Rect) MaxX() int {
	return r.X + r.Width - 1
}

// MaxY returns maximum Y belonging to rect.
func (r *Rect) MaxY() int {
	return r.Y + r.Height - 1
}

type Sheet struct {
	Idx      int
	Title    string
	Cursor   Cursor
	Viewport Viewport
	Size     Rect
	Segments []Segment

	colSizes map[int]int
	rowSizes map[int]int
}

func New(idx int, name string) *Sheet {
	return &Sheet{
		Idx:      idx,
		Title:    name,
		Cursor:   Cursor{0, 0},
		Viewport: Viewport{0, 0},
		Size:     Rect{0, 0, 0, 0},

		colSizes: make(map[int]int),
		rowSizes: make(map[int]int),
	}
}

// ColSize returns width of a column in pixels.
func (s *Sheet) ColSize(n int) int {
	if val, ok := s.colSizes[n]; ok {
		return val
	}
	return CellDefaultWidth
}

// SetColSize sets the new width for a column in pixels.
func (s *Sheet) SetColSize(n, size int) {
	if size < 1 || size > CellMaxWidth {
		return
	}
	s.colSizes[n] = size
}

// RowSize returns height of a row in pixels.
func (s *Sheet) RowSize(n int) int {
	if val, ok := s.rowSizes[n]; ok {
		return val
	}
	return CellDefaultHeight
}

// AddStaticSegment creates a new Static segment and will with the given cells matrix.
// TODO: check intersections. Will be such a case?
// TODO: new segment need to be merged with the existing if possible
func (s *Sheet) AddStaticSegment(x, y, width, height int, cells [][]Cell) Segment {
	if len(cells) == 0 || len(cells[0]) == 0 {
		return nil
	}
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

// Cell returns the cell for given X and Y.
func (s *Sheet) Cell(x, y int) *Cell {
	for _, segment := range s.Segments {
		if segment.Contains(x, y) {
			return segment.Cell(x, y)
		}
	}
	return nil
}

// CellUnderCursor returns the cell the cursor points to.
func (s *Sheet) CellUnderCursor() *Cell {
	return s.Cell(s.Cursor.X, s.Cursor.Y)
}

// SetCell fills the cell with new data.
// If no such cell exists yet, created a new segment.
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

// FindSegment iterates over segments to find the one containing cell with given X and Y.
func (s *Sheet) FindSegment(x, y int) Segment {
	for _, segment := range s.Segments {
		if segment.Contains(x, y) {
			return segment
		}
	}
	return nil
}

func (s *Sheet) InsertRow(y int) {
	for _, segment := range s.Segments {
		size := segment.Size()
		if segment.ContainsY(y) {
			// segments that inserting line intersects need to be splitted or expanded
			segment.InsertRow(y - size.Y)
		} else if size.Y > y {
			// segments laying below inserting line need to be shifted down
			segment.Move(size.X, size.Y+1)
		}
	}
}

func (s *Sheet) InsertCol(x int) {
	for _, segment := range s.Segments {
		size := segment.Size()
		if segment.ContainsX(x) {
			// segments that inserting column intersects need to be splitted or expanded
			segment.InsertCol(x - size.X)
		} else if size.X > x {
			// segments laying right of inserting column need to be shifted right
			segment.Move(size.X+1, size.Y)
		}
	}
}

func (s *Sheet) DeleteRow(y int) {
	for _, segment := range s.Segments {
		size := segment.Size()
		if segment.ContainsY(y) {
			// segments that inserting line intersects need to be splitted or expanded
			segment.DeleteRow(y - size.Y)
		} else if size.Y > y {
			// segments laying below inserting line need to be shifted down
			segment.Move(size.X, size.Y-1)
		}
	}
}

func (s *Sheet) DeleteCol(x int) {
	for _, segment := range s.Segments {
		size := segment.Size()
		if segment.ContainsX(x) {
			// segments that inserting column intersects need to be splitted or expanded
			segment.DeleteCol(x - size.X)
		} else if size.X > x {
			// segments laying right of inserting column need to be shifted right
			segment.Move(size.X-1, size.Y)
		}
	}
}

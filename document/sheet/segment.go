package sheet

// Сегмент - это прямоугольник из ячеек. Сегменты могут быть разных типов и не обязательно
// хранят все значения ячеек в явном виде, вместо этого некоторые ячейки могут вычисляться
// при запросе их значения.
//
// Лист может состоять из неограниченного количества сегментов. Сегменты не могут пересекаться,
// то есть одна ячейка не может принадлежать более чем одному сегменту.

type Segment interface {
	Contains(x, y int) bool
	ContainsX(x int) bool
	ContainsY(y int) bool
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

// Проверяет, пересекает ли сегмент колонку x.
func (s *baseSegment) ContainsX(x int) bool {
	return x >= s.size.X && x < s.size.X+s.size.Width
}

// Проверяет, пересекает ли сегмент строку y.
func (s *baseSegment) ContainsY(y int) bool {
	return y >= s.size.Y && y < s.size.Y+s.size.Height
}

// Size returns rect of segment bounds.
func (s *baseSegment) Size() Rect {
	return s.size
}

// Меняет положение левого верхнего угла сегмента на листе.
func (s *baseSegment) Move(newX, newY int) {
	s.size.X = newX
	s.size.Y = newY
}

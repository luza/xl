package document

import (
	"fmt"
	"regexp"
	"strconv"

	"xl/document/sheet"
	"xl/document/value"
)

type Document struct {
	Sheets        []*sheet.Sheet
	CurrentSheet  *sheet.Sheet
	CurrentSheetN int

	maxSheetIdx int

	value.LinkRegistryInterface
	linksRegistry map[int]map[int]map[int]*value.Link
}

var cellNamePattern = regexp.MustCompile("^([A-Z]+)([0-9]+)$")

func New() *Document {
	return &Document{
		linksRegistry: make(map[int]map[int]map[int]*value.Link),
	}
}

func NewWithEmptySheet() *Document {
	s := sheet.New(1, "Sheet 1")
	return &Document{
		Sheets:        []*sheet.Sheet{s},
		CurrentSheet:  s,
		CurrentSheetN: 0,
		maxSheetIdx:   1,
		linksRegistry: make(map[int]map[int]map[int]*value.Link),
	}
}

// The maximum sheet name length is 31 characters. If the sheet name length is exceeded an error is thrown.
// These special characters are also not allowed: : \ / ? * [ ]
func (d *Document) NewSheet(title string) (*sheet.Sheet, error) {
	if title == "" {
		title = fmt.Sprintf("Sheet %d", d.maxSheetIdx+1)
	} else {
		// ensure title is unique
		for _, s := range d.Sheets {
			if s.Title == title {
				return nil, value.NewError(value.ErrorKindName, "duplicating sheet name")
			}
		}
	}
	s := sheet.New(d.maxSheetIdx+1, title)
	d.Sheets = append(d.Sheets, s)
	d.maxSheetIdx++
	return s, nil
}

func (d *Document) InsertRow(n int) {
	d.CurrentSheet.Cursor.Y += n
	d.CurrentSheet.InsertRow(d.CurrentSheet.Cursor.Y)
	// TODO: relinking
}

func (d *Document) InsertCol(n int) {
	d.CurrentSheet.Cursor.X += n
	d.CurrentSheet.InsertCol(d.CurrentSheet.Cursor.X)
	// TODO: relinking
}

func (d *Document) DeleteRow() {
	d.CurrentSheet.DeleteRow(d.CurrentSheet.Cursor.Y)
	// TODO: relinking
}

func (d *Document) DeleteCol() {
	d.CurrentSheet.DeleteCol(d.CurrentSheet.Cursor.X)
	// TODO: relinking
}

func (d *Document) FindCell(cellName string) (int, int, error) {
	// TODO: accept sheet name in request
	return cellNameToXY(cellName)
}

func (d *Document) sheetByIdx(idx int) *sheet.Sheet {
	for _, s := range d.Sheets {
		if s.Idx == idx {
			return s
		}
	}
	return nil
}

func cellNameToXY(name string) (int, int, error) {
	res := cellNamePattern.FindStringSubmatch(name)
	if len(res) < 3 {
		return 0, 0, value.NewError(value.ErrorKindName, "malformed cell name")
	}
	col, row := res[1], res[2]
	x, p := 0, 1
	for c := len(col) - 1; c >= 0; c-- {
		x += int(col[c]-'A'+1) * p
		p *= 26
	}
	y, _ := strconv.Atoi(row)
	if x < 1 || y < 1 {
		return 0, 0, value.NewError(value.ErrorKindName, "malformed cell name")
	}
	return x - 1, y - 1, nil
}

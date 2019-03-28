package document

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"xl/document/sheet"
	"xl/document/value"
)

const (
	maxSheetTitleLength = 31
)

type Document struct {
	Sheets        []*sheet.Sheet
	CurrentSheet  *sheet.Sheet
	CurrentSheetN int

	maxSheetIdx int

	value.LinkRegistryInterface
	linksRegistry map[int]map[int]map[int]*value.Link
}

var cellNamePattern = regexp.MustCompile(`^\$?([A-Z]+)\$?([0-9]+)$`)

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

// NewSheet creates a new sheet in the document. If title is not present, generated one will be used.
func (d *Document) NewSheet(title string) (*sheet.Sheet, error) {
	if title != "" {
		if len(title) > maxSheetTitleLength {
			return nil, value.NewError(value.ErrorKindName, "sheet title must be up to 31 characters long")
		}
		if strings.ContainsAny(title, ":\\/?*[]") {
			return nil, value.NewError(value.ErrorKindName, "sheet title can not include : \\ / ? * [ ]")
		}
		// ensure title is unique
		for _, s := range d.Sheets {
			if s.Title == title {
				return nil, value.NewError(value.ErrorKindName, "duplicating sheet title")
			}
		}
	} else {
		title = fmt.Sprintf("Sheet %d", d.maxSheetIdx+1)
	}
	s := sheet.New(d.maxSheetIdx+1, title)
	d.Sheets = append(d.Sheets, s)
	d.maxSheetIdx++
	return s, nil
}

// InsertEmptyRow inserts new empty row at position of cursor plus N.
func (d *Document) InsertEmptyRow(n int) {
	d.CurrentSheet.Cursor.Y += n
	d.CurrentSheet.InsertEmptyRow(d.CurrentSheet.Cursor.Y)
	// TODO: relinking
}

// InsertEmptyCol inserts new empty column at position of cursor plus N.
func (d *Document) InsertEmptyCol(n int) {
	d.CurrentSheet.Cursor.X += n
	d.CurrentSheet.InsertEmptyCol(d.CurrentSheet.Cursor.X)
	// TODO: relinking
}

// DeleteRow deletes row under cursor.
func (d *Document) DeleteRow() {
	d.CurrentSheet.DeleteRow(d.CurrentSheet.Cursor.Y)
	// TODO: relinking
}

// DeleteCol deletes column under cursor.
func (d *Document) DeleteCol() {
	d.CurrentSheet.DeleteCol(d.CurrentSheet.Cursor.X)
	// TODO: relinking
}

// FindCell finds position of the cell with given name.
func (d *Document) FindCell(cellName string) (int, int, error) {
	// TODO: accept sheet name in request
	return cellNameToXY(cellName)
}

// sheetByIdx returns sheet by its index.
func (d *Document) sheetByIdx(idx int) *sheet.Sheet {
	for _, s := range d.Sheets {
		if s.Idx == idx {
			return s
		}
	}
	return nil
}

// cellNameToXY transforms cell name into X, Y coordinates.
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

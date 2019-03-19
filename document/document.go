package document

import (
	"errors"
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
	s := sheet.New(1, "New Sheet")
	return &Document{
		Sheets:        []*sheet.Sheet{s},
		CurrentSheet:  s,
		CurrentSheetN: 0,
		maxSheetIdx:   1,
		linksRegistry: make(map[int]map[int]map[int]*value.Link),
	}
}

func (d *Document) NewSheet(title string) (*sheet.Sheet, error) {
	// FIXME: title must be unique
	if title == "" {
		title = fmt.Sprintf("Sheet %d", d.maxSheetIdx+1)
	}
	s := sheet.New(d.maxSheetIdx+1, title)
	d.Sheets = append(d.Sheets, s)
	d.maxSheetIdx++
	return s, nil
}

func (d *Document) sheetByIdx(idx int) *sheet.Sheet {
	for _, s := range d.Sheets {
		if s.Idx == idx {
			return s
		}
	}
	return nil
}

//func (d *Document) FindCell(sheetTitle string, x, y int) *sheet.Cell {
//	for _, s := range d.Sheets {
//		if s.Title == sheetTitle {
//			return s.Cell(x, y)
//		}
//	}
//	return nil
//}

func cellNameToXY(name string) (int, int, error) {
	res := cellNamePattern.FindStringSubmatch(name)
	if len(res) < 3 {
		return 0, 0, errors.New("malformed cell name")
	}
	col, row := res[1], res[2]
	x, p := 0, 1
	for c := len(col) - 1; c >= 0; c-- {
		x += int(col[c]-'A'+1) * p
		p *= 26
	}
	y, _ := strconv.Atoi(row)
	if x < 1 || y < 1 {
		return 0, 0, errors.New("malformed cell name")
	}
	return x - 1, y - 1, nil
}

package document

import (
	"xl/document/sheet"
)

type Document struct {
	Sheets        []*sheet.Sheet
	CurrentSheet  *sheet.Sheet
	CurrentSheetN int
}

func New() *Document {
	return &Document{}
}

func NewWithEmptySheet() *Document {
	s := sheet.New("New Sheet")
	return &Document{
		Sheets:        []*sheet.Sheet{s},
		CurrentSheet:  s,
		CurrentSheetN: 0,
	}
}

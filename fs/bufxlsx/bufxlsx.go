package bufxlsx

import (
	"xl/document"
	"xl/document/sheet"
	"xl/fs"

	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type BufXLSX struct {
	fs.FileInterface
	filename string
}

func NewWithFilename(filename string) *BufXLSX {
	return &BufXLSX{
		filename: filename,
	}
}

func (b *BufXLSX) Open() (*document.Document, error) {
	xlsx, err := excelize.OpenFile(b.filename)
	if err != nil {
		return nil, err
	}

	d := document.New()

	for _, name := range xlsx.GetSheetMap() {
		data, err := xlsx.GetRows(name)
		if err != nil {
			return nil, err
		}

		if len(data) == 0 || len(data[0]) == 0 {
			continue
		}

		width, height := len(data[0]), len(data)

		// make cells & transpose
		cells := make([][]sheet.Cell, len(data[0]))
		for x := 0; x < width; x++ {
			cells[x] = make([]sheet.Cell, height)
			for y := 0; y < height; y++ {
				cells[x][y] = *sheet.NewCellUntyped(data[y][x])
			}
		}

		s, err := d.NewSheet(name)
		if err != nil {
			return nil, err
		}

		s.AddStaticSegment(0, 0, width, height, cells)
	}

	return d, nil
}

// Write writes display values for the current sheet into CSV file.
func (b *BufXLSX) Write(doc *document.Document) error {
	file, err := os.Create(b.filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// TODO

	return nil
}

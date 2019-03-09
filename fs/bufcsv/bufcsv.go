package bufcsv

import (
	"xl/document"
	"xl/document/sheet"
	"xl/fs"

	"encoding/csv"
	"fmt"
	"os"
)

type BufCSV struct {
	fs.FileInterface
	filename string
}

func NewWithFilename(filename string) *BufCSV {
	return &BufCSV{
		filename: filename,
	}
}

func (b *BufCSV) Open() (*document.Document, error) {
	file, err := os.Open(b.filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	r := csv.NewReader(file)

	// TODO: how to make it variable? Auto detection?
	r.Comma = ';'

	data, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	// assume all rows have equal length
	height := len(data)
	width := 0
	for l, row := range data {
		if width > 0 && width != len(row) {
			err := fmt.Errorf("failed to read CSV at line %d: unexpected number of columns %d (expected %d)",
				l, len(row), width)
			return nil, err
		}
		width = len(row)
	}

	// make cells & transpose
	cells := make([][]sheet.Cell, width)
	for x := 0; x < width; x++ {
		cells[x] = make([]sheet.Cell, height)
		for y := 0; y < height; y++ {
			cells[x][y] = *sheet.NewCellText(data[y][x])
		}
	}

	// make a sheet
	s := sheet.New(b.filename)
	s.AddStaticSegment(0, 0, width, height, cells)

	d := document.New()
	d.Sheets = []*sheet.Sheet{s}
	return d, nil
}

// Write writes display values for the current sheet into CSV file.
func (b *BufCSV) Write(doc *document.Document) error {
	file, err := os.Create(b.filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	w := csv.NewWriter(file)
	defer w.Flush()

	// TODO: how to make it variable?
	w.Comma = ';'

	s := doc.CurrentSheet
	row := make([]string, s.Size.X+s.Size.Width)

	x, y := 0, 0
	for y < s.Size.Height+s.Size.Y {
		// write empty rows at top of sheet
		if y < s.Size.Y {
			if err = w.Write(row); err != nil {
				return err
			}
			y++
			continue
		}
		for x < s.Size.Width+s.Size.X {
			// fill empty cells at begging of row
			if x < s.Size.X {
				row[x] = ""
				x++
				continue
			}
			// iterate over segments
			segment := s.FindSegment(x, y)
			if segment == nil {
				row[x] = ""
				x++
				continue
			}
			size := segment.Size()
			// copy cells from found segment into row
			for x < size.Right() {
				row[x] = segment.Cell(x, y).DisplayText()
				x++
			}
		}
		if err = w.Write(row); err != nil {
			return err
		}
		y++
		x = 0
	}

	return nil
}

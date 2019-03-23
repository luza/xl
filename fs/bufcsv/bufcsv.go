package bufcsv

import (
	"xl/document"
	"xl/document/sheet"
	"xl/document/value"
	"xl/fs"

	"encoding/csv"
	"fmt"
	"os"
)

type BufCSV struct {
	fs.FileInterface
	filename string
	comma    rune
}

func NewWithFilename(filename string) *BufCSV {
	return &BufCSV{
		filename: filename,
		comma:    ',',
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
	b.comma = r.Comma

	// FIXME: read line-by-line
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

	d := document.New()

	// make cells & transpose
	cells := make([][]sheet.Cell, width)
	for x := 0; x < width; x++ {
		cells[x] = make([]sheet.Cell, height)
		for y := 0; y < height; y++ {
			cells[x][y] = *sheet.NewCellUntyped(data[y][x])
		}
	}

	// make a sheet
	s, _ := d.NewSheet(b.filename)
	s.AddStaticSegment(0, 0, width, height, cells)
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

	// the same as on reading
	w.Comma = b.comma

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
			for x <= size.MaxX() {
				row[x], _ = segment.Cell(x, y).StringValue(value.NewEvalContext(doc))
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

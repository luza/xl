package document

import (
	"xl/document/eval"
	"xl/document/sheet"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCellAxis(t *testing.T) {
	testCases := []struct {
		name      string
		x         int
		y         int
		anchoredX bool
		anchoredY bool
	}{
		{`A1`, 0, 0, false, false},
		{`B1`, 1, 0, false, false},
		{`B2`, 1, 1, false, false},
		{`Z1`, 25, 0, false, false},
		{`AA1`, 26, 0, false, false},
		{`AB1`, 27, 0, false, false},
		{`AZ1`, 51, 0, false, false},
		{`BA1`, 52, 0, false, false},
		{`AAA999`, 702, 998, false, false},
		{`AAAA1`, 18278, 0, false, false},
		{`$A1`, 0, 0, true, false},
		{`A$1`, 0, 0, false, true},
		{`$A$1`, 0, 0, true, true},
	}
	for _, c := range testCases {
		x, y, anchoredX, anchoredY, err := CellAxis(c.name)
		assert.NoErrorf(t, err, "case %s: must not fail on parse %s", c.name, err)
		assert.Equalf(t, c.x, x, "case %s: must be true X: %d==%d", c.name, c.x, x)
		assert.Equalf(t, c.y, y, "case %s: must be true Y: %d==%d", c.name, c.y, y)
		assert.Equalf(t, c.anchoredX, anchoredX, "case %s, anchoredX must match")
		assert.Equalf(t, c.anchoredY, anchoredY, "case %s, anchoredY must match")
	}
}

func TestColName(t *testing.T) {
	testCases := []struct {
		n    int
		name string
	}{
		{0, "A"},
		{25, "Z"},
		{26, "AA"},
		{51, "AZ"},
		{52, "BA"},
		{700, "ZY"},
		{701, "ZZ"},
		{702, "AAA"},
		{18277, "ZZZ"},
		{18278, "AAAA"},
	}
	for _, c := range testCases {
		name := ColName(c.n)
		assert.Equalf(t, c.name, name, "case %d: must be equal", c.n)
	}
}

func TestCellStringValue(t *testing.T) {
	d := NewWithEmptySheet()
	d.CurrentSheet.SetCell(0, 0, sheet.NewCellUntyped("1"))

	v, err := d.CurrentSheet.Cell(0, 0).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.NoError(t, err)
	assert.Equal(t, "1", v)
}

func TestCellSimpleArithmetic(t *testing.T) {
	d := NewWithEmptySheet()
	d.CurrentSheet.SetCell(0, 0, sheet.NewCellUntyped("1"))
	d.CurrentSheet.SetCell(0, 1, sheet.NewCellUntyped("=A1+1"))

	v, err := d.CurrentSheet.Cell(0, 1).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.NoError(t, err)
	assert.Equal(t, "2", v)
}

func TestCellSimpleCircularReferencing(t *testing.T) {
	d := NewWithEmptySheet()
	d.CurrentSheet.SetCell(0, 0, sheet.NewCellUntyped("=A1"))

	_, err := d.CurrentSheet.Cell(0, 0).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.EqualError(t, err, "circular reference")
}

func TestCellSimpleCircularReferencingWithArithmetic(t *testing.T) {
	d := NewWithEmptySheet()
	d.CurrentSheet.SetCell(0, 0, sheet.NewCellUntyped("=A1+1"))

	_, err := d.CurrentSheet.Cell(0, 0).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.EqualError(t, err, "circular reference")
}

func TestCellSimpleCircularReferencingWithArithmetic2(t *testing.T) {
	d := NewWithEmptySheet()
	d.CurrentSheet.SetCell(0, 0, sheet.NewCellUntyped("=1+A1"))

	_, err := d.CurrentSheet.Cell(0, 0).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.EqualError(t, err, "circular reference")
}

func TestCellMultipleReferencingCellInArithmetic(t *testing.T) {
	d := NewWithEmptySheet()
	d.CurrentSheet.SetCell(0, 0, sheet.NewCellUntyped("1"))
	d.CurrentSheet.SetCell(0, 1, sheet.NewCellUntyped("=A1+A1"))

	v, err := d.CurrentSheet.Cell(0, 1).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.NoError(t, err)
	assert.Equal(t, "2", v)
}

func TestCellEmptyCellArithmetic(t *testing.T) {
	d := NewWithEmptySheet()

	s2, err := d.NewSheet("sh2")
	assert.NoError(t, err)

	d.CurrentSheet.SetCell(0, 0, sheet.NewCellUntyped("='sh2'!A1+1"))
	s2.SetCell(0, 0, sheet.NewCellUntyped(""))

	v, err := d.CurrentSheet.Cell(0, 0).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.NoError(t, err)
	assert.Equal(t, "1", v)
}

func TestCellRangeRef(t *testing.T) {
	d := NewWithEmptySheet()
	d.CurrentSheet.AddStaticSegment(0, 0, 2, 2, [][]sheet.Cell{
		{*sheet.NewCellUntyped("1"), *sheet.NewCellUntyped("3")},
		{*sheet.NewCellUntyped("2"), *sheet.NewCellUntyped("4")},
	})
	d.CurrentSheet.SetCell(0, 2, sheet.NewCellUntyped("=SUM(A1:B2)"))

	v, err := d.CurrentSheet.Cell(0, 2).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.NoError(t, err)
	assert.Equal(t, "10", v)
}

func TestCellRangeRefWithInvalidRange(t *testing.T) {
	d := NewWithEmptySheet()
	d.CurrentSheet.SetCell(0, 0, sheet.NewCellUntyped("=SUM(C2:B1)"))

	_, err := d.CurrentSheet.Cell(0, 0).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.EqualError(t, err, "invalid range bounds")
}

func TestCellRangeRefWithCrossSheetsRange(t *testing.T) {
	d := NewWithEmptySheet()
	d.NewSheet("")
	d.CurrentSheet.SetCell(0, 0, sheet.NewCellUntyped("=SUM(Sheet1!A2:Sheet2!B3)"))

	_, err := d.CurrentSheet.Cell(0, 0).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.EqualError(t, err, "cross-sheets ranges are not allowed")
}

func TestCellRangeRefWithCircularRef(t *testing.T) {
	d := NewWithEmptySheet()
	d.CurrentSheet.AddStaticSegment(0, 0, 2, 2, [][]sheet.Cell{
		{*sheet.NewCellUntyped("1"), *sheet.NewCellUntyped("3")},
		{*sheet.NewCellUntyped("2"), *sheet.NewCellUntyped("4")},
	})
	d.CurrentSheet.SetCell(0, 2, sheet.NewCellUntyped("=SUM(A1:B3)"))

	_, err := d.CurrentSheet.Cell(0, 2).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.EqualError(t, err, "circular reference")
}

func TestCellRefWithAnchors(t *testing.T) {
	d := NewWithEmptySheet()
	d.CurrentSheet.SetCell(0, 0, sheet.NewCellUntyped("1"))
	d.CurrentSheet.SetCell(0, 1, sheet.NewCellUntyped("=$A$1"))

	v := d.CurrentSheet.Cell(0, 1).RawValue()
	assert.Equal(t, "=$A$1", v)

	v = d.CurrentSheet.Cell(0, 1).Expression(eval.NewContext(d, d.CurrentSheet.Idx)).String()
	assert.Equal(t, "=$A$1", v)

	v, err := d.CurrentSheet.Cell(0, 1).StringValue(eval.NewContext(d, d.CurrentSheet.Idx))
	assert.NoError(t, err)
	assert.Equal(t, "1", v)
}

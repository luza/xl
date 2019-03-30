package document

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCellNameToXY(t *testing.T) {
	testCases := []struct {
		name string
		x    int
		y    int
	}{
		{`A1`, 0, 0},
		{`B1`, 1, 0},
		{`B2`, 1, 1},
		{`Z1`, 25, 0},
		{`AA1`, 26, 0},
		{`AB1`, 27, 0},
		{`AZ1`, 51, 0},
		{`BA1`, 52, 0},
		{`AAA999`, 702, 998},
		{`AAAA1`, 18278, 0},
	}
	for _, c := range testCases {
		x, y, err := CellAxis(c.name)
		assert.NoErrorf(t, err, "case %s: must not fail on parse %s", c.name, err)
		assert.Equalf(t, c.x, x, "case %s: must be true X: %d==%d", c.name, c.x, x)
		assert.Equalf(t, c.y, y, "case %s: must be true Y: %d==%d", c.name, c.y, y)
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

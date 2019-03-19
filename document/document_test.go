package document

import "testing"

func TestCellNameToXY(t *testing.T) {
	testCases := []struct {
		name string
		x    int
		y    int
	}{
		{`A1`, 1, 1},
		{`B1`, 2, 1},
		{`B2`, 2, 2},
		{`Z1`, 26, 1},
		{`AA1`, 27, 1},
		{`AB1`, 28, 1},
		{`AAA999`, 703, 999},
	}
	for _, c := range testCases {
		x, y, err := cellNameToXY(c.name)
		if err != nil {
			t.Errorf("case %s: must not fail on parse %s", c.name, err)
			continue
		}
		if x != c.x || y != c.y {
			t.Errorf("case %s: must be true X: %d==%d, Y: %d==%d", c.name, c.x, x, c.y, y)
		}
	}
}

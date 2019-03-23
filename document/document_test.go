package document

import "testing"

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
		{`AAA999`, 702, 998},
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

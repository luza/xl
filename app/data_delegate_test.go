package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
		name := colName(c.n)
		assert.Equalf(t, c.name, name, "case %d: must be equal", c.n)
	}
}

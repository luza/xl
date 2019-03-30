package sheet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGuessCellType(t *testing.T) {
	testCases := []struct {
		value       string
		t           int
		castedValue interface{}
	}{
		{``, CellValueTypeEmpty, nil},
		{`TRUE`, CellValueTypeBool, true},
		{`false`, CellValueTypeBool, false},
		{`0`, CellValueTypeInteger, 0},
		{`1`, CellValueTypeInteger, 1},
		{`0.1`, CellValueTypeDecimal, nil},
		{`1.0`, CellValueTypeDecimal, nil},
		{`1.0e10`, CellValueTypeDecimal, nil},
		{`=FUNC()`, CellValueTypeFormula, nil},
		{`=`, CellValueTypeText, "="},
		{`abc`, CellValueTypeText, "abc"},
	}
	for _, c := range testCases {
		guessedType, castedValue := guessCellType(c.value)
		assert.Equalf(t, c.t, guessedType, "case %s", c.value)
		assert.Equalf(t, c.castedValue, castedValue, "case %s", c.value)
	}
}

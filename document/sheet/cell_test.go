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
		{``, cellValueTypeEmpty, nil},
		{`TRUE`, cellValueTypeBool, true},
		{`false`, cellValueTypeBool, false},
		{`0`, cellValueTypeInteger, 0},
		{`1`, cellValueTypeInteger, 1},
		{`0.1`, cellValueTypeDecimal, nil},
		{`1.0`, cellValueTypeDecimal, nil},
		{`1.0e10`, cellValueTypeDecimal, nil},
		{`=FUNC()`, cellValueTypeFormula, nil},
		{`=`, cellValueTypeString, "="},
		{`abc`, cellValueTypeString, "abc"},
	}
	for _, c := range testCases {
		guessedType, castedValue := guessCellType(c.value)
		assert.Equalf(t, c.t, guessedType, "case %s", c.value)
		assert.Equalf(t, c.castedValue, castedValue, "case %s", c.value)
	}
}

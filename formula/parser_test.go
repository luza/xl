package formula

import (
	"testing"

	"xl/document/value"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	testCases := []struct {
		f       string
		res     string
		varsNum int
	}{
		{`=1`, "1", 0},
		{`=(1)`, "1", 0},
		{`=1+1`, "2", 0},
		{`=(1+1)`, "2", 0},
		{`=2+2*2`, "6", 0},
		{`=1+1*1+1`, "3", 0},
		{`="string"`, "string", 0},
		{`="ap""""ple"`, `ap""ple`, 0},
		{`=tRUE`, "TRUE", 0},
		{`=TRUE+TRUE`, "2", 0},
		{`=-TRUE`, "-1", 0},
		{`=+TRUE`, "TRUE", 0},
		{`=1=1`, "TRUE", 0},
		{`=1<>1`, "FALSE", 0},
		{`=1>1`, "FALSE", 0},
		{`=1<1`, "FALSE", 0},
		{`=1>=1`, "TRUE", 0},
		{`=1<=1`, "TRUE", 0},
		{`=-1`, "-1", 0},
		{`=+1`, "1", 0},
		{`=--1`, "1", 0},
		{`=TRIM("  ggg  ")`, "ggg", 0},
		{`=SUM(1)`, "1", 0},
		{`=SUM(1, 2, 3)`, "6", 0},
		{`=A1`, "4", 1},
		{`=-A1`, "-4", 1},
		{`=A1+A1`, "10", 2},
		{`=Sheet!A1+Sheet2!A1`, "10", 2},
		{`='Sheet With Spaces'!A1+Sheet2!A1`, "10", 2},
		{`='Sheet ''With Spaces'!A1`, "4", 1},
		{`=A1:B200+A1:C300`, "10", 2},
		{`='Sheet With Spaces'!A1:'Sheet With Spaces'!B200+Sheet2!A1:Sheet2!C300`, "10", 2},
	}
	for _, c := range testCases {
		f, vb, err := Parse(c.f)
		assert.NoErrorf(t, err, "case %s: must not fail on parse %s", c.f, err)
		assert.Lenf(t, vb.Vars, c.varsNum, "case %s: must return %d variables (returned %d)", c.f, c.varsNum, len(vb.Vars))
		var lr value.LinkRegistryInterface
		ec := value.NewEvalContext(lr)
		v, err := f(ec, []value.Value{
			value.NewDecimalValue(decimal.NewFromFloat(4)),
			value.NewDecimalValue(decimal.NewFromFloat(6)),
		})
		assert.NoErrorf(t, err, "case %s: function must not fail, got %s", c.f, err)
		s, _ := v.StringValue(ec)
		assert.Equalf(t, c.res, s, "case %s: must be equal to %s, got %s", c.f, c.res, s)
	}
}

func TestParseErrors(t *testing.T) {
	testCases := []struct {
		f   string
		err string
	}{
		{`=`, `<source>:1:1: unexpected "<EOF>" (expected <cell>)`},
		{`=()`, `<source>:1:3: unexpected ")" (expected <cell>)`},
		{`=1+`, `<source>:1:3: unexpected token "+"`},
	}
	for _, c := range testCases {
		_, _, err := Parse(c.f)
		assert.Errorf(t, err, "case %s: must fail", c.f)
		assert.Equalf(t, c.err, err.Error(), "case %s: must fail with reason '%s', actual '%s'", c.f, c.err, err.Error())
	}
}

func TestExecuteErrors(t *testing.T) {
	testCases := []struct {
		f   string
		err string
	}{
		{`="a"+"b"`, `arithmetic (+) on string operand`},
		{`=1/0`, `division by zero`},
	}
	for _, c := range testCases {
		f, _, err := Parse(c.f)
		assert.NoErrorf(t, err, "case %s: must not fail on parse %s", c.f, err)
		var lr value.LinkRegistryInterface
		ec := value.NewEvalContext(lr)
		_, err = f(ec, []value.Value{
			value.NewDecimalValue(decimal.NewFromFloat(4)),
			value.NewDecimalValue(decimal.NewFromFloat(6)),
		})
		assert.Errorf(t, err, "case %s: execution must fail", c.f)
		assert.Equalf(t, c.err, err.Error(), "case %s: must fail with reason '%s', actual '%s'", c.f, c.err, err.Error())
	}
}
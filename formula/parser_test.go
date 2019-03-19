package formula

import (
	"testing"
	"xl/formula/value"

	"github.com/shopspring/decimal"
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
		if err != nil {
			t.Errorf("case %s: must not fail on parse %s", c.f, err)
			continue
		}
		if len(vb.Vars) != c.varsNum {
			t.Errorf("case %s: must return %d variables (returned %d)", c.f, c.varsNum, len(vb.Vars))
		}
		v, err := f([]value.Value{
			value.NewDecimalValue(decimal.NewFromFloat(4)),
			value.NewDecimalValue(decimal.NewFromFloat(6)),
		})
		if err != nil {
			t.Errorf("case %s: function must not fail, got %s", c.f, err)
		}
		s, _ := v.StringValue()
		if s != c.res {
			t.Errorf("case %s: must be equal to %s, got %s", c.f, c.res, s)
		}
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
		if err == nil {
			t.Errorf("case %s: must fail", c.f)
			continue
		}
		if err.Error() != c.err {
			t.Errorf("case %s: must fail with reason '%s', actual '%s'", c.f, c.err, err.Error())
		}
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
		if err != nil {
			t.Errorf("case %s: must not fail on parse %s", c.f, err)
			continue
		}
		_, err = f([]value.Value{
			value.NewDecimalValue(decimal.NewFromFloat(4)),
			value.NewDecimalValue(decimal.NewFromFloat(6)),
		})
		if err == nil {
			t.Errorf("case %s: execution must fail", c.f)
			continue
		}
		if err.Error() != c.err {
			t.Errorf("case %s: must fail with reason '%s', actual '%s'", c.f, c.err, err.Error())
		}
	}
}

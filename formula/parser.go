package formula

import (
	"xl/document/eval"

	"bytes"
	"strings"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
)

type Boolean bool
type Sheet string
type FuncName string
type String string

func (b *Boolean) Capture(values []string) error {
	*b = Boolean(strings.EqualFold(values[0], "TRUE"))
	return nil
}

func (s *Sheet) Capture(values []string) error {
	values[0] = strings.TrimRight(values[0], "!")
	// remove first and last char
	l := len(values[0])
	if l >= 2 && values[0][0] == '\'' && values[0][l-1] == '\'' {
		values[0] = values[0][1 : len(values[0])-1]
	}
	// replace double quotes with single quote
	values[0] = strings.Replace(values[0], `''`, `'`, -1)
	*s = Sheet(values[0])
	return nil
}

func (f *FuncName) Capture(values []string) error {
	*f = FuncName(strings.TrimRight(values[0], "("))
	return nil
}

func (s *String) Capture(values []string) error {
	// remove first and last char
	values[0] = values[0][1 : len(values[0])-1]
	// replace double quotes with single quote
	values[0] = strings.Replace(values[0], `""`, `"`, -1)
	*s = String(values[0])
	return nil
}

type Expression struct {
	Equality *Equality `"=" @@`
}

func (e *Expression) String() string {
	var buf bytes.Buffer
	e.Output(func(s string, i int) {
		buf.WriteString(s)
	})
	return buf.String()
}

type Equality struct {
	Comparison *Comparison `@@`
	Op         string      `[ @( "<>" | "=" )`
	Next       *Equality   `  @@ ]`
}

type Comparison struct {
	Addition *Addition   `@@`
	Op       string      `[ @( ">" | ">=" | "<" | "<=" )`
	Next     *Comparison `  @@ ]`
}

type Addition struct {
	Multiplication *Multiplication `@@`
	Op             string          `[ @( "-" | "+" )`
	Next           *Addition       `  @@ ]`
}

type Multiplication struct {
	Power *Power          `@@`
	Op    string          `[ @( "/" | "*" )`
	Next  *Multiplication `  @@ ]`
}

type Power struct {
	Base     *Unary   `@@`
	Exponent []*Unary `{ "^" @@ }`
}

type Unary struct {
	Op      string   `( @( "-" | "+" )`
	Unary   *Unary   `  @@ )`
	Primary *Primary `| @@`
}

type Primary struct {
	SubExpression *Equality `"(" @@ ")" `
	Number        *float64  `| @Number`
	String        *String   `| @String`
	Boolean       *Boolean  `| @("TRUE" | "FALSE")`
	Func          *Func     `| @@`
	Variable      *Variable `| @@`
}

type Func struct {
	Name      FuncName    `@FuncName`
	Arguments []*Equality `[ @@ { ";" @@ } ] ")"`
}

type Variable struct {
	Cell   *Cell `@@`
	CellTo *Cell `[ ":" @@ ]`
}

type Cell struct {
	Sheet    *Sheet `[ @Sheet ]`
	CellName string `@CellName`
}

var lex = lexer.Must(lexer.Regexp(
	`(\s+)` +
		`|^=` +
		`|(?P<Operators><>|<=|>=|[-+*/()=<>;:\^])` +
		`|(?P<Number>\d*\.?\d+([eE][-+]?\d+)?)` +
		`|(?P<String>"([^"]|"")*")` +
		`|(?P<Boolean>(?i)TRUE|FALSE)` +
		`|(?P<FuncName>[A-Za-z][A-Za-z0-9\.]+)\(` +
		`|(?P<Sheet>[A-Za-z0-9_]+|'([^']|'')*')!` +
		`|(?P<CellName>\$?[A-Za-z]+\$?[1-9][0-9]*)`,
))

// Parse parses the formula, extracts variables from it and builds
// functions chain that perform the expression representing by the formula..
func Parse(source string) (*Expression, error) {
	// TODO: do that once
	p, err := participle.Build(
		&Expression{},
		participle.Lexer(lex),
		participle.CaseInsensitive("Boolean"),
		participle.Upper("CellName"),
	)
	if err != nil {
		panic(err)
	}
	expression := &Expression{}
	if err = p.ParseString(source, expression); err != nil {
		return nil, eval.NewError(eval.ErrorKindFormula, err.Error())
	}
	return expression, nil
}

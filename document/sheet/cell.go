package sheet

import (
	"strconv"

	"xl/document/eval"
	"xl/formula"

	"github.com/shopspring/decimal"
)

// Ячейка - это структура, описывающая одну клетку листа.
// В ячейке хранится сырое значение, введенное пользователем или прочитанное из
// файла, а также преобразованное значение разного типа в зависимости от типа ячейки.
//
// В начальном состоянии ячейка имеет неопределнный тип и содержит только сырое значение.
// При первом обращении на основе сырого значения определяется тип ячейки и внутренняя структура
// преобразовывается в одну из структур, соответствующую определенному типу.
// В новую структуру сохраняется значение, сконвертированное в определенный тип.
//
// Если значение ячейки - формула, то формула парсится и хранится в ячейке в виде Выражения (Expression),
// функции, которая производит вычисления над Значениями в соответствии с заданной формулой, и
// ссылок, которые хранят зависимость этой ячейки от значений других ячеек, если они заданы в формуле.
// При запросе значения ячейки, если она задана формулой, в которой есть ссылки на другие значения,
// это значение вычисляется; для этого запрашивается (и вычисляется, если требуется) значение всех
// ссылаемых ячеек и выполняются арифметические операции, заданные формулой.

// Alloc = 1799 MiB, TotalAlloc = 2241 MiB, Sys = 2326 MiB, NumGC = 13
// Alloc = 1026 MiB, TotalAlloc = 1132 MiB, Sys = 1136 MiB, NumGC = 11

const (
	cellValueTypeEmpty = iota
	cellValueTypeString
	cellValueTypeInteger
	cellValueTypeDecimal
	cellValueTypeBool
	cellValueTypeFormula
)

type Cell struct {
	rawValue string
	// Внутренняя структура ячейки. Разная для разных типов.
	v interface{}
}

type untypedCell struct{}

type stringCell struct{}

type boolCell struct {
	Value bool
}

type intCell struct {
	Value int
}

type decimalCell struct {
	Value decimal.Decimal
}

type formulaCell struct {
	FormulaValue formula.Function
	Expression   *formula.Expression
	Refs         []ref
	offsetX      int
	offsetY      int
}

// Создает пустую ячейку.
func NewCellEmpty() *Cell {
	return &Cell{}
}

// Создает ячейку без определенного типа, фактический тип будет определен позднее.
func NewCellUntyped(v string) *Cell {
	return &Cell{
		rawValue: v,
		v:        untypedCell{},
	}
}

// Копирует значение ячейки и задает полученной копии смещение.
// Смещение используется при разрешении Ссылок, чтобы сдвинуть их относительно
// ключевой ячейки х-сегмента.
func NewCellAsCopyWithOffset(sourceCell *Cell, offsetX, offsetY int) *Cell {
	if v, ok := sourceCell.v.(formulaCell); ok {
		return &Cell{
			// make a copy
			v: formulaCell{
				FormulaValue: v.FormulaValue,
				Expression:   v.Expression,
				Refs:         v.Refs,
				offsetX:      offsetX,
				offsetY:      offsetY,
			},
		}
	} else {
		return sourceCell
	}
}

// RawValue returns raw cell value as string. No evaluation performed.
func (c *Cell) RawValue() string {
	return c.rawValue
}

// Возвращает выражение, построееное по формуле.
// Если в формуле есть Переменные, то они обновляются по актуальным значениям Ссылок.
func (c *Cell) Expression(ec *eval.Context) *formula.Expression {
	switch v := c.v.(type) {
	case untypedCell:
		if err := c.evaluateType(ec); err != nil {
			return nil
		}
		return c.Expression(ec)
	case formulaCell:
		if err := updateVars(ec, v.Expression, v.Refs, v.offsetX, v.offsetY); err != nil {
			return nil
		}
		return v.Expression
	default:
		return nil
	}
}

// BoolValue returns evaluated cell rawValue as boolean.
func (c *Cell) BoolValue(ec *eval.Context) (bool, error) {
	if c.v == nil {
		return false, nil
	}
	switch v := c.v.(type) {
	case untypedCell:
		if err := c.evaluateType(ec); err != nil {
			return false, err
		}
		return c.BoolValue(ec)
	case stringCell:
		return false, eval.NewError(eval.ErrorKindCasting, "unable to cast string to bool")
	case boolCell:
		return v.Value, nil
	case intCell:
		return v.Value != 0, nil
	case decimalCell:
		return !v.Value.Equal(decimal.Zero), nil
	case formulaCell:
		val, err := v.FormulaValue(ec, refsToValues(v.Refs, v.offsetX, v.offsetY))
		if err != nil {
			return false, err
		}
		return val.BoolValue(ec)
	default:
		panic("unsupported type")
	}
}

// DecimalValue returns evaluated cell value as decimal.
func (c *Cell) DecimalValue(ec *eval.Context) (decimal.Decimal, error) {
	if c.v == nil {
		return decimal.Zero, nil
	}
	switch v := c.v.(type) {
	case untypedCell:
		if err := c.evaluateType(ec); err != nil {
			return decimal.Zero, err
		}
		return c.DecimalValue(ec)
	case stringCell:
		return decimal.Zero, eval.NewError(eval.ErrorKindCasting, "unable to cast string to decimal")
	case boolCell:
		return decimal.Zero, eval.NewError(eval.ErrorKindCasting, "unable to cast bool to decimal")
	case intCell:
		return decimal.New(int64(v.Value), 0), nil
	case decimalCell:
		return v.Value, nil
	case formulaCell:
		val, err := v.FormulaValue(ec, refsToValues(v.Refs, v.offsetX, v.offsetY))
		if err != nil {
			return decimal.Zero, err
		}
		return val.DecimalValue(ec)
	default:
		panic("unsupported type")
	}
}

// StringValue returns evaluated cell rawValue as string.
func (c *Cell) StringValue(ec *eval.Context) (string, error) {
	if c.v == nil {
		return "", nil
	}
	switch v := c.v.(type) {
	case untypedCell:
		if err := c.evaluateType(ec); err != nil {
			return "", err
		}
		return c.StringValue(ec)
	case stringCell:
		return c.rawValue, nil
	case boolCell:
		return c.rawValue, nil
	case intCell:
		return c.rawValue, nil
	case decimalCell:
		return c.rawValue, nil
	case formulaCell:
		val, err := v.FormulaValue(ec, refsToValues(v.Refs, v.offsetX, v.offsetY))
		if err != nil {
			return "", err
		}
		return val.StringValue(ec)
	default:
		panic("unsupported type")
	}
}

// Возвращает значение ячейки как Значение для формулы.
func (c *Cell) Value(ec *eval.Context) (eval.Value, error) {
	if c.v == nil {
		return eval.NewEmptyValue(), nil
	}
	switch v := c.v.(type) {
	case untypedCell:
		if err := c.evaluateType(ec); err != nil {
			return eval.NewEmptyValue(), err
		}
		return c.Value(ec)
	case stringCell:
		return eval.NewStringValue(c.rawValue), nil
	case boolCell:
		return eval.NewBoolValue(v.Value), nil
	case intCell:
		return eval.NewDecimalValue(decimal.New(int64(v.Value), 0)), nil
	case decimalCell:
		return eval.NewDecimalValue(v.Value), nil
	case formulaCell:
		return v.FormulaValue(ec, refsToValues(v.Refs, v.offsetX, v.offsetY))
	default:
		panic("unsupported type")
	}
}

// Сбрасывает значение ячейки на пустое.
// FIXME: оставляет осиротевшие Ссылки в refRegistry.
func (c *Cell) SetValueEmpty() {
	c.rawValue = ""
	c.v = nil
}

// SetValueUntyped fill new cell value with no any type associated with it.
// Type will be determined later on demand.
func (c *Cell) SetValueUntyped(v string) {
	c.rawValue = v
	c.v = untypedCell{}
}

// Вычисляет тип ячейки на осное ее сырого значение и крнвертирует внутреннюю структуру в нужный тип.
func (c *Cell) evaluateType(ec *eval.Context) error {
	t, castedV := guessCellType(c.rawValue)
	switch t {
	case cellValueTypeEmpty:
		c.v = nil
	case cellValueTypeString:
		c.v = stringCell{}
	case cellValueTypeInteger:
		c.v = intCell{
			Value: castedV.(int),
		}
	case cellValueTypeDecimal:
		d, _ := decimal.NewFromString(c.rawValue)
		c.v = decimalCell{
			Value: d,
		}
	case cellValueTypeBool:
		c.v = boolCell{
			Value: castedV.(bool),
		}
	case cellValueTypeFormula:
		expr, err := formula.Parse(c.rawValue)
		if err != nil {
			return err
		}
		c.rawValue = expr.String() // need this?
		formulaValue, _ := expr.BuildFunc()
		refs, err := makeRefs(ec, expr.Variables())
		if err != nil {
			return err
		}
		c.v = formulaCell{
			FormulaValue: formulaValue,
			Expression:   expr,
			Refs:         refs,
		}
	default:
		panic("unsupported type")
	}
	return nil
}

// guessCellType detects cell type based on its rawValue.
// Returns detected type and either casted rawValue or nil if casting wasn't done.
func guessCellType(v string) (int, interface{}) {
	if len(v) == 0 {
		return cellValueTypeEmpty, nil
	} else if v[0] == '=' && len(v) > 1 {
		return cellValueTypeFormula, nil
	} else {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return cellValueTypeInteger, int(i)
		}
		if _, err := strconv.ParseFloat(v, 64); err == nil {
			return cellValueTypeDecimal, nil
		}
		if b, err := strconv.ParseBool(v); err == nil {
			return cellValueTypeBool, b
		}
	}
	return cellValueTypeString, v
}

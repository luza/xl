package eval

import (
	"github.com/shopspring/decimal"
)

// Хранит адрес ячейки.
type CellAddress struct {
	SheetIdx int
	X        int
	Y        int
}

// Ссылка хранит адрес ячейки, на которую ссылается, и параметры.
type CellReference struct {
	CellAddress
	// Флаг, заданы ли координаты ячейки как неизменяемые при экстраполяции значений ячеек.
	AnchoredX bool
	AnchoredY bool
}

type RefRegistryInterface interface {
	// Работа с Сылками.
	AddRef(cell CellReference)
	FromAddress(cell CellReference) (string, string, error)
	ToAddress(sheetTitle, cellName string) (CellReference, error)

	// Получение значений по адресу ячейки.
	Value(ec *Context, cell CellAddress) (Value, error)
	BoolValue(ec *Context, cell CellAddress) (bool, error)
	DecimalValue(ec *Context, cell CellAddress) (decimal.Decimal, error)
	StringValue(ec *Context, cell CellAddress) (string, error)

	// Получение значений для диапазона ячеек.
	IterateBoolValues(ec *Context, cell, cellTo CellAddress, f func(bool) error) error
	IterateDecimalValues(ec *Context, cell, cellTo CellAddress, f func(decimal.Decimal) error) error
	IterateStringValues(ec *Context, cell, cellTo CellAddress, f func(string) error) error
}

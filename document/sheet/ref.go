package sheet

import (
	"xl/document/eval"
	"xl/formula"
)

// Ссылка на другую ячейку или диапазон ячеек.
// Если задана только Cell, то это ссылка на ячейку.
// Если заданы оба Cell и CellTo, то это ссылка на диапазон, где Cell - левый верхний угол,
// CellTo - правый нижний.
type ref struct {
	Cell   eval.CellReference
	CellTo *eval.CellReference
}

// Преобразовывает распарсенное лист!имя ячейки из формулы в ее адрес.
func toAddress(ec *eval.Context, c *formula.Cell) (eval.CellReference, error) {
	var sheetTitle string
	if c.Sheet != nil {
		sheetTitle = string(*c.Sheet)
	}
	return ec.DataProvider.ToAddress(sheetTitle, c.CellName)
}

// Преобразовывает адрес ячейки обратно в ее лист!имя, из которго можно составить формулу.
func fromAddress(ec *eval.Context, ca eval.CellReference) (*formula.Cell, error) {
	sheetTitle, cellName, err := ec.DataProvider.FromAddress(ca)
	if err != nil {
		return nil, err
	}
	var s *formula.Sheet
	if sheetTitle != "" {
		fs := formula.Sheet(sheetTitle)
		s = &fs
	}
	c := &formula.Cell{
		Sheet:    s,
		CellName: cellName,
	}
	return c, nil
}

// Делает массив ссылок на основе Переменных из распарсенной формулы.
func makeRefs(ec *eval.Context, vars []*formula.Variable) ([]ref, error) {
	refs := make([]ref, len(vars))
	for i, v := range vars {
		ca, err := toAddress(ec, v.Cell)
		if err != nil {
			return nil, err
		}
		ec.DataProvider.AddRef(ca)
		ref := ref{
			Cell: ca,
		}
		if v.CellTo != nil {
			ca, err = toAddress(ec, v.CellTo)
			if err != nil {
				return nil, err
			}
			ec.DataProvider.AddRef(ca)
			ref.CellTo = &ca
		}
		refs[i] = ref
	}
	return refs, nil
}

// Делает массив Значений из массива Ссылок.
func refsToValues(refs []ref, offsetX, offsetY int) []eval.Value {
	values := make([]eval.Value, len(refs))
	for i, r := range refs {
		cell := applyOffsetToRef(r.Cell, offsetX, offsetY)
		var cellTo *eval.CellReference
		if r.CellTo != nil {
			c := applyOffsetToRef(*r.CellTo, offsetX, offsetY)
			cellTo = &c

		}
		values[i] = eval.NewRefValue(cell, cellTo)
	}
	return values
}

// Изменяет Переменные в Выражении так, чтобы они получили значение в соответствии с акутальным
// значением ссылки.
func updateVars(ec *eval.Context, x *formula.Expression, refs []ref, offsetX, offsetY int) error {
	for i, v := range x.Variables() {
		cell := applyOffsetToRef(refs[i].Cell, offsetX, offsetY)
		c, err := fromAddress(ec, cell)
		if err != nil {
			return err
		}
		v.Cell = c
		if refs[i].CellTo != nil {
			cellTo := applyOffsetToRef(*refs[i].CellTo, offsetX, offsetY)
			c, err = fromAddress(ec, cellTo)
			if err != nil {
				return err
			}
			v.CellTo = c
		}
	}
	return nil
}

// Применяет смещение к ссылке. Смещение использутеся в экстраполяционном сегменте для
// указания, насколько запрашиваемая ячейка отстоит от ключевой.
func applyOffsetToRef(cell eval.CellReference, offsetX, offsetY int) eval.CellReference {
	if !cell.AnchoredX {
		cell.X += offsetX
	}
	if !cell.AnchoredY {
		cell.Y += offsetY
	}
	return cell
}

package sheet

const (
	CellValueTypeGeneral = iota
	CellValueTypeText
)

type Cell struct {
	valueType   int
	staticValue string
}

func NewCellGeneral() *Cell {
	return &Cell{
		valueType: CellValueTypeGeneral,
	}
}

func NewCellText(text string) *Cell {
	return &Cell{
		valueType:   CellValueTypeText,
		staticValue: text,
	}
}

func (c *Cell) DisplayText() string {
	if c.valueType == CellValueTypeText {
		return c.staticValue
	}
	return ""
}

func (c *Cell) Value() string {
	// TODO: other types of values
	return c.staticValue
}

func (c *Cell) SetValueText(v string) {
	c.valueType = CellValueTypeText
	c.staticValue = v
}

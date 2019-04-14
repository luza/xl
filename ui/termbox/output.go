package termbox

import (
	"bytes"
	"xl/ui"

	"github.com/gdamore/tcell"
)

const (
	pixelsInCharX = 6
	pixelsInCharY = 20

	sheetNameMaxWidth = 10
	statusLineHeight  = 1
	hRulerHeight      = 1
	formulaLineHeight = 1
)

func (t *Termbox) SetDataDelegate(delegate ui.DataDelegateInterface) {
	t.dataDelegate = delegate
}

func (t *Termbox) SetDirty(f ui.DirtyFlag) {
	t.dirty |= f
}

func (t *Termbox) SetStatus(msg string, flags int) {
	t.statusMessage = msg
	t.statusFlags = flags
	t.SetDirty(ui.DirtyStatusLine)
}

func (t *Termbox) RefreshView() {
	docView := t.dataDelegate.DocView()
	sheetView := t.dataDelegate.SheetView()

	// formula line
	if t.dirty&ui.DirtyFormulaLine > 0 {
		formulaLineView := sheetView.FormulaLineView
		currentCellName := t.dataDelegate.CellView(sheetView.Cursor.X, sheetView.Cursor.Y).Name
		t.drawCell(0, 0, t.screenWidth, formulaLineHeight, currentCellName, tcell.ColorYellow, tcell.ColorBlack)
		text := formulaLineView.DisplayText
		if formulaLineView.Expression != nil {
			var buf bytes.Buffer
			formulaLineView.Expression.Output(func(s string, t int) {
				buf.WriteString(s)
			})
			text = buf.String()
		}
		t.drawCell(len(currentCellName)+1, 0, t.screenWidth, formulaLineHeight, text, tcell.ColorWhite, tcell.ColorBlack)
	}

	// vertical ruler
	if t.dirty&ui.DirtyVRuler > 0 {
		screenY := formulaLineHeight + hRulerHeight
		cellY := sheetView.Viewport.Top
		t.vRulerWidth = 0
		for screenY < t.screenHeight-statusLineHeight {
			rowView := t.dataDelegate.RowView(cellY)
			heightChars := pixelsToCharsY(rowView.Height)
			fg := tcell.ColorWhite
			if cellY == sheetView.Cursor.Y {
				fg = tcell.ColorYellow
			}
			t.drawCell(0, screenY, len(rowView.Name)+1+1, heightChars, rowView.Name, fg, tcell.ColorBlack)
			if len(rowView.Name)+1 > t.vRulerWidth {
				t.vRulerWidth = len(rowView.Name) + 1
			}
			cellY++
			screenY += heightChars
		}
		t.calculatedViewportHeight = cellY - sheetView.Viewport.Top
	}

	// horizontal ruler
	if t.dirty&ui.DirtyHRuler > 0 {
		screenX := t.vRulerWidth
		screenY := formulaLineHeight
		cellX := sheetView.Viewport.Left
		for screenX < t.screenWidth {
			colView := t.dataDelegate.ColView(cellX)
			widthChars := pixelsToCharsX(colView.Width)
			fg := tcell.ColorWhite
			if cellX == sheetView.Cursor.X {
				fg = tcell.ColorYellow
			}
			t.drawCell(screenX, screenY, widthChars, hRulerHeight, colView.Name, fg, tcell.ColorBlack)
			cellX++
			screenX += widthChars
		}
		t.calculatedViewportWidth = cellX - sheetView.Viewport.Left
	}

	// grid
	if t.dirty&ui.DirtyGrid > 0 {
		cellY := sheetView.Viewport.Top
		screenY := formulaLineHeight + hRulerHeight
		for screenY < t.screenHeight-statusLineHeight {
			cellX := sheetView.Viewport.Left
			screenX := t.vRulerWidth
			heightChars := pixelsToCharsY(t.dataDelegate.RowView(cellY).Height)
			for screenX < t.screenWidth {
				widthChars := pixelsToCharsX(t.dataDelegate.ColView(cellX).Width)
				c := t.dataDelegate.CellView(cellX, cellY)
				text := c.DisplayText

				bgColor := tcell.ColorBlack
				if cellX%2 != 0 || cellY%2 == 0 {
					bgColor = tcell.Color235
				}
				if cellX%2 != 0 && cellY%2 == 0 {
					bgColor = tcell.Color238
				}
				if cellX == sheetView.Cursor.X && cellY == sheetView.Cursor.Y {
					t.lastCursorX = screenX
					t.lastCursorY = screenY
					t.screen.ShowCursor(screenX, screenY)
				}
				if c.Error != nil {
					text = *c.Error
					bgColor = tcell.ColorRed
				}
				t.drawCell(screenX, screenY, widthChars, heightChars, text, tcell.ColorSilver, bgColor)
				cellX++
				screenX += widthChars
			}
			cellY++
			screenY += heightChars
		}
	}

	// status line
	if t.dirty&ui.DirtyStatusLine > 0 {
		screenX := 0
		screenY := t.screenHeight - statusLineHeight
		for i, s := range docView.Sheets {
			bgColor := tcell.ColorBlack
			fgColor := tcell.ColorWhite
			if i == docView.CurrentSheetIdx {
				bgColor = tcell.ColorWhite
				fgColor = tcell.ColorBlack
			}
			t.drawCell(screenX, screenY, sheetNameMaxWidth, statusLineHeight, s, fgColor, bgColor)
			screenX += sheetNameMaxWidth
		}
		fgColor := tcell.ColorWhite
		bgColor := tcell.ColorBlack
		if t.statusFlags&ui.StatusFlagError > 0 {
			bgColor = tcell.ColorRed
		}
		t.drawCell(screenX, screenY, t.screenWidth-screenX, statusLineHeight, t.statusMessage, fgColor, bgColor)
	}
	t.dirty = 0
	t.screen.Show()
}

func (t *Termbox) drawCell(x int, y int, width int, height int, text string, fg tcell.Color, bg tcell.Color) {
	var st tcell.Style
	st = st.Background(bg)
	textAsRunes := []rune(text)
	textLen := len(textAsRunes)
	for cursorY := y; cursorY < y+height; cursorY++ {
		indexX := 0
		for cursorX := x; cursorX < x+width; cursorX++ {
			char := ' '
			st = st.Foreground(fg)
			if cursorY == y && indexX < textLen {
				if textLen > width && cursorX == x+width-1 {
					char = '>'
					st = st.Foreground(tcell.ColorYellow)
				} else {
					char = textAsRunes[indexX]
				}
			}
			t.screen.SetContent(cursorX, cursorY, char, nil, st)
			indexX++
		}
	}
}

func pixelsToCharsX(pixels int) int {
	res := pixels / pixelsInCharX
	if res < 1 {
		res = 1
	}
	return res
}

func pixelsToCharsY(pixels int) int {
	res := pixels / pixelsInCharY
	if res < 1 {
		res = 1
	}
	return res
}

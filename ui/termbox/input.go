package termbox

import (
	"xl/ui"

	"errors"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/termbox"
)

// ReadKey blocks until new key is read. Returns the key read.
func (t *Termbox) ReadKey() (ui.InputEventInterface, error) {
	event := termbox.PollEvent()
	if event.Type == termbox.EventKey {
		e := ui.KeyEvent{
			Mod: tcell.ModMask(event.Mod),
			Key: tcell.Key(event.Key),
			Ch:  event.Ch,
		}
		return e, nil
	} else {
		return nil, errors.New("unknown event")
	}
}

func (t *Termbox) EditCellValue(oldValue string) (string, error) {
	w, _ := termbox.Size()
	v, err := t.enterEditorMode(&editorConfig{
		X:        0,
		Y:        0,
		Width:    w,
		Height:   formulaLineHeight,
		MaxLines: 1,
		FgColor:  colorWhite,
		BgColor:  colorBlack,
		Value:    oldValue,
	})
	if err != nil {
		return "", err
	}
	return v, nil
}

func (t *Termbox) InputCommand() (string, error) {
	w, h := termbox.Size()
	drawCell(0, h-statusLineHeight, 1, statusLineHeight, ":", colorWhite, colorBlack)
	v, err := t.enterEditorMode(&editorConfig{
		X:        1,
		Y:        h - statusLineHeight,
		Width:    w - 1,
		Height:   statusLineHeight,
		MaxLines: 1,
		FgColor:  colorWhite,
		BgColor:  colorBlack,
	})
	if err != nil {
		return "", err
	}
	return v, nil
}

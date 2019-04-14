package termbox

import (
	"xl/ui"

	"errors"

	"github.com/gdamore/tcell"
)

// ReadKey blocks until new key is read. Returns the key read.
func (t *Termbox) ReadKey() (ui.InputEventInterface, error) {
	ev := t.screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		e := ui.KeyEvent{
			Mod: ev.Modifiers(),
			Key: ev.Key(),
			Ch:  ev.Rune(),
		}
		return e, nil

	case *tcell.EventResize:
	//handling resize
	case *tcell.EventMouse:
	//handling mouse
	case *tcell.EventError:
		return nil, errors.New("unknown event")
	}
	return nil, nil
}

func (t *Termbox) EditCellValue(oldValue string) (string, error) {
	w, _ := t.screen.Size()
	v, err := t.enterEditorMode(&editorConfig{
		Tbox:     t,
		X:        0,
		Y:        0,
		Width:    w,
		Height:   formulaLineHeight,
		MaxLines: 1,
		FgColor:  tcell.ColorWhite,
		BgColor:  tcell.ColorBlack,
		Value:    oldValue,
	})
	if err != nil {
		return "", err
	}
	return v, nil
}

func (t *Termbox) InputCommand() (string, error) {
	w, h := t.screen.Size()
	t.drawCell(0, h-statusLineHeight, 1, statusLineHeight, ":", tcell.ColorWhite, tcell.ColorBlack)
	v, err := t.enterEditorMode(&editorConfig{
		Tbox:     t,
		X:        1,
		Y:        h - statusLineHeight,
		Width:    w - 1,
		Height:   statusLineHeight,
		MaxLines: 1,
		FgColor:  tcell.ColorWhite,
		BgColor:  tcell.ColorBlack,
	})
	if err != nil {
		return "", err
	}
	return v, nil
}

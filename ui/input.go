package ui

import (
	"github.com/nsf/termbox-go"
)

type InputInterface interface {
	ReadKey() (InputEventInterface, error)
}

type InputEventInterface interface {
}

type KeyEvent struct {
	InputEventInterface

	Mod termbox.Modifier // one of Mod* constants or 0
	Key termbox.Key      // one of Key* constants, invalid if 'Ch' is not 0
	Ch  rune             // a unicode character
}

package ui

import (
	"github.com/gdamore/tcell"
)

type InputInterface interface {
	ReadKey() (InputEventInterface, error)
}

type InputEventInterface interface {
}

type KeyEvent struct {
	InputEventInterface

	Mod tcell.ModMask
	Key tcell.Key
	Ch  rune
}

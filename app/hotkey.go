package app

import (
	"github.com/gdamore/tcell"
)

type Key struct {
	Mod tcell.ModMask
	Key tcell.Key
	Ch  rune
}

var HotKeys = map[string]Key{
	"a": {tcell.ModNone, tcell.KeyRune, 'a'},
	"b": {tcell.ModNone, tcell.KeyRune, 'b'},
	"c": {tcell.ModNone, tcell.KeyRune, 'c'},
	"d": {tcell.ModNone, tcell.KeyRune, 'd'},
	"e": {tcell.ModNone, tcell.KeyRune, 'e'},
	"f": {tcell.ModNone, tcell.KeyRune, 'f'},
	"g": {tcell.ModNone, tcell.KeyRune, 'g'},
	"h": {tcell.ModNone, tcell.KeyRune, 'h'},
	"i": {tcell.ModNone, tcell.KeyRune, 'i'},
	"j": {tcell.ModNone, tcell.KeyRune, 'j'},
	"k": {tcell.ModNone, tcell.KeyRune, 'k'},
	"l": {tcell.ModNone, tcell.KeyRune, 'l'},
	"m": {tcell.ModNone, tcell.KeyRune, 'm'},
	"n": {tcell.ModNone, tcell.KeyRune, 'n'},
	"o": {tcell.ModNone, tcell.KeyRune, 'o'},
	"p": {tcell.ModNone, tcell.KeyRune, 'p'},
	"q": {tcell.ModNone, tcell.KeyRune, 'q'},
	"r": {tcell.ModNone, tcell.KeyRune, 'r'},
	"s": {tcell.ModNone, tcell.KeyRune, 's'},
	"t": {tcell.ModNone, tcell.KeyRune, 't'},
	"u": {tcell.ModNone, tcell.KeyRune, 'u'},
	"v": {tcell.ModNone, tcell.KeyRune, 'v'},
	"w": {tcell.ModNone, tcell.KeyRune, 'w'},
	"x": {tcell.ModNone, tcell.KeyRune, 'x'},
	"y": {tcell.ModNone, tcell.KeyRune, 'y'},
	"z": {tcell.ModNone, tcell.KeyRune, 'z'},
	">": {tcell.ModNone, tcell.KeyRune, '>'},
	"}": {tcell.ModNone, tcell.KeyRune, '}'},
	"{": {tcell.ModNone, tcell.KeyRune, '{'},
}

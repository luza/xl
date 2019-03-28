package app

import "github.com/gdamore/tcell/termbox"

type Key struct {
	Mod termbox.Modifier // one of Mod* constants or 0
	Key termbox.Key      // one of Key* constants, invalid if 'Ch' is not 0
	Ch  rune             // a unicode character
}

var HotKeys = map[string]Key{
	"a": {0, 0, 'a'},
	"b": {0, 0, 'b'},
	"c": {0, 0, 'c'},
	"d": {0, 0, 'd'},
	"e": {0, 0, 'e'},
	"f": {0, 0, 'f'},
	"g": {0, 0, 'g'},
	"h": {0, 0, 'h'},
	"i": {0, 0, 'i'},
	"j": {0, 0, 'j'},
	"k": {0, 0, 'k'},
	"l": {0, 0, 'l'},
	"m": {0, 0, 'm'},
	"n": {0, 0, 'n'},
	"o": {0, 0, 'o'},
	"p": {0, 0, 'p'},
	"q": {0, 0, 'q'},
	"r": {0, 0, 'r'},
	"s": {0, 0, 's'},
	"t": {0, 0, 't'},
	"u": {0, 0, 'u'},
	"v": {0, 0, 'v'},
	"w": {0, 0, 'w'},
	"x": {0, 0, 'x'},
	"y": {0, 0, 'y'},
	"z": {0, 0, 'z'},
	">": {0, 0, '>'},
	"}": {0, 0, '}'},
	"{": {0, 0, '{'},
}

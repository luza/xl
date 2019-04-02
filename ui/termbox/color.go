package termbox

import (
	"github.com/gdamore/tcell"
)

const (
	colorGrey   int32 = 0x9d9d9d
	colorRed    int32 = 0xff0000
	colorYellow int32 = 0xffff00
	colorWhite  int32 = 0xffffff
	colorBlack  int32 = 0x000000

	colorGrey236 int32 = 0xe3e3e3
	colorGrey239 int32 = 0xdcdcdc
)

// Init initializes the screen for use.
func Init() (error, tcell.Screen) {

	if s, e := tcell.NewScreen(); e != nil {
		return e, nil
	} else if e = s.Init(); e != nil {
		return e, nil
	} else {
		return nil, s
	}
}

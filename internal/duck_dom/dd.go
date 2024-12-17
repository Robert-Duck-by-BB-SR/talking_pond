package duckdom

import (
	"fmt"
)

const (
	FG_KEY        = "\033[38;2;"
	BG_KEY        = "\033[48;2;"
	INVERT_STYLES = "\033[7m"
	RED_COLOR     = "\033[31m"
	RESET_STYLES  = "\033[0m"

	// constant commands
	CLEAR_SCREEN                 = "\033[2J"
	MOVE_CURSOR_TO_THE_BENINGING = "\033[H"
	MOVE_CURSOR_TO_POSITION      = "\033[%d;%dH"
	CLEAR_ROW                    = "\033[2K"
	HIDE_CURSOR                  = "\x1b[?25l"

	// NOTE: DEBUG ONLY. IF YOU USE IT IN PROD I WILL FIND YOU AND MAKE YOU SMELL MY SOCKS
	DEBUG_STYLES = "\033[30;43m"
)

type Position struct {
	StartingRow, StartingCol uint
}

type Screen struct {
	MaxRows        int
	MaxCols        int
	CursorPosition Position
	ActiveWindowId int

	// fuck Windows, all my homies use Linux
	Windows     []Window
	RenderQueue []string
}

func (self *Screen) Render() {
	for renderable := range self.RenderQueue {
		fmt.Print(renderable)
	}
}

func ClearScreen() {
	fmt.Printf(CLEAR_SCREEN)
	fmt.Printf(MOVE_CURSOR_TO_THE_BENINGING)
}

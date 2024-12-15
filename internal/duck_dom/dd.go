package duckdom

import (
	"fmt"
)

const (
	FG_KEY = "\033[38;2;"
	BG_KEY = "\033[48;2;"
	INVERT_STYLES = "\033[7m"
	RED_COLOR= "\033[31m"
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
	CursorPos      Position
	ActiveWindowId int

	// fuck Windows, all my homies use Linux
	Windows     []Renderable
	RenderQueue []Renderable
}

func (self *Screen) GetPos() Position { return Position{} }
func (self *Screen) Active() Renderable { return self.Windows[self.ActiveWindowId] }
func (self *Screen) SetActive(id int)   { self.ActiveWindowId = id }
func (self *Screen) ActiveIndex() int   { return self.ActiveWindowId }

type Renderable interface {
	Stylable
	Render() string
	Active() Renderable
	SetActive(int)
	ActiveIndex() int
	GetPos() Position
}

type Stylable interface {
	SetWidth(int) Stylable
	SetHeight(int) Stylable
	SetBackground(string) Stylable
	SetBorder(Border) Stylable
}

func (self *Screen) Render() string {
	// NOTE: maybe make it fill the render q?
	return ""
}

func ClearScreen() {
	fmt.Printf(CLEAR_SCREEN)
	fmt.Printf(MOVE_CURSOR_TO_THE_BENINGING)
}

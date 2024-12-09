package duckdom

import "fmt"

const (
	// styles
	INVERT_STYLES = "\033[7m"
	RESET_STYLES  = "\033[0m"

	// constant commands
	CLEAR_SCREEN                 = "\033[2J"
	MOVE_CURSOR_TO_THE_BENINGING = "\033[H"
	CLEAR_ROW                    = "\033[2K"
)

type Position struct {
	Row, Col uint
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

type Renderable interface {
	Render() string
	SetStyle(string)
	Active() Renderable
	SetActive(int)
	ActiveIndex() int
}

func (self *Screen) Render() string {
	// NOTE: maybe make it fill the render q?
	return ""
}

func (self *Screen) SetStyle(string) {
}

func (self *Screen) Active() Renderable { return self.Windows[self.ActiveWindowId] }
func (self *Screen) SetActive(id int)   { self.ActiveWindowId = id }
func (self *Screen) ActiveIndex() int   { return self.ActiveWindowId }

func DebugMeDaddy(screen *Screen, content string) {
	fmt.Printf("\033[%d;1H", screen.MaxRows)
	fmt.Printf(CLEAR_ROW)
	fmt.Printf("\033[%d;%dH\033[30;43m%s\033[0m", screen.MaxRows, 1, "DebugDuck: "+content)
}

func ClearScreen() {
	fmt.Printf(CLEAR_SCREEN)
	fmt.Printf(MOVE_CURSOR_TO_THE_BENINGING)
}

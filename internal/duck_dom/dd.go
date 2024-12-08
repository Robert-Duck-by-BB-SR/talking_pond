package duckdom

import "fmt"

type Position struct {
	Row, Col uint
}

type Screen struct {
	Max_rows           int
	Max_cols           int
	Cursor_pos         [2]uint
	Active_window_indx int
	// fuck Windows, all my homies use Linux
	Windows      []Renderable
	Render_queue []Renderable
}

type Renderable interface {
	Render() string
	SetPos(Position)
	GetPos() Position
	SetStyle(string)
}

const (
	// styles
	INVERT_STYLES = "\033[7m"
	RESET_STYLES  = "\033[0m"

	// constant commands
	CLEAR_SCREEN                 = "\033[2J"
	MOVE_CURSOR_TO_THE_BENINGING = "\033[H"
	CLEAR_ROW                    = "\033[2K"
)

func Debug_me_daddy(screen *Screen, content string) {
	fmt.Printf("\033[%d;1H", screen.Max_rows)
	fmt.Printf(CLEAR_ROW)
	fmt.Printf("\033[%d;%dH\033[30;43m%s\033[0m", screen.Max_rows, 1, "DebugDuck: "+content)
}

func Clear_screen() {
	fmt.Printf(CLEAR_SCREEN)
	fmt.Printf(MOVE_CURSOR_TO_THE_BENINGING)
}
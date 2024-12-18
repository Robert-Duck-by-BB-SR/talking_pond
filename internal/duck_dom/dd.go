package duckdom

import (
	"fmt"
)

const (
	NORMAL  = "NORMAL"
	INSERT  = "INSERT"
	COMMAND = "COMMAND"
	WINDOW  = "WINDOW"
)

const (
	FG_KEY        = "\033[38;2;"
	BG_KEY        = "\033[48;2;"
	INVERT_STYLES = "\033[7m"
	RED_COLOR     = "\033[31m"
	RESET_STYLES  = "\033[0m"
	RGB           = "%d;%d;%dm"

	// constant commands
	CLEAR_SCREEN                 = "\033[2J"
	MOVE_CURSOR_TO_THE_BENINGING = "\033[H"
	MOVE_CURSOR_TO_POSITION      = "\033[%d;%dH"
	CLEAR_ROW                    = "\033[2K"
	HIDE_CURSOR                  = "\x1b[?25l"
	SHOW_CURSOR                  = "\x1b[?25h"

	// NOTE: DEBUG ONLY. IF YOU USE IT IN PROD I WILL FIND YOU AND MAKE YOU SMELL MY SOCKS
	DEBUG_STYLES = "\033[30;43m"
)

type Position struct {
	StartingRow, StartingCol uint
}

type Screen struct {
	// Width = max number of columns for terminal window
	Width int
	// Height = max number of row for terminal window
	Height             int
	CursorPosition     Position
	ActiveWindowId     int
	EventLoopIsRunning bool
	State
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
	fmt.Print(CLEAR_SCREEN)
	fmt.Print(MOVE_CURSOR_TO_THE_BENINGING)
	fmt.Print(HIDE_CURSOR)
}

type State interface {
	HandleKeypress(*Screen, []byte)
}

type NormalMode struct{}

var Normal NormalMode

func cycle_index(new, len int) int {
	if new < 0 {
		return len - 1
	}
	if new >= len {
		return 0
	}
	return new
}

func (self *Screen) change_window(direction int) {
	DebugMeDaddy(self, fmt.Sprint(len(self.Windows)))
	// TODO: select active item and deselect it visually
	prev_window := self.Windows[self.ActiveWindowId]
	prev_window.Border.Style = RoundedBorder
	self.ActiveWindowId = cycle_index(self.ActiveWindowId+direction, len(self.Windows))
	new_window := self.Windows[self.ActiveWindowId]
	new_window.Border.Style = BoldBorder
	self.RenderQueue = append(
		self.RenderQueue,
		render_border(prev_window.Position, &prev_window.Styles),
		render_border(new_window.Position, &new_window.Styles),
	)
	if len(prev_window.Components) > 0 {
		self.RenderQueue = append(self.RenderQueue, prev_window.Components[new_window.ActiveComponentId].Render())
	}

	if len(new_window.Components) > 0 {
		self.RenderQueue = append(self.RenderQueue, new_window.Components[new_window.ActiveComponentId].Render())
	}
}

func (self *Screen) change_component(direction int) {
	active_window := &self.Windows[self.ActiveWindowId]
	if len(active_window.Components) > 0 {
		prev_component := active_window.Components[active_window.ActiveComponentId]
		self.RenderQueue = append(
			self.RenderQueue,
			RESET_STYLES+prev_component.Render())

		active_window.ActiveComponentId = cycle_index(active_window.ActiveComponentId+direction, len(active_window.Components))
		if len(active_window.Components) > 0 {
			new_component := active_window.Components[active_window.ActiveComponentId]
			self.RenderQueue = append(
				self.RenderQueue,
				INVERT_STYLES+new_component.Render(),
			)
		}
	}
}

func (self *Screen) Activate() {
	self.change_component(0)
}

func (*NormalMode) HandleKeypress(screen *Screen, keys []byte) {
	// big ass switch case
	switch keys[0] {
	case 'q':
		screen.EventLoopIsRunning = false
	// Ctrl+W (ASCII 23 or \x17)
	case 23:
		screen.State = &WM
	case 'l':
		fallthrough
	case 'j':
		screen.change_component(+1)
	case 'k':
		fallthrough
	case 'h':
		screen.change_component(-1)
	// switching modes
	case ':':
		screen.State = &Command
	case 'i':
		screen.State = &Insert
	}
}

type InsertMode struct{}

var Insert InsertMode

func (*InsertMode) HandleKeypress(screen *Screen, keys []byte) {
	switch keys[0] {
	case '':
		screen.State = &Normal
	case 'j':
		screen.RenderQueue = append(screen.RenderQueue, "jjjjjjjj")
	case 'i':
		screen.RenderQueue = append(screen.RenderQueue, "iiiiiii")
	}
}

type CommandMode struct{}

var Command CommandMode

func (*CommandMode) HandleKeypress(screen *Screen, keys []byte) {
	switch keys[0] {
	case '':
		screen.State = &Normal
	}
}

type WindowMode struct{}

// FIXME: PAPI RENAME
var WM WindowMode

func (*WindowMode) HandleKeypress(screen *Screen, keys []byte) {
	switch keys[0] {
	case '':
		screen.State = &Normal
	case 'l':
		fallthrough
	case 'j':
		screen.change_window(+1)
	case 'k':
		fallthrough
	case 'h':
		screen.change_window(-1)
	}
}

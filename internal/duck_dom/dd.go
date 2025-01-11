package duckdom

import (
	"fmt"

	tpc "github.com/Robert-Duck-by-BB-SR/talking_pond/internal/tps_client"
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
	Row, Col int
}

type Screen struct {
	// Width = max number of columns for terminal window
	Width int
	// Height = max number of row for terminal window
	Height             int
	CursorPosition     Position
	ActiveWindowId     int
	EventLoopIsRunning bool
	StatusBar          Window
	State
	// fuck Windows, all my homies use Linux
	Windows     []*Window
	RenderQueue []string
	Client      tpc.Client
}

func (self *Screen) Render() {
	for _, window := range self.Windows {
		self.RenderQueue = append(self.RenderQueue, window.Render())
	}

	self.RenderQueue = append(self.RenderQueue, self.StatusBar.Render())

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

func cycle_index(new, len int) int {
	if new < 0 {
		return len - 1
	}
	if new >= len {
		return 0
	}
	return new
}

func (self *Screen) change_window(id int) {
	old_window := self.Windows[self.ActiveWindowId]
	old_window.Border.Color = PRIMARY_THEME.SecondaryTextColor
	old_window.Border.Style = NormalBorder
	old_window.Active = false
	self.ActiveWindowId = id
	new_window := self.Windows[self.ActiveWindowId]
	new_window.Border.Color = PRIMARY_THEME.ActiveTextColor
	new_window.Border.Style = BoldBorder
	new_window.Active = true
	self.RenderQueue = append(
		self.RenderQueue,
		render_border(old_window.Position, old_window.Active, &old_window.Styles),
		render_border(new_window.Position, new_window.Active, &new_window.Styles),
	)
	if len(old_window.Components) > 0 {
		old_active := old_window.Components[old_window.ActiveComponentId]
		old_active.Active = false
		self.RenderQueue = append(
			self.RenderQueue,
			old_active.Render(),
		)
	}
	if len(new_window.Components) > 0 {
		new_active := new_window.Components[new_window.ActiveComponentId]
		new_active.Active = true
		self.RenderQueue = append(
			self.RenderQueue,
			new_active.Render(),
		)
	}
}

func (self *Screen) change_component(id int) {
	active_window := self.Windows[self.ActiveWindowId]
	if len(active_window.Components) > 0 {
		prev_component := active_window.Components[active_window.ActiveComponentId]
		prev_component.Active = false
		self.RenderQueue = append(
			self.RenderQueue,
			prev_component.Render(),
		)

		active_window.ActiveComponentId = id
		if len(active_window.Components) > 0 {
			new_component := active_window.Components[active_window.ActiveComponentId]
			new_component.Active = true
			self.RenderQueue = append(
				self.RenderQueue,
				new_component.Render(),
			)
		}
	}
}

func (self *Screen) Activate() {
	self.change_window(0)
	self.change_component(0)
}

func (self *Screen) AddWindow(w *Window) {
	w.Index = len(self.Windows)
	w.Oldfart = self
	self.Windows = append(self.Windows, w)
}

type NormalMode struct{}

var Normal NormalMode

func (*NormalMode) HandleKeypress(screen *Screen, keys []byte) {
	// big ass switch case
	active_window := screen.Windows[screen.ActiveWindowId]
	switch keys[0] {
	case 'q':
		screen.EventLoopIsRunning = false
	case '':
		screen.change_state(&WM, WINDOW)
	case 'l':
		fallthrough
	case 'j':
		index := cycle_index(active_window.ActiveComponentId+1, len(active_window.Components))
		screen.change_component(index)
	case 'k':
		fallthrough
	case 'h':
		index := cycle_index(active_window.ActiveComponentId-1, len(active_window.Components))
		screen.change_component(index)
	case ':':
		screen.change_state(&Command, COMMAND)
	case 'i':
		active_window := screen.Windows[screen.ActiveWindowId]
		active_component := active_window.Components[active_window.ActiveComponentId]
		if active_component.Inputable {
			screen.change_state(&Insert, INSERT)
		}
	case 'I':
		screen.change_state(&Insert, INSERT)
		screen.change_window(len(screen.Windows) - 1)
		screen.change_component(0)
	case '':
		active_window := screen.Windows[screen.ActiveWindowId]
		active_component := active_window.Components[active_window.ActiveComponentId]
		active_component.Action()
	}
}

type InsertMode struct{}

var Insert InsertMode

func (*InsertMode) HandleKeypress(screen *Screen, keys []byte) {
	switch keys[0] {
	case '':
		fallthrough
	case '':
		screen.change_state(&Normal, NORMAL)
	default:
		active_window := screen.Windows[screen.ActiveWindowId]
		active_component := active_window.Components[active_window.ActiveComponentId]
		active_component.Buffer += string(keys[0])
		screen.RenderQueue = append(screen.RenderQueue, active_window.Render())
	}
}

func (screen *Screen) change_state(state State, state_name string) {
	screen.State = state
	screen.StatusBar.Components[0].Buffer = state_name
	screen.RenderQueue = append(screen.RenderQueue, screen.StatusBar.Render())
}

type CommandMode struct{}

var Command CommandMode

func (*CommandMode) HandleKeypress(screen *Screen, keys []byte) {
	switch keys[0] {
	case '':
		fallthrough
	case '':
		screen.change_state(&Normal, NORMAL)
	}
}

type WindowMode struct{}

var WM WindowMode

func (*WindowMode) HandleKeypress(screen *Screen, keys []byte) {
	switch keys[0] {
	case '':
		fallthrough
	case '':
		fallthrough
	case '':
		screen.change_state(&Normal, NORMAL)
	case 'l':
		fallthrough
	case 'j':
		id := cycle_index(screen.ActiveWindowId+1, len(screen.Windows))
		screen.change_window(id)
	case 'k':
		fallthrough
	case 'h':
		id := cycle_index(screen.ActiveWindowId-1, len(screen.Windows))
		screen.change_window(id)
	}
}

package duckdom

import (
	"bufio"
	"log"
	"os"
	"strings"

	tpc "github.com/Robert-Duck-by-BB-SR/talking_pond/internal/tps_client"
)

var DEBUG_MODE = false

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
	StatusBar Window
	// fuck Windows, all my homies use Linux
	Client      tpc.Client
	Windows     []*Window
	RenderQueue strings.Builder

	CursorPosition Position
	State

	// Width = max number of columns for terminal window
	Width int
	// Height = max number of row for terminal window
	Height             int
	ActiveWindowId     int
	EventLoopIsRunning bool
	ModalIsActive      bool
}

// Renders everything there is in screen. Uses screen.Render
func (self *Screen) RenderFull() {
	for _, window := range self.Windows {
		window.Render(&self.RenderQueue)
	}

	self.StatusBar.Render(&self.RenderQueue)

	if self.ModalIsActive {
		self.ActivateModal()
	} else {
		self.Activate()
	}
	self.Render()
}

// Dumps everything there is in RenderQueue into stdout and resets the RenderQueue.
func (self *Screen) Render() {
	writer := bufio.NewWriter(os.Stdout)
	if _, err := writer.WriteString(self.RenderQueue.String()); err != nil {
		log.Fatalln(err)
	}
	if err := writer.Flush(); err != nil {
		log.Fatalln(err)
	}
	self.RenderQueue.Reset()
}

func ClearScreen() {
	writer := bufio.NewWriter(os.Stdout)
	if _, err := writer.WriteString(CLEAR_SCREEN); err != nil {
		log.Fatalln(err)
	}
	if _, err := writer.WriteString(MOVE_CURSOR_TO_THE_BENINGING); err != nil {
		log.Fatalln(err)
	}
	if _, err := writer.WriteString(HIDE_CURSOR); err != nil {
		log.Fatalln(err)
	}
	if err := writer.Flush(); err != nil {
		log.Fatalln(err)
	}
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
	render_border(&self.RenderQueue, old_window.Position, old_window.Active, &old_window.Styles)
	render_border(&self.RenderQueue, new_window.Position, new_window.Active, &new_window.Styles)
	if len(old_window.Components) > 0 {
		old_active := old_window.Components[old_window.ActiveComponentId]
		old_active.Active = false
		old_active.Render(&self.RenderQueue)
	}
	if len(new_window.Components) > 0 {
		new_active := new_window.Components[new_window.ActiveComponentId]
		new_active.Active = true
		new_active.Render(&self.RenderQueue)
	}
}

func (self *Screen) change_component(id int) {
	active_window := self.Windows[self.ActiveWindowId]
	if len(active_window.Components) > 0 {
		prev_component := active_window.Components[active_window.ActiveComponentId]
		prev_component.Active = false
		prev_component.Render(&self.RenderQueue)

		active_window.ActiveComponentId = id
		if len(active_window.Components) > 0 {
			new_component := active_window.Components[active_window.ActiveComponentId]
			new_component.Active = true
			new_component.Render(&self.RenderQueue)
		}
	}
}

// rerenders first window and its active component
func (self *Screen) Activate() {
	self.change_window(0)
	self.change_component(0)
}

// rerenders the last window and its active component
func (self *Screen) ActivateModal() {
	self.ModalIsActive = true
	self.change_window(len(self.Windows) - 1)
	self.change_component(0)
	active_window := self.get_active_window()
	active_window.Render(&self.RenderQueue)
}

func (self *Screen) AddWindow(w *Window) {
	w.Index = len(self.Windows)
	w.Oldfart = self
	self.Windows = append(self.Windows, w)
}

func (self *Screen) get_active_component() *Component {
	active_window := self.Windows[self.ActiveWindowId]
	return active_window.Components[active_window.ActiveComponentId]
}

func (self *Screen) get_active_window() *Window {
	active_window := self.Windows[self.ActiveWindowId]
	return active_window
}

func (screen *Screen) CloseModal() {
	active_window := screen.Windows[screen.ActiveWindowId]
	screen.change_window(0)
	active_window = screen.Windows[0]
	screen.change_component(active_window.ActiveComponentId)
	screen.ModalIsActive = false
	screen.Windows = screen.Windows[:len(screen.Windows)-1]
	screen.RenderFull()
}

type NormalMode struct{}

var Normal NormalMode

func (*NormalMode) HandleKeypress(screen *Screen, keys []byte) {
	// big ass switch case
	active_window := screen.Windows[screen.ActiveWindowId]
	switch keys[0] {

	case 'q':
		if screen.ModalIsActive {
			screen.CloseModal()
		}
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
	case '':
		active_component := screen.get_active_component()
		if active_component.ScrollType == VERTICAL {
			active_component.BufferVerticalFrom -= 1
			active_component.Render(&screen.RenderQueue)
		}
	case '':
		active_component := screen.get_active_component()
		if active_component.ScrollType == VERTICAL {
			active_component.BufferVerticalFrom += 1
			active_component.Render(&screen.RenderQueue)
		}
	case 'w':
		active_component := screen.get_active_component()
		if active_component.ScrollType == HORIZONTAL {
			active_component.BufferHorizontalFrom += 1
			active_component.Render(&screen.RenderQueue)
		}
	case 'b':
		active_component := screen.get_active_component()
		if active_component.ScrollType == HORIZONTAL {
			active_component.BufferHorizontalFrom -= 1
			active_component.Render(&screen.RenderQueue)
		}
	case ':':
		screen.change_state(&Command, CLEAR_ROW+":")
	case 'i':
		active_component := screen.get_active_component()
		if active_component.Inputable {
			screen.change_state(&Insert, INSERT)
		}
	case 'I':
		screen.change_state(&Insert, INSERT)
		screen.change_window(len(screen.Windows) - 1)
		screen.change_component(0)
	case '
':
		active_component := screen.get_active_component()
		if active_component.Action != nil {
			active_component.Action()
		}
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
	case 8, 127:
		active_component := screen.get_active_component()
		if len(active_component.Buffer) != 0{
			active_component.Buffer = active_component.Buffer[:len(active_component.Buffer) - 1]
			active_window := screen.Windows[screen.ActiveWindowId]
			active_window.Render(&screen.RenderQueue)
		}
	default:
		active_component := screen.get_active_component()
		active_component.Buffer += string(keys[0])
		active_component.Render(&screen.RenderQueue)
	}
}

func (screen *Screen) change_state(state State, state_name string) {
	screen.State = state
	screen.StatusBar.Components[0].Buffer = state_name
	screen.StatusBar.Render(&screen.RenderQueue)
}

type CommandMode struct{}

var Command CommandMode

func (*CommandMode) HandleKeypress(screen *Screen, keys []byte) {
	status_line := screen.StatusBar.Components[0]
	switch keys[0] {
	case '':
		fallthrough
	case '':
		screen.change_state(&Normal, NORMAL)
	case '
':
		status_line.Action()
		screen.change_state(&Normal, NORMAL)
	default:
		status_line.Buffer += string(keys[0])
		status_line.Render(&screen.RenderQueue)
	}
}

type WindowMode struct{}

var WM WindowMode

func (*WindowMode) HandleKeypress(screen *Screen, keys []byte) {
	switch keys[0] {
	case '':
		fallthrough
	case '
':
		fallthrough
	case '':
		screen.change_state(&Normal, NORMAL)
	case '':
		// move half page up
	case '':
		// move half page down
	case 'l':
		fallthrough
	case 'j':
		if !screen.ModalIsActive {
			id := cycle_index(screen.ActiveWindowId+1, len(screen.Windows))
			screen.change_window(id)
		}
	case 'k':
		fallthrough
	case 'h':
		if !screen.ModalIsActive {
			id := cycle_index(screen.ActiveWindowId-1, len(screen.Windows))
			screen.change_window(id)
		}
	}
}

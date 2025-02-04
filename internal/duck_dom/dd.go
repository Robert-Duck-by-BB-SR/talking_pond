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
	HIDDEN_CURSOR                = "\033[?25l"
	VISIBLE_CURSOR               = "\033[?25h"

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
	WriteToQ           chan string
	ReadFromQ          chan q_reader
}

type q_reader struct {
	response chan string
}

func (self *Screen) RenderQueueStart() {
	self.WriteToQ = make(chan string)
	self.ReadFromQ = make(chan q_reader)
	for {
		select {
		case text := <-self.WriteToQ:
			self.RenderQueue.WriteString(text)
		case reader := <-self.ReadFromQ:
			reader.response <- self.RenderQueue.String()
			self.RenderQueue.Reset()
		}
	}
}

// Renders everything there is in screen. Uses screen.Render
func (self *Screen) RenderFull() {
	for _, window := range self.Windows {
		self.WriteToQ <- window.Render()
	}

	self.WriteToQ <- self.StatusBar.Render()

	if self.ModalIsActive {
		self.ActivateModal()
	} else {
		self.Activate(self.ActiveWindowId)
	}
	self.Render()
}

// Dumps everything there is in RenderQueue into stdout and resets the RenderQueue.
func (self *Screen) Render() {
	screen := q_reader{}
	screen.response = make(chan string)
	self.ReadFromQ <- screen
	text := <-screen.response
	if len(text) != 0 {
		writer := bufio.NewWriter(os.Stdout)
		if _, err := writer.WriteString(text); err != nil {
			log.Fatalln(err)
		}
		if err := writer.Flush(); err != nil {
			log.Fatalln(err)
		}
	}
}

func ClearScreen() {
	writer := bufio.NewWriter(os.Stdout)
	if _, err := writer.WriteString(CLEAR_SCREEN); err != nil {
		log.Fatalln(err)
	}
	if _, err := writer.WriteString(MOVE_CURSOR_TO_THE_BENINGING); err != nil {
		log.Fatalln(err)
	}
	if _, err := writer.WriteString(HIDDEN_CURSOR); err != nil {
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
	var builder strings.Builder
	defer builder.Reset()
	old_window := self.Windows[self.ActiveWindowId]
	old_window.Border.Color = PRIMARY_THEME.SecondaryTextColor
	old_window.Border.Style = NormalBorder
	old_window.Active = false
	self.ActiveWindowId = id
	new_window := self.Windows[self.ActiveWindowId]
	new_window.Border.Color = PRIMARY_THEME.ActiveTextColor
	new_window.Border.Style = BoldBorder
	new_window.Active = true
	render_border(&builder, old_window.Position, old_window.Active, &old_window.Styles)
	render_border(&builder, new_window.Position, new_window.Active, &new_window.Styles)
	if len(old_window.Components) > 0 {
		old_active := old_window.Components[old_window.ActiveComponentId]
		old_active.Active = false
		builder.WriteString(old_active.Render())
	}
	if len(new_window.Components) > 0 {
		new_active := new_window.Components[new_window.ActiveComponentId]
		new_active.Active = true
		builder.WriteString(new_active.Render())
	}
	self.WriteToQ <- builder.String()
}

func (self *Screen) change_component(id int) {
	var builder strings.Builder
	defer builder.Reset()
	active_window := self.Windows[self.ActiveWindowId]
	if len(active_window.Components) > 0 {
		prev_component := active_window.Components[active_window.ActiveComponentId]
		prev_component.Active = false
		builder.WriteString(prev_component.Render())

		active_window.ActiveComponentId = id
		if len(active_window.Components) > 0 {
			new_component := active_window.Components[active_window.ActiveComponentId]
			new_component.Active = true
			builder.WriteString(new_component.Render())
		}
	}
	self.WriteToQ <- builder.String()
}

// rerenders first window and its active component
// FIXME: should not take an argument?
func (self *Screen) Activate(i int) {
	self.change_window(i)
	self.change_component(i)
}

// rerenders the last window and its active component
func (self *Screen) ActivateModal() {
	self.ModalIsActive = true
	self.change_window(len(self.Windows) - 1)
	self.change_component(0)
	active_window := self.get_active_window()
	self.WriteToQ <- active_window.Render()
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
			active_component.buffer_vertical_scroll_from -= 1
			screen.WriteToQ <- active_component.Render()
		}
	case '':
		active_component := screen.get_active_component()
		if active_component.ScrollType == VERTICAL {
			active_component.buffer_vertical_scroll_from += 1
			screen.WriteToQ <- active_component.Render()
		}
	case 'w':
		active_component := screen.get_active_component()
		if active_component.ScrollType == HORIZONTAL {
			active_component.buffer_horizontal_scroll_from += 1
			screen.WriteToQ <- active_component.Render()
		}
	case 'b':
		active_component := screen.get_active_component()
		if active_component.ScrollType == HORIZONTAL {
			active_component.buffer_horizontal_scroll_from -= 1
			screen.WriteToQ <- active_component.Render()
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
	case '':
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
	case '':
		active_component := screen.get_active_component()
		if active_component.Inputable && active_component.Action != nil {
			active_component.Action()
		}
	case 8, 127:
		active_component := screen.get_active_component()
		if len(active_component.Buffer) != 0 {
			active_component.Buffer = active_component.Buffer[:len(active_component.Buffer)-1]
			active_component.render_content(&screen.RenderQueue)
		}
	default:
		active_component := screen.get_active_component()
		active_component.Buffer += string(keys[0])
		active_component.buffer_vertical_scroll_from += 1
		active_component.render_content(&screen.RenderQueue)
	}
}

func (screen *Screen) change_state(state State, state_name string) {
	if screen.State == &Insert || screen.State == &Command {
		screen.RenderQueue.WriteString(HIDDEN_CURSOR)
	}
	if state == &Insert || state == &Command {
		screen.RenderQueue.WriteString(VISIBLE_CURSOR)
	}
	screen.State = state
	screen.StatusBar.Components[0].Buffer = state_name
	screen.WriteToQ <- screen.StatusBar.Render()
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
	case '':
		status_line.Action()
		screen.change_state(&Normal, NORMAL)
	case 8, 127:
		if len(status_line.Buffer) != 0 && status_line.Buffer[len(status_line.Buffer)-1] != ':' {
			status_line.Buffer = status_line.Buffer[:len(status_line.Buffer)-1]
			screen.WriteToQ <- status_line.Render()
		}
	default:
		status_line.Buffer += string(keys[0])
		screen.WriteToQ <- status_line.Render()
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

func CreateMessages(content *Window, conversation string, message []string) {
	// for _, m := range message {
	// 	data := strings.Split(m, string([]byte{255}))
	// 	if len(data) < 5 {
	// 		continue
	// 	}
	// 	_ = data[0] // type
	// 	user := data[1]
	// 	convo := data[2]
	// 	datetime := data[3]
	// 	text := data[4]
	// 	if convo == conversation {
	// 		item := CreateComponent(text,
	// 			Styles{
	// 				MaxWidth:   content.Width - 4,
	// 				TextColor:  PRIMARY_THEME.SecondaryTextColor,
	// 				Background: PRIMARY_THEME.ActiveBg,
	// 				Border:     Border{Style: RoundedBorder, Color: RED_COLOR},
	// 			},
	// 		)
	// 		content.AddComponent(item)
	// 		user_date := CreateComponent(fmt.Sprintf("%s | %s", user, datetime), Styles{
	// 			MaxWidth:   content.Width - 4,
	// 			MaxHeight:  1,
	// 			TextColor:  PRIMARY_THEME.SecondaryTextColor,
	// 			Background: PRIMARY_THEME.ActiveBg,
	// 		})
	// 		content.AddComponent(user_date)
	// 		item.Render(&content.Oldfart.RenderQueue)
	// 		user_date.Render(&content.Oldfart.RenderQueue)
	// 	}
	// }
}

func (screen *Screen) Receive() {
	for {
		if err, messages := tpc.Receive(screen.Client.Conn); err != nil {
			// CreateMessages(screen.Windows[1], screen.Client.Conversation, messages)
			// response := messages[1:]
			switch messages[0] {
			case 250:
			case 251:
			case 252:
			case 253:
			default:
			}
		}
	}
}

package duckdom

import "fmt"

type Sidebar struct {
	ActiveChildId int
	Children      []Renderable
	Pos           Position
	Content       string
	Styles        string
}

func (self *Sidebar) OnClick() {}

func (self *Sidebar) Render() string {
	buffer := ""
	for i := range self.Pos.Row{
		for j := range self.Pos.Col{
			buffer += fmt.Sprintf(MOVE_CURSOR_TO_POSITION+"%s%s"+RESET_STYLES, i, j, DEBUG_STYLES, " ")
		}
	}

	return buffer
	// return fmt.Sprintf(MOVE_CURSOR_TO_POSITION+"%s%s"+RESET_STYLES, self.Pos.Row, self.Pos.Col, self.Styles, self.Content)
}

func (self *Sidebar) SetStyle(styles string) {
	self.Styles = styles
}

func (self *Sidebar) Active() Renderable { return self.Children[self.ActiveChildId] }
func (self *Sidebar) SetActive(id int)   { self.ActiveChildId = id }
func (self *Sidebar) ActiveIndex() int   { return self.ActiveChildId }

func (self *Sidebar) GetPos() Position { return self.Pos }

package duckdom

import "fmt"

type Button struct {
	ActiveChildId int
	Children      []Renderable
	Pos           Position
	Content       string
	Styles        string
}

func (self *Button) OnClick() {}

func (self *Button) Render() string {
	return fmt.Sprintf(MOVE_CURSOR_TO_POSITION+"%s%s"+RESET_STYLES, self.Pos.StartingRow, self.Pos.StartingCol, self.Styles, self.Content)
}
func (self *Button) SetStyle(styles string) {
	self.Styles = styles
}

func (self *Button) Active() Renderable { return self.Children[self.ActiveChildId] }
func (self *Button) SetActive(id int)   { self.ActiveChildId = id }
func (self *Button) ActiveIndex() int   { return self.ActiveChildId }

func (self *Button) GetPos() Position { return self.Pos }

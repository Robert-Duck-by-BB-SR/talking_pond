package duckdom

import "fmt"

type Button struct {
	Active_child_indx int
	Children          []Renderable
	Pos               Position
	Content           string
	Styles            string
}

func (self *Button) OnClick() {}

func (self *Button) Render() string {
	return fmt.Sprintf("\033[%d;%dH%s%s\033[0m", self.Pos.Row, self.Pos.Col, self.Styles, self.Content)
}

func (self *Button) SetPos(p Position) {
	self.Pos = p
}

func (self *Button) GetPos() Position {
	return self.Pos
}
func (self *Button) SetStyle(styles string) {
	self.Styles = styles
}

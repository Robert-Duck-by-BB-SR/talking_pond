package duckdom

import (
	"fmt"
	"strings"
)

type Window struct {
	ActiveChildId int
	Children      []Renderable
	Pos           Position
	Content       string
	Styles        Styles
}

func (self *Window) SetWidth(w int) Stylable{
	self.Styles.Width = w;
	return self
}

func (self *Window) SetHeight(h int) Stylable{
	self.Styles.Height = h
	return self
}

func (self *Window) SetBackground(b string) Stylable{
	self.Styles.Background = b
	return self
}

func (self *Window) Render() string {
	var string_builder strings.Builder
	string_builder.WriteString(self.Styles.Background)
	fillament := strings.Repeat(" ", self.Styles.Width) 
	for i := self.Pos.Row; i < uint(self.Styles.Height); i+= 1{
		// replace 0 by starting col
		string_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, i, self.Pos.Col))
		string_builder.WriteString(fillament)
	}
	string_builder.WriteString(RESET_STYLES)

	return string_builder.String()
}

func (self *Window) SetStyle(styles Styles) {
	self.Styles = styles
}

func (self *Window) Active() Renderable { return self.Children[self.ActiveChildId] }
func (self *Window) SetActive(id int)   { self.ActiveChildId = id }
func (self *Window) ActiveIndex() int   { return self.ActiveChildId }
func (self *Window) GetPos() Position { return self.Pos }

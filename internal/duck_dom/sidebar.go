package duckdom

import (
	"fmt"
	"strings"
)

type Sidebar struct {
	ActiveChildId int
	Children      []Renderable
	Pos           Position
	Content       string
	Styles        Styles
}

func (self *Sidebar) SetWidth(w int) Stylable{
	self.Styles.Width = w;
	return self
}

func (self *Sidebar) SetHeight(h int) Stylable{
	self.Styles.Height = h
	return self
}

func (self *Sidebar) SetBackground(b string) Stylable{
	self.Styles.Background = b
	return self
}

func (self *Sidebar) Render() string {
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

func (self *Sidebar) SetStyle(styles Styles) {
	self.Styles = styles
}

func (self *Sidebar) Active() Renderable { return self.Children[self.ActiveChildId] }
func (self *Sidebar) SetActive(id int)   { self.ActiveChildId = id }
func (self *Sidebar) ActiveIndex() int   { return self.ActiveChildId }
func (self *Sidebar) GetPos() Position { return self.Pos }

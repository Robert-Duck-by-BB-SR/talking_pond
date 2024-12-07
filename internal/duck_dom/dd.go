package duckdom

import "fmt"

type Screen struct{
	Cursor_pos [2]uint
	Active_child_indx uint
	Children []Item
}

type Item struct{
	Active_child_indx uint
	Children []Item
	Row uint
	Col uint
	Content string
	Styles string
}

func (self* Item) Render() string{
	return fmt.Sprintf("\r\033[%d;%dH%s%s\033[0m", self.Row, self.Col, self.Styles, self.Content)
}


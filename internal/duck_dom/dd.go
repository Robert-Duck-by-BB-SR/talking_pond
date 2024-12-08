package duckdom

import "fmt"

type Screen struct{
	Max_rows int
	Max_cols int
	Cursor_pos [2]uint
	Screen_active_child_indx int
	Children []Item
}

type Item struct{
	Active_child_indx int
	Children []Item
	Row uint
	Col uint
	Content string
	Styles string
	counter int
}

func (self* Item) Render() string{
	return fmt.Sprintf("\033[%d;%dH\033[2K%s%s\033[0m", self.Row, self.Col, self.Styles, self.Content)
}

func Debug_me_daddy(screen *Screen, content string){
	fmt.Printf("\033[%d;1H", screen.Max_rows)
	fmt.Printf("\033[2K")
	fmt.Printf("\033[%d;%dH\033[30;43m%s\033[0m", screen.Max_rows, 1, "DebugDuck: " + content)
}


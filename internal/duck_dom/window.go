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

	if(self.Styles.Border.Width != 0){
		return render_with_border(self, &string_builder)
	} 

	return render_default(self, &string_builder)
}

func render_default(self *Window, string_builder *strings.Builder) string{
	string_builder.WriteString(self.Styles.Background)
	fillament := strings.Repeat(" ", self.Styles.Width)

	for row := self.Pos.StartingRow; row < uint(self.Styles.Height); row += 1{
		string_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, self.Pos.StartingCol))
		string_builder.WriteString(fillament)
	}
	string_builder.WriteString(RESET_STYLES)

	return string_builder.String()
}

func render_with_border(self *Window, string_builder *strings.Builder) string{
	string_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Pos.StartingRow, self.Pos.StartingCol))
	
	// top border
	var top_bor_builder strings.Builder
	top_bor_builder.WriteString(self.Styles.Border.Color)
	top_bor_builder.WriteString(normalBorder.TopLeft)
	top_bor_builder.WriteString(strings.Repeat(normalBorder.Top, self.Styles.Width - 2))
	top_bor_builder.WriteString(normalBorder.TopRight)

	// bottom border
	var bot_bor_builder strings.Builder
	bot_bor_builder.WriteString(normalBorder.BottomLeft)
	bot_bor_builder.WriteString(strings.Repeat(normalBorder.Bottom, self.Styles.Width - 2))
	bot_bor_builder.WriteString(normalBorder.BottomRight)

	// general filament 
	var ver_bor_filament strings.Builder
	ver_bor_filament.WriteString(self.Styles.Border.Color)
	ver_bor_filament.WriteString(normalBorder.Left)
	ver_bor_filament.WriteString(self.Styles.Background)
	ver_bor_filament.WriteString(strings.Repeat(" ", self.Styles.Width - 2))
	ver_bor_filament.WriteString(self.Styles.Border.Color)
	ver_bor_filament.WriteString(normalBorder.Right)

	string_builder.WriteString(top_bor_builder.String())
	string_builder.WriteString(self.Styles.Background)

	for row := self.Pos.StartingRow + 1; row < uint(self.Styles.Height) - 1; row += 1{
		string_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, self.Pos.StartingCol))
		string_builder.WriteString(ver_bor_filament.String())
	}

	string_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Styles.Height - 1, self.Pos.StartingCol))
	string_builder.WriteString(self.Styles.Border.Color)
	string_builder.WriteString(bot_bor_builder.String())
	string_builder.WriteString(RESET_STYLES)

	// FIXME: GREY BIRDER
	return string_builder.String()
}

func (self *Window) SetStyle(styles Styles) {
	self.Styles = styles
}

func (self *Window) Active() Renderable { return self.Children[self.ActiveChildId] }
func (self *Window) SetActive(id int)   { self.ActiveChildId = id }
func (self *Window) ActiveIndex() int   { return self.ActiveChildId }
func (self *Window) GetPos() Position { return self.Pos }

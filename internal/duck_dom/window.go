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
	var bor_builder strings.Builder

	if(self.Styles.Border.Width != 0){
		return render_with_border(self, &bor_builder)
	} 

	return render_default(self, &bor_builder)
}

func render_default(self *Window, bor_builder *strings.Builder) string{
	bor_builder.WriteString(self.Styles.Background)
	fillament := strings.Repeat(" ", self.Styles.Width)

	for row := self.Pos.StartingRow; row < uint(self.Styles.Height); row += 1{
		bor_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, self.Pos.StartingCol))
		bor_builder.WriteString(fillament)
	}
	bor_builder.WriteString(RESET_STYLES)

	return bor_builder.String()
}

func render_with_border(self *Window, main_border_build *strings.Builder) string{
	main_border_build.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Pos.StartingRow, self.Pos.StartingCol))
	
	// top border
	var top_border_builder strings.Builder
	top_border_builder.WriteString(self.Styles.Border.Color)
	top_border_builder.WriteString(normalBorder.TopLeft)
	top_border_builder.WriteString(strings.Repeat(normalBorder.Top, self.Styles.Width - 2))
	top_border_builder.WriteString(normalBorder.TopRight)
	top_border_builder.WriteString(RESET_STYLES)

	// general filament + border 
	var vertical_border_filament strings.Builder
	vertical_border_filament.WriteString(self.Styles.Border.Color)
	vertical_border_filament.WriteString(normalBorder.Left)
	vertical_border_filament.WriteString(RESET_STYLES)
	vertical_border_filament.WriteString(self.Styles.Background)
	vertical_border_filament.WriteString(strings.Repeat(" ", self.Styles.Width - 2))
	vertical_border_filament.WriteString(RESET_STYLES)
	vertical_border_filament.WriteString(self.Styles.Border.Color)
	vertical_border_filament.WriteString(normalBorder.Right)
	vertical_border_filament.WriteString(RESET_STYLES)

	// bottom border
	var bottom_border_build strings.Builder
	bottom_border_build.WriteString(self.Styles.Border.Color)
	bottom_border_build.WriteString(normalBorder.BottomLeft)
	bottom_border_build.WriteString(strings.Repeat(normalBorder.Bottom, self.Styles.Width - 2))
	bottom_border_build.WriteString(self.Styles.Border.Color)
	bottom_border_build.WriteString(normalBorder.BottomRight)
	bottom_border_build.WriteString(RESET_STYLES)

	// combining my ass
	main_border_build.WriteString(top_border_builder.String())

	for row := self.Pos.StartingRow + 1; row < uint(self.Styles.Height) - 1; row += 1{
		main_border_build.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, self.Pos.StartingCol))
		main_border_build.WriteString(vertical_border_filament.String())
	}

	main_border_build.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Styles.Height - 1, self.Pos.StartingCol))
	main_border_build.WriteString(bottom_border_build.String())

	// FIXME: GREY BIRDER
	main_border_build.WriteString(RESET_STYLES)
	return main_border_build.String()
}

func (self *Window) SetStyle(styles Styles) {
	self.Styles = styles
}

func (self *Window) Active() Renderable { return self.Children[self.ActiveChildId] }
func (self *Window) SetActive(id int)   { self.ActiveChildId = id }
func (self *Window) ActiveIndex() int   { return self.ActiveChildId }
func (self *Window) GetPos() Position { return self.Pos }

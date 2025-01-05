package duckdom

import (
	"fmt"
	"strings"
)

type Component struct {
	Position
	Styles
	Content         string
	Buffer          string
	Parent          *Window
	ChildComponents []Component
	Active          bool
	Index			int
}

func CreateComponent(buffer string, styles Styles) *Component {
	// assert_component_dimensions(&styles)

	// // width auto
	// if styles.Width == 0 {
	// 	styles.Width = len(buffer)
	// }
	//
	// if styles.MaxWidth != 0 && len(buffer) > styles.MaxWidth {
	// 	styles.Width = styles.MaxWidth
	// }
	//
	// allowed_horizontal_space := self.Width - shift_cursor_by_border*2 - self.Styles.Paddding*2
	// allowed_vertical_space := self.Height - shift_cursor_by_border*2 - self.Styles.Paddding*2
	//
	// lines_used := 0
	// content := self.Buffer
	// for len(content) > allowed_horizontal_space {
	// 	// buffer_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))
	// 	// buffer_builder.WriteString(content[:allowed_horizontal_space])
	// 	content = content[allowed_horizontal_space:]
	// 	lines_used += 1
	// }

	component := Component{
		Buffer: buffer,
		Styles: styles,
	}

	return &component
}

// func assert_component_dimensions(styles *Styles) {
// 	if styles.Border.Style != NoBorder && styles.Width < 3 ||
// 		styles.Border.Style != NoBorder && styles.Height < 3 {
// 		panic("Component width and height should be at least 3 when border was added")
// 	}
//
// 	if styles.Width < 0 || styles.Height < 0 {
// 		panic("Component width and height should be bigger than -1")
// 	}
// }

func (self *Component) Render() string {
	self.rearrange_component()

	var builder strings.Builder
	if self.Active {
		builder.WriteString(INVERT_STYLES)
	} else {
		builder.WriteString(RESET_STYLES)
	}
	builder.WriteString(self.Styles.Compile())
	builder.WriteString(self.render_background())
	builder.WriteString(self.render_buffer())
	builder.WriteString(RESET_STYLES)

	// TODO: test me
	if self.Styles.Border.Style != NoBorder {
		builder.WriteString(render_border(self.Position, self.Active, &self.Styles))
	}

	self.Content = builder.String()
	return self.Content
}

func (self *Component) rearrange_component(){
	if self.Index == 0 {
		move_row_by := 0 
		move_col_by := 0 

		if self.Parent.Styles.Border.Style != NoBorder{
			move_col_by += 1
			move_row_by += 1
		}

		move_col_by += self.Parent.Styles.Paddding
		move_row_by += self.Parent.Styles.Paddding

		self.Position = Position{Row: self.Parent.Row + move_row_by, Col: self.Parent.Col + move_col_by}

	} else{
		siblings := self.Parent.Components
		if self.Styles.Direction == Block {
			last_component := siblings[len(siblings)-1]
			new_row := last_component.Row + last_component.Height
			new_col := last_component.Col

			DebugMeDaddy(self.Parent.Parent, string(new_col))

			// move to the next line in case it doesnt fit

			// rows_will_take := new_row + self.Styles.Height
			// cols_will_take := new_col + self.Width
			// assert_component_placement(rows_will_take, cols_will_take, c, self)

			self.Position = Position{Row: new_row, Col: new_col}
		} else {
			last_component := siblings[len(siblings)-1]
			new_row := last_component.Row
			new_col := last_component.Col + last_component.Width
			// assert_component_placement(new_row+self.Height, new_col+self.Width, c, self)

			self.Position = Position{Row: new_row, Col: new_col}
		}
		// arrange based on other siblings
	}

	// if len(siblings) == 1 {
	// 	if self.Border.Style != NoBorder {
	// 		self.Position = Position{StartingRow: self.StartingRow + 1, StartingCol: self.StartingCol + 1}
	// 	}
	// 	// assert_component_placement( self.StartingRow+self.Styles.Height, self.StartingCol+self.Width, self. self)
	// } else {
	// 	if self.Styles.Direction == Block {
	// 		last_component := siblings[len(siblings)-1]
	// 		new_row := last_component.StartingRow + last_component.Height
	// 		new_col := last_component.StartingCol
	//
	// 		// rows_will_take := new_row + self.Styles.Height
	// 		// cols_will_take := new_col + self.Width
	// 		// assert_component_placement(rows_will_take, cols_will_take, c, self)
	//
	// 		self.Position = Position{StartingRow: new_row, StartingCol: new_col}
	// 	} else {
	// 		last_component := siblings[len(siblings)-1]
	// 		new_row := last_component.StartingRow
	// 		new_col := last_component.StartingCol + last_component.Width
	// 		// assert_component_placement(new_row+self.Height, new_col+self.Width, c, self)
	//
	// 		self.Position = Position{StartingRow: new_row, StartingCol: new_col}
	// 	}
	// }
}

func (self *Component) render_background() string {
	var bg_builder strings.Builder
	fillament := strings.Repeat(" ", self.Styles.Width)

	for i := 0; i < self.Styles.Height; i += 1 {
		bg_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Row+i, self.Position.Col))
		bg_builder.WriteString(fillament)
	}

	return bg_builder.String()
}

func (self *Component) render_buffer() string {
	var buffer_builder strings.Builder

	shift_cursor_by_border := 0
	if self.Styles.Border.Style != NoBorder {
		shift_cursor_by_border += 1
	}

	moved_row := self.Row + shift_cursor_by_border + self.Styles.Paddding
	moved_col := self.Col + shift_cursor_by_border + self.Styles.Paddding
	// allowed_horizontal_space := self.Width - shift_cursor_by_border*2 - self.Styles.Paddding*2
	// allowed_vertical_space := self.Height - shift_cursor_by_border*2 - self.Styles.Paddding*2
	//
	// lines_used := 0
	// content := self.Buffer
	// for len(content) > allowed_horizontal_space {
	// 	// buffer_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))
	// 	// buffer_builder.WriteString(content[:allowed_horizontal_space])
	// 	content = content[allowed_horizontal_space:]
	// 	lines_used += 1
	// }

	buffer_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row, moved_col))
	buffer_builder.WriteString(self.Buffer)

	return buffer_builder.String()
}

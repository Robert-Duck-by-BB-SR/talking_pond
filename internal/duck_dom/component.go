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
	Index           int
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

func (self *Component) rearrange_component() {
	if self.Index == 0 {
		self.Row = self.Parent.Row
		self.Col = self.Parent.Col

		if self.Parent.Border.Style != NoBorder {
			self.Row += 1
			self.Col += 1
		}

		self.Row += self.Parent.Paddding
		self.Col += self.Parent.Paddding

		return
	}

	prev_component := self.Parent.Components[self.Index-1]
	if self.Parent.Direction == BLOCK {
		self.Row = prev_component.Row + prev_component.Height
		self.Col = prev_component.Col
	} else {
		self.Row = prev_component.Row
		self.Col = prev_component.Col + prev_component.Width
	}
}

func (self *Component) calculate_dimensions() {
	// width auto
	styles := self.Styles
	if styles.Width == 0 {
		styles.Width = len(self.Buffer)
	}
	if styles.MaxWidth != 0 && len(self.Buffer) > styles.MaxWidth {
		styles.Width = styles.MaxWidth
	}

	shift_cursor_by_border := 0
	if self.Styles.Border.Style != NoBorder {
		shift_cursor_by_border += 1
	}

	allowed_horizontal_space := self.Width - shift_cursor_by_border*2 - self.Styles.Paddding*2
	lines_used := 0
	content := self.Buffer
	for len(content) > allowed_horizontal_space {
		// buffer_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))
		// buffer_builder.WriteString(content[:allowed_horizontal_space])
		content = content[allowed_horizontal_space:]
		lines_used += 1
	}
	self.Height = lines_used
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

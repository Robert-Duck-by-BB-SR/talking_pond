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
	component := Component{
		Buffer: buffer,
		Styles: styles,
	}

	return &component
}


func (self *Component) Render() string {
	self.rearrange_component()
	self.calculate_dimensions()
	assert_component_dimensions(&self.Styles)

	var builder strings.Builder
	if self.Active {
		builder.WriteString(INVERT_STYLES)
	} else {
		builder.WriteString(RESET_STYLES)
	}
	builder.WriteString(self.Styles.Compile())
	builder.WriteString(self.render_background())
	builder.WriteString(self.Content)
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
	// component allows only block direction
	self.Row = prev_component.Row + prev_component.Height
	self.Col = prev_component.Col
}

// Changes dimentions of a component based on content
// Content will be updated if it does not fit into one line
func (self *Component) calculate_dimensions() {
	shift_cursor_by_border := 0
	if self.Styles.Border.Style != NoBorder {
		shift_cursor_by_border += 1
	}

	// width auto
	if self.Width == 0 {
		self.Width = len(self.Buffer)
		self.Width += self.Paddding*2 + shift_cursor_by_border * 2
	}

	if self.MaxWidth != 0 && len(self.Buffer) > self.MaxWidth {
		self.Width = self.MaxWidth
	}

	moved_row := self.Row + shift_cursor_by_border + self.Styles.Paddding
	moved_col := self.Col + shift_cursor_by_border + self.Styles.Paddding

	allowed_horizontal_space := self.Width - shift_cursor_by_border*2 - self.Styles.Paddding*2
	allowed_vertical_space := self.MaxHeight - shift_cursor_by_border*2 - self.Styles.Paddding*2
	lines_used := 0
	content := self.Buffer
	var content_builder strings.Builder
	for len(content) > allowed_horizontal_space {
		content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))
		content_builder.WriteString(content[:allowed_horizontal_space])
		content = content[allowed_horizontal_space:]
		lines_used += 1
		if allowed_vertical_space > 0 && lines_used > allowed_vertical_space {
			break
		}
	}
	content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))
	content_builder.WriteString(content)
	self.Content = content_builder.String()
	lines_used+= 1
	self.Height = lines_used + shift_cursor_by_border*2 + self.Styles.Paddding*2
}

func assert_component_dimensions(styles *Styles) {
	if styles.Border.Style != NoBorder && styles.Width < 3{
		panic("Component width should be at least 3 when border was added")
	}
	if styles.Border.Style != NoBorder && styles.Height < 3 {
		panic("Component height should be at least 3 when border was added")
	}

	if styles.Width < 0 || styles.Height < 0 {
		panic("Component width and height should be bigger than -1")
	}
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

package duckdom

import (
	"fmt"
	"strings"
)

type Component struct {
	Position
	Styles
	Content   string
	Buffer    string
	Parent    *Window
	Active    bool
	Inputable bool
	Scrollable bool
	BufferStartsFrom int
	Index     int
	Action    func()
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
	self.assert_component_dimensions()

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
	if self.Styles.Border != NoBorder {
		builder.WriteString(render_border(self.Position, self.Active, &self.Styles))
	}

	self.Content = builder.String()
	return self.Content
}

func (self *Component) rearrange_component() {
	if self.Index == 0 {
		self.Row = self.Parent.Row
		self.Col = self.Parent.Col

		if self.Parent.Border != NoBorder {
			self.Row += 1
			self.Col += 1
		}

		self.Row += self.Parent.Paddding
		self.Col += self.Parent.Paddding

		return
	}

	prev_component := self.Parent.Components[self.Index-1]
	if self.Direction == BLOCK {
		self.Row = prev_component.Row + prev_component.Height
		self.Col = prev_component.Col
	} else {
		self.Row = prev_component.Row
		self.Col = prev_component.Col + prev_component.Width
	}
}

// Changes dimentions of a component based on content
// Content will be updated if it does not fit into one line
func (self *Component) calculate_dimensions() {
	shift_cursor_by_border := 0

	if self.Styles.Border != NoBorder {
		shift_cursor_by_border += 1
	}

	if self.MaxWidth == 0 && self.Width == 0 {
		panic("MaxWidth and Width some of them at least should not be 0")
	}

	if self.MaxWidth != 0 && self.Width <= self.MaxWidth {
		self.Width = len(self.Buffer)
		if self.MinWidth > self.Width {
			self.Width = self.MinWidth
		}
		self.Width += self.Paddding*2 + shift_cursor_by_border*2
	}

	if self.MaxWidth != 0 && len(self.Buffer) > self.MaxWidth {
		self.Width = self.MaxWidth
	}

	if self.MaxHeight <= 0 {
		self.MaxHeight = self.Parent.Styles.Height - self.Row
	}
	moved_row := self.Row + shift_cursor_by_border + self.Styles.Paddding
	moved_col := self.Col + shift_cursor_by_border + self.Styles.Paddding

	allowed_horizontal_space := self.Width - shift_cursor_by_border*2 - self.Styles.Paddding*2
	if allowed_horizontal_space <= 0{
		panic(fmt.Sprintf("Allowed horizontal space should be bigger than 0, actual value: %d", allowed_horizontal_space))
	}

	allowed_vertical_space := self.MaxHeight - shift_cursor_by_border*2 - self.Styles.Paddding*2
	if self.MaxHeight > 0 && allowed_vertical_space <= 0{
		panic(fmt.Sprintf("Allowed vertical space should be bigger than 0, actual value: %d", allowed_vertical_space))
	}
	lines_used := 0
	content := self.Buffer 
	limit := len(content)/allowed_horizontal_space
	if self.BufferStartsFrom > limit {
		self.BufferStartsFrom = limit
	}

	if self.BufferStartsFrom < 0 {
		self.BufferStartsFrom = 0
	}
	content = content[allowed_horizontal_space*self.BufferStartsFrom:]
	var content_builder strings.Builder
	for len(content) > allowed_horizontal_space {
		content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))
		content_builder.WriteString(content[:allowed_horizontal_space])
		content = content[allowed_horizontal_space:]
		lines_used += 1
		if allowed_vertical_space > 0 && lines_used >= allowed_vertical_space {
			break
		}
	}

	if len(content) <= allowed_horizontal_space && lines_used < allowed_vertical_space{
		content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))
		content_builder.WriteString(content)
		lines_used += 1
	}

	self.Content = content_builder.String()

	if self.Height == 0 {
		self.Height = lines_used + shift_cursor_by_border*2 + self.Styles.Paddding*2
	}
}

func (self *Component) assert_component_dimensions() {
	if self.Border != NoBorder && self.Width < 3 {
		panic(fmt.Sprintf("Window %d component %d width should be at least 3 when border was added", self.Parent.Index, self.Index))
	}
	if self.Border != NoBorder && self.Height < 3 {
		panic(fmt.Sprintf("Window %d component %d height should be at least 3 when border was added", self.Parent.Index, self.Index))
	}

	if self.Width < 0 || self.Height < 0 {
		panic(fmt.Sprintf("Window %d component %d width and height should be bigger than -1", self.Parent.Index, self.Index))
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

package duckdom

import (
	"fmt"
	"strings"
)

type Component struct {
	Position
	Styles
	content                       string
	Buffer                        string
	Parent                        *Window
	Active                        bool
	Inputable                     bool
	ScrollType                    ScrollType
	buffer_vertical_scroll_from   int
	buffer_horizontal_scroll_from int
	Index                         int
	allowed_vertical_space        int
	allowed_horizontal_space      int
	Action                        func()
}

type ScrollType int

const (
	NONE ScrollType = iota
	VERTICAL
	HORIZONTAL
)

func CreateComponent(buffer string, styles Styles) *Component {
	component := Component{
		Buffer: buffer,
		Styles: styles,
	}

	return &component
}

func (self *Component) reverse_rearange() string {
	var builder strings.Builder
	defer builder.Reset()

	parent := self.Parent

	// NOTE: supporting display block for now only

	if self.Index == parent.scroll_to {
		self.Row = parent.Row + parent.Height - self.Height
		self.Col = parent.Col

		if parent.Border != NoBorder {
			self.Row -= 1
			self.Col += 1
		}

		self.Row -= parent.Paddding
		self.Col += parent.Paddding

		return builder.String()
	}

	if self.Index < len(parent.Components)-1 {
		next := parent.Components[self.Index+1]
		self.Row = next.Row - self.Height
		self.Col = next.Col
		return builder.String()
	}

	return ""
}

func (self *Component) Render() string {
	var builder strings.Builder
	defer builder.Reset()
	if !self.Parent.ReverseRenderable {
		self.rearrange_component()
		self.calculate_dimensions()
	} else {
		self.calculate_dimensions()
		self.reverse_rearange()
	}
	self.assert_component_dimensions()
	builder.WriteString(self.render_content())

	if self.Styles.Border != NoBorder {
		render_border(&builder, self.Position, self.Active, &self.Styles)
	}
	builder.WriteString(RESET_STYLES)
	return builder.String()
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

func (self *Component) render_background(content_builder *strings.Builder) {
	full_line_fillament := strings.Repeat(" ", self.Width)
	for row := self.Row; row < self.Row+self.Height; row += 1 {
		position := fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, self.Col)
		content_builder.WriteString(position)
		content_builder.WriteString(full_line_fillament)
	}
}

// Changes dimentions of a component based on content
func (self *Component) calculate_dimensions() {
	if self.Width == 0 && self.Inputable {
		panic("Inputable components must have static width")
	}
	if self.Height == 0 && self.Inputable {
		panic("Inputable components must have static height")
	}

	if self.allowed_horizontal_space != 0 && self.allowed_vertical_space != 0 {
		return
	}

	if self.MaxWidth == 0 && self.Width == 0 {
		panic("MaxWidth and Width some of them at least should not be 0")
	}

	shift_cursor_by_border := 0

	if self.Styles.Border != NoBorder {
		shift_cursor_by_border += 1
	}

	// Width auto
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

	// NOTE: auto max height???
	if self.MaxHeight <= 0 {
		self.MaxHeight = self.Parent.Row + self.Parent.Styles.Height - self.Row
	}

	self.allowed_horizontal_space = self.Width - shift_cursor_by_border*2 - self.Styles.Paddding*2
	if self.allowed_horizontal_space <= 0 {
		panic(fmt.Sprintf(
			"[W|C][%d|%d] Allowed horizontal space should be bigger than 0, actual value: %d",
			self.Parent.Index, self.Index, self.allowed_horizontal_space,
		))
	}

	max_height := self.MaxHeight
	if self.Height != 0 {
		max_height = self.Height
	}

	self.allowed_vertical_space = max_height - shift_cursor_by_border*2 - self.Styles.Paddding*2
	if self.MaxHeight > 0 && self.allowed_vertical_space <= 0 {
		panic(fmt.Sprintf("[W|C][%d|%d] Allowed vertical space should be bigger than 0, actual value: %d",
			self.Parent.Index, self.Index, self.allowed_vertical_space,
		))
	}

	content := self.Buffer
	lines_used := 0

	for len(content) > self.allowed_horizontal_space {
		content = content[self.allowed_horizontal_space:]
		lines_used += 1
		if self.allowed_vertical_space > 0 && lines_used >= self.allowed_vertical_space {
			break
		}
	}

	if len(content) <= self.allowed_horizontal_space && lines_used < self.allowed_vertical_space {
		lines_used += 1
	}

	if self.Height == 0 {
		self.Height = lines_used + shift_cursor_by_border*2 + self.Styles.Paddding*2
	}

}

// Content will be updated if it does not fit into one line
func (self *Component) render_buffer(content_builder *strings.Builder) {
	if self.allowed_horizontal_space == 0 && self.allowed_vertical_space == 0 {
		panic("you forgot to calculate dimensions you fucking donkey")
	}
	shift_cursor_by_border := 0

	if self.Styles.Border != NoBorder {
		shift_cursor_by_border += 1
	}

	moved_row := self.Row + shift_cursor_by_border + self.Styles.Paddding
	moved_col := self.Col + shift_cursor_by_border + self.Styles.Paddding
	if self.buffer_vertical_scroll_from < 0 {
		self.buffer_vertical_scroll_from = 0
	}
	if self.buffer_horizontal_scroll_from < 0 {
		self.buffer_horizontal_scroll_from = 0
	}

	content := self.Buffer
	vertical_scroll_limit := len(content)/self.allowed_horizontal_space - self.allowed_vertical_space + 1

	if self.buffer_vertical_scroll_from > vertical_scroll_limit {
		self.buffer_vertical_scroll_from = vertical_scroll_limit
	}
	if self.ScrollType == VERTICAL && self.allowed_vertical_space > 0 && self.buffer_vertical_scroll_from > 0 {
		content = content[self.allowed_horizontal_space*self.buffer_vertical_scroll_from:]
	}

	if self.allowed_horizontal_space+self.buffer_horizontal_scroll_from > len(content) {
		self.buffer_horizontal_scroll_from -= 1
	}
	if self.ScrollType == HORIZONTAL && self.allowed_vertical_space == 1 {
		content = content[self.buffer_horizontal_scroll_from:]
	}

	lines_used := 0
	for len(content) > self.allowed_horizontal_space {
		content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))

		content_builder.WriteString(content[:self.allowed_horizontal_space])
		content = content[self.allowed_horizontal_space:]
		lines_used += 1
		if self.allowed_vertical_space > 0 && lines_used >= self.allowed_vertical_space {
			break
		}
	}
	if len(content) <= self.allowed_horizontal_space && lines_used < self.allowed_vertical_space {
		content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))
		content_builder.WriteString(content)
		lines_used += 1
	}
}

func (self *Component) render_content() string {
	var content_builder strings.Builder
	defer content_builder.Reset()
	if self.Active {
		content_builder.WriteString(INVERT_STYLES)
	} else {
		content_builder.WriteString(RESET_STYLES)
	}
	self.Styles.Compile(&content_builder)
	self.render_background(&content_builder)
	self.render_buffer(&content_builder)
	return content_builder.String()
}

func (self *Component) assert_component_dimensions() {
	if self.Border != NoBorder && self.Width < 3 {
		panic(fmt.Sprintf("Window %d component %d width should be at least 3 when border was added: %+v", self.Parent.Index, self.Index, self))
	}
	if self.Border != NoBorder && self.Height < 3 {
		panic(fmt.Sprintf("Window %d component %d height should be at least 3 when border was added: %+v", self.Parent.Index, self.Index, self))
	}

	if self.Width < 0 || self.Height < 0 {
		panic(fmt.Sprintf("Window %d component %d width and height should be bigger than -1: %+v", self.Parent.Index, self.Index, self))
	}
}

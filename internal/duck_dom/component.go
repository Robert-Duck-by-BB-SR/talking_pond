package duckdom

import (
	"fmt"
	"strings"
)

type Component struct {
	Position
	Styles
	Content              string
	Buffer               string
	Parent               *Window
	Active               bool
	Inputable            bool
	ScrollType           ScrollType
	BufferVerticalFrom   int
	BufferHorizontalFrom int
	Index                int
	Action               func()
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

func (self *Component) Render(builder *strings.Builder) {
	self.rearrange_component()
	if self.Active {
		builder.WriteString(INVERT_STYLES)
	} else {
		builder.WriteString(RESET_STYLES)
	}
	self.Styles.Compile(builder)
	self.render_background(builder)
	self.calculate_dimensions(builder)
	self.assert_component_dimensions()

	if self.Styles.Border != NoBorder {
		render_border(builder, self.Position, self.Active, &self.Styles)
	}
	builder.WriteString(RESET_STYLES)
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
func (self *Component) calculate_dimensions(content_builder *strings.Builder) {
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
		self.MaxHeight = self.Parent.Row + self.Parent.Styles.Height - self.Row
	}
	moved_row := self.Row + shift_cursor_by_border + self.Styles.Paddding
	moved_col := self.Col + shift_cursor_by_border + self.Styles.Paddding

	allowed_horizontal_space := self.Width - shift_cursor_by_border*2 - self.Styles.Paddding*2
	if allowed_horizontal_space <= 0 {
		panic(fmt.Sprintf("[W|C][%d|%d]Allowed horizontal space should be bigger than 0, actual value: %d", self.Parent.Index, self.Index, allowed_horizontal_space))
	}

	allowed_vertical_space := self.MaxHeight - shift_cursor_by_border*2 - self.Styles.Paddding*2
	if self.MaxHeight > 0 && allowed_vertical_space <= 0 {
		panic(fmt.Sprintf("[W|C][%d|%d]Allowed vertical space should be bigger than 0, actual value: %d", self.Parent.Index, self.Index, allowed_vertical_space))
	}

	if self.BufferVerticalFrom < 0 {
		self.BufferVerticalFrom = 0
	}
	if self.BufferHorizontalFrom < 0 {
		self.BufferHorizontalFrom = 0
	}

	content := self.Buffer
	vertical_scroll_limit := len(content) / allowed_horizontal_space
	if self.BufferVerticalFrom > vertical_scroll_limit {
		self.BufferVerticalFrom = vertical_scroll_limit
	}
	if self.ScrollType == VERTICAL {
		content = content[allowed_horizontal_space*self.BufferVerticalFrom:]
	}

	if allowed_horizontal_space+self.BufferHorizontalFrom > len(content) {
		self.BufferHorizontalFrom -= 1
	}
	if self.ScrollType == HORIZONTAL && allowed_vertical_space == 1 {
		content = content[self.BufferHorizontalFrom:]
	}

	//######### - row col + width
	//###123### - movedrow and movedcol = bg * padding - col, bg * padding - col + content line len
	//######### - row + height and col + width

	full_line_fillament := strings.Repeat(" ", self.Styles.Width)
	FileDebugMeDaddy(fmt.Sprintln(self.Styles.Width))
	if self.Paddding != 0 {
		for i := range self.Paddding {
			content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Row+shift_cursor_by_border+i, self.Col))
			content_builder.WriteString(full_line_fillament)
		}
	}

	if self.Buffer == "" && self.MinWidth != 0 {
		content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Row+shift_cursor_by_border+self.Paddding, self.Col))
		content_builder.WriteString(full_line_fillament)
	}

	lines_used := 0
	for len(content) > allowed_horizontal_space {
		if self.Paddding != 0 {
			content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, self.Col))
			content_builder.WriteString(full_line_fillament)
		}
		content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))

		content_builder.WriteString(content[:allowed_horizontal_space])
		content = content[allowed_horizontal_space:]
		lines_used += 1
		if allowed_vertical_space > 0 && lines_used >= allowed_vertical_space {
			break
		}
	}

	if len(content) <= allowed_horizontal_space && lines_used < allowed_vertical_space {
		if self.Paddding != 0 {
			content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, self.Col))
			content_builder.WriteString(full_line_fillament)
		}
		content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+lines_used, moved_col))
		content_builder.WriteString(content)
		lines_used += 1
	}

	if self.Height == 0 {
		self.Height = lines_used + shift_cursor_by_border*2 + self.Styles.Paddding*2
	}

	if self.Paddding != 0 {
		for i := range self.Paddding {
			content_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Row+self.Height-shift_cursor_by_border*2-i, self.Col))
			content_builder.WriteString(full_line_fillament)
		}
	}
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

func (self *Component) render_background(bg_builder *strings.Builder) {
}

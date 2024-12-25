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
	// NOTE: we should really think about it
	// maybe it would be better if we just made a bunch of functions
	// that take *Component as an input and does some actions with it
	Action func()
}

func CreateComponent(buffer string, styles Styles) *Component {
	// check styles
	// change width of component in case if buffer is bigger than provided

	assert_component_dimensions(&styles)

	if styles.Border.Style != NoBorder {
		// if we use border, min width and height should be 2
		// why tho?? your assert says it has to be bigger than 3????
		if styles.Width < 3 {
			styles.Width = 3
		}

		if styles.Height < 3 {
			styles.Height = 3
		}

	}

	component := Component{
		Buffer: buffer,
		Styles: styles,
	}

	return &component
}

func assert_component_dimensions(styles *Styles) {
	if styles.Border.Style != NoBorder && styles.Width < 3 || styles.Border.Style != NoBorder && styles.Height < 3 {
		panic("Component width and height should be at least 3 when border was added")
	}

	if styles.Width < 1 || styles.Height < 1 {
		panic("Component width and height should be bigger than 0")
	}
}

func (self *Component) ExecuteAction() {
	self.Action()
}

func (self *Component) Render() string {
	// TODO: test me
	// later somewhere here I will implement
	// 1. padding -> deez nuts
	// 2. text-align

	var builder strings.Builder
	builder.WriteString(self.Styles.Compile())
	builder.WriteString(self.render_background())
	builder.WriteString(self.render_buffer())
	builder.WriteString(RESET_STYLES)
	self.Content = builder.String()

	// TODO: test me
	if self.Styles.Border.Style != NoBorder {
		self.Content += render_border(self.Position, &self.Styles)
	}

	return self.Content
}

func (self *Component) render_background() string {
	var bg_builder strings.Builder
	fillament := strings.Repeat(" ", self.Styles.Width)

	for i := 0; i < self.Styles.Height; i += 1 {
		bg_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.StartingRow+i, self.Position.StartingCol))
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

	moved_row := self.StartingRow + shift_cursor_by_border + self.Styles.Paddding
	moved_col := self.StartingCol + shift_cursor_by_border + self.Styles.Paddding
	allowed_space := self.Width - shift_cursor_by_border*2 - self.Styles.Paddding*2

	// if the whole word cannot fit in one line -> truncate
	if len(self.Buffer) > allowed_space {
		// splited_buffer := strings.Split(self.Buffer, " ")
		// for i := 0; i < len(splited_buffer); i += 1 {
		// 	// do something in case that single word is still too big
		// 	buffer_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row+i, moved_col))
		// 	buffer_builder.WriteString(splited_buffer[i])
		// }
		//
		// NOTE: proposal

		buffer_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row, moved_col))
		buffer_builder.WriteString(self.Buffer[:allowed_space])
        // this way we actually truncate by the character and not the space so we don't have to worry about moving to 
        // next line, FIXME: HOWEVER we should consider the fact that we should be able to increase the height of the component
        // depending on the content
        // the way it's done now any message in chat what will be greater than the width of the component would be 
        // truncated which is not the desired behaviour for a chat
	} else {
		buffer_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, moved_row, moved_col))
		buffer_builder.WriteString(self.Buffer)
	}

	return buffer_builder.String()
}

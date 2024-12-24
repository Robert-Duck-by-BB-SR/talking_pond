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

	assert_component_dimentions(&styles)

	if styles.Border.Style != NoBorder {
		// right now I'm concerned about it, future me will be mad
		// if we use border, min width and height should be 2
		switch {
		case styles.Height < 2:
			styles.Height += 2
		case styles.Height < 2:
			styles.Height += 1
		}

		switch {
		case styles.Width < 2:
			styles.Height += 2
		case styles.Width < 2:
			styles.Height += 1
		}

		// if len(buffer) 
	}


	component := Component{
		Buffer: buffer,
		Styles: styles,
	}

	return &component
}

func assert_component_dimentions(styles *Styles) {
	if styles.Border.Style != NoBorder && styles.Width < 3 || styles.Border.Style != NoBorder && styles.Height < 3{
		panic("Component width and height should be bigger than 3 when border was added")
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
	// 1. padding 
	// 2. text-align

	// TODO: test me 
	// We use +1 because border takes one char around

	// OKEY +1 MEANS THAT I MOVE BY ONE BUT BORDER TAKES 2 CHARS FROM BOTH SIDES
	// we should have two chars free from both sides (top bottom, left right)

	var builder strings.Builder
	builder.WriteString(self.Styles.Compile())

	shift_cursor_by_border := 0
	if self.Styles.Border.Style != NoBorder {
		shift_cursor_by_border += 1
	}

	builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Position.StartingRow + shift_cursor_by_border, self.Position.StartingCol + shift_cursor_by_border))


	fillament := strings.Repeat(" ", self.Styles.Width - shift_cursor_by_border)

	fillament_minus_text := ""
	// space (padding)
	// text (might be truncated)
	// space (padding + lefovers)
	if self.Styles.Paddding != 0 {
		// some bullshit here
	} else{
		fillament_minus_text = self.Buffer + strings.Repeat(" ", self.Styles.Width - len(self.Buffer) - shift_cursor_by_border)
	}
	
	for i := shift_cursor_by_border; i < self.Styles.Height; i += 1 {
		if i == shift_cursor_by_border{
			builder.WriteString(fillament_minus_text)
		} else{
			builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.StartingRow+i, self.Position.StartingCol))
			builder.WriteString(fillament)
		}
	}
	builder.WriteString(RESET_STYLES)
	self.Content = builder.String()

	// TODO: test me 
	if self.Styles.Border.Style != NoBorder {
		self.Content += render_border(self.Position, &self.Styles)
	}
	
	return self.Content
}


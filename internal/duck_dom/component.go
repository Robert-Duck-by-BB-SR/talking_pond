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

	assert_component_dimentions(styles.Width, styles.Height)

	if styles.Border.Style != NoBorder {
		// right now I'm concerned about it, future me will be mad
		switch {
		case styles.Height < 2:
			styles.Height += 2
		case styles.Height < 3:
			styles.Height += 1
		}

		switch {
		case styles.Width < 2:
			styles.Height += 2
		case styles.Width < 3:
			styles.Height += 1
		}
	}

	component := Component{
		Buffer: buffer,
		Styles: styles,
	}

	return &component
}

func assert_component_dimentions(w, h int) {
	if w <= 0 || h <= 0 {
		panic("Component width and height should be bigger than 0")
	}
}

func (self *Component) ExecuteAction() {
	self.Action()
}

func (self *Component) Render() string {
	// component := self.render_background()
	//
	// if self.Styles.Border.Style != NoBorder {
	// 	component += render_border(self.Position, &self.Styles)
	// }

	// later somewhere here I will implement
	// 1. padding 
	// 2. text-align




	// return component
	var builder strings.Builder
	// WE USE + 1 because border takes one char around
	// NEEDS TO BE UPDATED
	builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Position.StartingRow + 1, self.Position.StartingCol + 1))
	builder.WriteString(self.Styles.Compile())
	builder.WriteString(self.Buffer)
	builder.WriteString(strings.Repeat(" ", self.Styles.Width - len(self.Buffer)))
	builder.WriteString(RESET_STYLES)
	self.Content = builder.String()
	return self.Content + render_border(self.Position, &self.Styles)
}

func (self *Component) render_background() string {
	var bg_builder strings.Builder
	bg_builder.WriteString(self.Styles.Background)
	bg_builder.WriteString(self.Styles.Color)
	fillament := strings.Repeat(" ", self.Styles.Width)

	for i := 0; i < self.Styles.Height; i += 1 {
		bg_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.StartingRow+i, self.Position.StartingCol))
		bg_builder.WriteString(fillament)
	}
	bg_builder.WriteString(RESET_STYLES)

	return bg_builder.String()
}

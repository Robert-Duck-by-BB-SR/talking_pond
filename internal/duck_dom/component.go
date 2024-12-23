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

	if styles.Border.Style != NoBorder{
		// right now I'm concerned about it, future me will be mad
		if styles.Height < 2 {
			styles.Height += 2
		}
		if styles.Height < 3 {
			styles.Height += 1
		}

		if styles.Width < 2 {
			styles.Width += 2
		}
		if styles.Width < 3 {
			styles.Width += 1
		}
	}

	component := Component{
		Buffer: buffer,
		Styles: styles,
	}

	// &dd.Component{
	// 	// Position: dd.Position{StartingRow: 3, StartingCol: uint(sidebar.StartingCol) + 2},
	// 	Buffer: "|Deez nuts|",
	// 	Styles: dd.Styles{
	// 		Width: len("|Deez nuts|"),
	// 		// Width:      screen.Width - sidebar.Styles.Width - 1,
	// 		Height:     1,
	// 		Background: dd.MakeRGBBackground(80, 40, 100),
	// 		// Border: dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
	// 	},
	// }

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

	// if self.Styles.Border.Style != NoBorder {
	// 	component += render_border(self.Position, &self.Styles)
	// }

	// return component
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Position.StartingRow, self.Position.StartingCol))
	builder.WriteString(self.Styles.Compile())
	builder.WriteString(self.Buffer)
	builder.WriteString(RESET_STYLES)
	self.Content = builder.String()
	return self.Content
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


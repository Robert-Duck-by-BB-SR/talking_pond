package duckdom

import (
	"fmt"
	"strings"
)

type Window struct {
	Position
	Parent *Screen
	Styles
	ActiveComponentId int
	Components        []*Component
}

func CreateWindow(styles Styles) *Window {
	assert_window_dimensions(styles.Width, styles.Height)

	if styles.Border.Style != NoBorder {
		if styles.Width < 3 {
			styles.Width = 3
		}

		if styles.Height < 3 {
			styles.Height = 3
		}
	}

	return &Window{
		Position: Position{StartingRow: 1, StartingCol: 1},
		Styles:   styles,
	}
}

func assert_window_dimensions(w, h int) {
	if w <= 0 || h <= 0 {
		panic("Window width and height should be bigger than 0")
	}
}

func (self *Window) Render() string {
	var window_with_components strings.Builder
	window_with_components.WriteString(self.render_background())

	if self.Styles.Border.Style != NoBorder {
		window_with_components.WriteString(render_border(self.Position, &self.Styles))
	}

	for _, component := range self.Components {
		window_with_components.WriteString(component.Render())
	}

	return window_with_components.String()
}

func (self *Window) AddComponent(c *Component) {
	c.Parent = self

	if len(self.Components) == 0 {
		if self.Border.Style != NoBorder {
			c.Position = Position{StartingRow: self.StartingRow + 1, StartingCol: self.StartingCol + 1}
		}
		assert_component_placement(c.StartingRow+c.Styles.Height, c.StartingCol+c.Width, self)
	} else {
		if c.Styles.Direction == Block {
			last_component := self.Components[len(self.Components)-1]
			new_row := last_component.StartingRow + last_component.Height
			new_col := last_component.StartingCol
			assert_component_placement(new_row+c.Styles.Height, new_col+c.Width, self)

			c.Position = Position{StartingRow: new_row, StartingCol: new_col}
		} else {
			last_component := self.Components[len(self.Components)-1]
			new_row := last_component.StartingRow
			new_col := last_component.StartingCol + last_component.Width
			assert_component_placement(new_row+c.Height, new_col+c.Width, self)

			c.Position = Position{StartingRow: new_row, StartingCol: new_col}
		}
	}

	self.Components = append(self.Components, c)
}

func (self *Window) render_background() string {
	var bg_builder strings.Builder
	bg_builder.WriteString(self.Styles.Background)
	fillament := strings.Repeat(" ", self.Styles.Width)

	for i := 0; i < self.Styles.Height; i += 1 {
		bg_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.StartingRow+i, self.Position.StartingCol))
		bg_builder.WriteString(fillament)
	}
	bg_builder.WriteString(RESET_STYLES)

	return bg_builder.String()
}

func assert_component_placement(r, c int, w *Window) {
	if r > w.Height {
		panic("Component height will not fit, do a math, you dumbass")
	}

	if c > w.Width {
		panic("Component width will not fit, do a math, you dumbass")
	}
}

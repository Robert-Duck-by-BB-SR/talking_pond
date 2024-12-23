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

func CreateWindow(styles Styles) *Window{
	assert_window_dimentions(styles.Width, styles.Height)

	if styles.Border.Style != NoBorder{
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

	return &Window{
		Position: Position{StartingRow: 1, StartingCol: 1},
		Styles: styles,
	}
}

func assert_window_dimentions(w, h int) {
	if w <= 0 || h <= 0 {
		panic("Window width and height should be bigger than 0")
	}
}

func (self *Window) Render() string {
	window_with_components := self.render_background()

	if self.Styles.Border.Style != NoBorder {
		window_with_components += render_border(self.Position, &self.Styles)
	}

	for _, component := range self.Components {
		window_with_components += component.Render()
	}

	return window_with_components
}

func (self *Window) AddComponent(c *Component) {
	c.Parent = self

	if len(self.Components) == 0 {
		if self.Border.Style == NoBorder {
			c.Position = Position{StartingRow: 1, StartingCol: 1}
		} else {
			c.Position = Position{StartingRow: 2, StartingCol: 2}
		}
		assert_component_placement(c.StartingRow + c.Styles.Height, c.StartingCol + c.Width, self)
	} else {
		if c.Styles.Direction == Block {
			last_component := self.Components[len(self.Components)-1]
			lastRow := last_component.StartingRow + last_component.Height - 1
			// can be 1 or 2
			lastCol := last_component.StartingCol
			newRow := lastRow + c.Height
			assert_component_placement(newRow, lastCol + c.Width, self)

			c.Position = Position{StartingRow: newRow, StartingCol: lastCol}
		} else {
			last_component := self.Components[len(self.Components)-1]
			lastRow := last_component.StartingRow
			lastCol := last_component.StartingCol + last_component.Width
			assert_component_placement(lastRow + c.Height, lastCol + c.Width, self)

			c.Position = Position{StartingRow: lastRow, StartingCol: lastCol}
		}
	}

	self.Components = append(self.Components, c)
}

func (self *Window) render_background() string {
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

func assert_component_placement(r, c int, w *Window) {
	if r > w.Height {
		panic("Component height will not fit, do a math, you dumbass")
	}

	if c > w.Width {
		panic("Component width will not fit, do a math, you dumbass")
	}
}

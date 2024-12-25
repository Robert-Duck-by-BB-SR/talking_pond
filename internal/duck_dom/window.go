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
	assert_window_dimensions(&styles)

	return &Window{
		Position: Position{StartingRow: 1, StartingCol: 1},
		Styles:   styles,
	}
}

func assert_window_dimensions(styles *Styles) {
	if styles.Border.Style != NoBorder && styles.Width < 3 ||
		styles.Border.Style != NoBorder && styles.Height < 3 {
		panic("Component width and height should be at least 3 when border was added")
	}

	if styles.Width < 1 || styles.Height < 1 {
		panic("Component width and height should be bigger than 0")
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
		assert_component_placement(c.StartingRow+c.Styles.Height, c.StartingCol+c.Width, c, self)
	} else {
		if c.Styles.Direction == Block {
			last_component := self.Components[len(self.Components)-1]
			new_row := last_component.StartingRow + last_component.Height
			new_col := last_component.StartingCol

			rows_will_take := new_row+c.Styles.Height
			cols_will_take := new_col+c.Width
			assert_component_placement(rows_will_take, cols_will_take, c, self)

			c.Position = Position{StartingRow: new_row, StartingCol: new_col}
		} else {
			last_component := self.Components[len(self.Components)-1]
			new_row := last_component.StartingRow
			new_col := last_component.StartingCol + last_component.Width
			assert_component_placement(new_row+c.Height, new_col+c.Width, c, self)

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

func assert_component_placement(rows_will_take, cols_will_take int, c *Component, w *Window) {
	if cols_will_take > w.Width {
		panic(fmt.Sprintf(
			"Component width is too big: window [r,c][%d,%d] [w,h][%d,%d] will not fit component [r,c][%d,%d] [w,h][%d,%d]",
			w.StartingRow, w.StartingCol, w.Styles.Width, w.Styles.Height, c.StartingRow, c.StartingCol, c.Styles.Width, c.Styles.Height,
		))
	}
	if rows_will_take > w.Height {
		panic(fmt.Sprintf(
			"Component height is too big: window [r,c][%d,%d] [w,h][%d,%d] will not fit component [r,c][%d,%d] [w,h][%d,%d]",
			w.StartingRow, w.StartingCol, w.Styles.Width, w.Styles.Height, c.StartingRow, c.StartingCol, c.Styles.Width, c.Styles.Height,
		))
	}
}

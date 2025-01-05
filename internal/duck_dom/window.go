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
	Active            bool
}

func CreateWindow(styles Styles) *Window {
	assert_window_dimensions(&styles)

	return &Window{
		Position: Position{Row: 1, Col: 1},
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
		window_with_components.WriteString(render_border(self.Position, self.Active, &self.Styles))
	}

	for _, component := range self.Components {
		window_with_components.WriteString(component.Render())
	}

	return window_with_components.String()
}

func (self *Window) AddComponent(c *Component) {
	c.Index = len(self.Components)
	c.Parent = self
	self.Components = append(self.Components, c)
}

func assert_component_placement(rows_will_take, cols_will_take int, c *Component, w *Window) {
	if cols_will_take > w.Width {
		panic(fmt.Sprintf(
			"Component width is too big: window [r,c][%d,%d] [w,h][%d,%d] will not fit component [r,c][%d,%d] [w,h][%d,%d]",
			w.Row, w.Col, w.Styles.Width, w.Styles.Height, c.Row, c.Col, c.Styles.Width, c.Styles.Height,
		))
	}
	if rows_will_take > w.Height {
		panic(fmt.Sprintf(
			"Component height is too big: window [r,c][%d,%d] [w,h][%d,%d] will not fit component [r,c][%d,%d] [w,h][%d,%d]",
			w.Row, w.Col, w.Styles.Width, w.Styles.Height, c.Row, c.Col, c.Styles.Width, c.Styles.Height,
		))
	}
}

func (self *Window) render_background() string {
	var bg_builder strings.Builder
	bg_builder.WriteString(self.Styles.Background)
	fillament := strings.Repeat(" ", self.Styles.Width)

	for i := 0; i < self.Styles.Height; i += 1 {
		bg_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Row+i, self.Position.Col))
		bg_builder.WriteString(fillament)
	}
	bg_builder.WriteString(RESET_STYLES)

	return bg_builder.String()
}

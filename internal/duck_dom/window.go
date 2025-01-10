package duckdom

import (
	"fmt"
	"strings"
)

type Window struct {
	Position
	Styles
	Oldfart           *Screen
	ActiveComponentId int
	Components        []*Component
	Active            bool
	Index             int
}

func CreateWindow(styles Styles) *Window {
	assert_window_dimensions(&styles)

	return &Window{
		Styles: styles,
	}
}

func assert_window_dimensions(styles *Styles) {
	if styles.Border != NoBorder && styles.Width < 3 ||
		styles.Border != NoBorder && styles.Height < 3 {
		panic("Component width and height should be at least 3 when border was added")
	}

	if styles.Width < 1 || styles.Height < 1 {
		panic("Component width and height should be bigger than 0")
	}
}

func (self *Window) Render() string {
	self.rearange_window()
	var window_with_components strings.Builder
	window_with_components.WriteString(self.render_background())

	if self.Styles.Border != NoBorder {
		window_with_components.WriteString(render_border(self.Position, self.Active, &self.Styles))
	}

	for _, component := range self.Components {
		window_with_components.WriteString(component.Render())
	}

	return window_with_components.String()
}

func (self *Window) rearange_window() {
	// position absolute
	if self.Row != 0 && self.Col != 0 {
		return
	}

	if self.Index == 0 {
		self.Row = 1
		self.Col = 1
		return
	}

	prev_window := self.Oldfart.Windows[self.Index-1]
	if self.Direction == BLOCK {
		self.Row = prev_window.Row + prev_window.Height
		self.Col = prev_window.Col
	} else {
		self.Row = prev_window.Row
		self.Col = prev_window.Col + prev_window.Width
	}
}

func (self *Window) AddComponent(c *Component) {
	c.Index = len(self.Components)
	c.Parent = self
	self.Components = append(self.Components, c)
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

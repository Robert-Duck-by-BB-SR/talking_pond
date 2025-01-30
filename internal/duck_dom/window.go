package duckdom

import (
	"fmt"
	"strings"
)

type Window struct {
	Styles
	Components []*Component
	Position
	Oldfart           *Screen
	Index             int
	ActiveComponentId int
	Active            bool
	OnRender          func()
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

func (self *Window) Render(builder *strings.Builder) {
	self.rearange_window()
	self.render_background(builder)

	if self.Styles.Border != NoBorder {
		render_border(builder, self.Position, self.Active, &self.Styles)
	}

	if self.OnRender != nil {
		self.OnRender()
	}

	for _, component := range self.Components {
		component.Render(builder)
	}

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
	c.Styles.BorderBackground = self.Background
	self.Components = append(self.Components, c)
}

func (self *Window) render_background(bg_builder *strings.Builder) {
	bg_builder.WriteString(self.Styles.Background)
	fillament := strings.Repeat(" ", self.Styles.Width)

	for i := 0; i < self.Styles.Height; i += 1 {
		bg_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Row+i, self.Position.Col))
		bg_builder.WriteString(fillament)
	}
	bg_builder.WriteString(RESET_STYLES)
}

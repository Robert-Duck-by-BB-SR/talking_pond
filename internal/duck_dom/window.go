package duckdom

import (
	"fmt"
	"strings"
)

type Window struct {
	Position
	Parent            *Screen
	Styles
	ActiveComponentId int
	Components        []Component
}

func (self *Window) Render() string {
	if self.Styles.Border.Color != "" && self.Height < 3 {
		panic("Min height with border should be 3")
	}
	window_with_components := self.render_background()

	// TODO: use better way if border is assigned
	if self.Styles.Border.Color != "" {
		window_with_components += render_border(self.Position, &self.Styles)
	}

	for _, component := range self.Components {
		window_with_components += component.Render()
	}

	return window_with_components
}

func (self *Window) render_background() string {
	var bg_builder strings.Builder
	bg_builder.WriteString(self.Styles.Background)
	bg_builder.WriteString(self.Styles.Color)
	fillament := strings.Repeat(" ", self.Styles.Width)

	for i := 0; i < self.Styles.Height; i += 1 {
		bg_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.StartingRow+uint(i), self.Position.StartingCol))
		bg_builder.WriteString(fillament)
	}
	bg_builder.WriteString(RESET_STYLES)

	return bg_builder.String()
}

func assert_dimentions(w *Window) {
	if w.Styles.Width <= 0 || w.Styles.Height <= 0 {
		panic("Width and height should be bigger than 0")
	}

	if w.Styles.Border.Color != "" && w.Height < 3 {
		panic("Min height with border should be 3")
	}
}

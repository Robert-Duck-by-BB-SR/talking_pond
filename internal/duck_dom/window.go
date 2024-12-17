package duckdom

import (
	"fmt"
	"strings"
)

type Window struct {
	Position
	Styles
	ActiveComponentId int
	Components        []Component
}

func (self *Window) Render() string {
	win := self.render_background()

	if self.Styles.Border.Color != "" {
		win += render_border(self.Position, &self.Styles)
	}

	return win
}

func (self *Window) render_background() string {
	var bg_builder strings.Builder
	bg_builder.WriteString(self.Styles.Background)
	bg_builder.WriteString(self.Styles.Color)
	fillament := strings.Repeat(" ", self.Styles.Width)

	for row := self.Position.StartingRow; row < uint(self.Styles.Height); row += 1 {
		bg_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, self.Position.StartingCol))
		bg_builder.WriteString(fillament)
	}
	bg_builder.WriteString(RESET_STYLES)

	return bg_builder.String()
}

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
	var bor_builder strings.Builder
	bor_builder.WriteString(self.Styles.Background)
	bor_builder.WriteString(self.Styles.Color)
	fillament := strings.Repeat(" ", self.Styles.Width)

	for row := self.Position.StartingRow; row < uint(self.Styles.Height); row += 1 {
		bor_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, self.Position.StartingCol))
		bor_builder.WriteString(fillament)
	}
	bor_builder.WriteString(RESET_STYLES)

	return bor_builder.String()
}

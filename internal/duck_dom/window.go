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
	var bor_builder strings.Builder

	if self.Styles.Border.Width != 0 {
		return render_with_border(self, &bor_builder)
	}

	return render_default(self, &bor_builder)
}

func render_default(self *Window, bor_builder *strings.Builder) string {
	bor_builder.WriteString(self.Styles.Background)
	fillament := strings.Repeat(" ", self.Styles.Width)

	for row := self.Position.StartingRow; row < uint(self.Styles.Height); row += 1 {
		bor_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, self.Position.StartingCol))
		bor_builder.WriteString(fillament)
	}
	bor_builder.WriteString(RESET_STYLES)

	return bor_builder.String()
}

func render_with_border(self *Window, main_border_build *strings.Builder) string {
	main_border_build.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Position.StartingRow, self.Position.StartingCol))

	// top border
	var top_border_builder strings.Builder
	top_border_builder.WriteString(self.Styles.Border.Color)
	top_border_builder.WriteString(normalBorder.TopLeft)
	top_border_builder.WriteString(strings.Repeat(normalBorder.Top, self.Styles.Width-2))
	top_border_builder.WriteString(normalBorder.TopRight)
	top_border_builder.WriteString(RESET_STYLES)

	// general filament + border
	var vertical_border_filament strings.Builder
	vertical_border_filament.WriteString(self.Styles.Border.Color)
	vertical_border_filament.WriteString(normalBorder.Left)
	vertical_border_filament.WriteString(RESET_STYLES)
	vertical_border_filament.WriteString(self.Styles.Background)
	vertical_border_filament.WriteString(strings.Repeat(" ", self.Styles.Width-2))
	vertical_border_filament.WriteString(RESET_STYLES)
	vertical_border_filament.WriteString(self.Styles.Border.Color)
	vertical_border_filament.WriteString(normalBorder.Right)
	vertical_border_filament.WriteString(RESET_STYLES)

	// bottom border
	var bottom_border_build strings.Builder
	bottom_border_build.WriteString(self.Styles.Border.Color)
	bottom_border_build.WriteString(normalBorder.BottomLeft)
	bottom_border_build.WriteString(strings.Repeat(normalBorder.Bottom, self.Styles.Width-2))
	bottom_border_build.WriteString(self.Styles.Border.Color)
	bottom_border_build.WriteString(normalBorder.BottomRight)
	bottom_border_build.WriteString(RESET_STYLES)

	// combining my ass
	main_border_build.WriteString(top_border_builder.String())

	for row := self.Position.StartingRow + 1; row < uint(self.Styles.Height)-1; row += 1 {
		main_border_build.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, self.Position.StartingCol))
		main_border_build.WriteString(vertical_border_filament.String())
	}

	main_border_build.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, self.Styles.Height-1, self.Position.StartingCol))
	main_border_build.WriteString(bottom_border_build.String())

	return main_border_build.String()
}

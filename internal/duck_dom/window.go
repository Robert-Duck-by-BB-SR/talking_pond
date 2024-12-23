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

func (self *Window) Render() string {
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

func (self *Window) AddComponent(c *Component) {
	c.Parent = self

	if c.Styles.Display == Block {

	}
	// ASSERT DIMENTIONS BUT I WILL DO IT LATER

	if len(self.Components) == 0 {
		if self.Border.Color == "" {
			c.Position = Position{StartingRow: 1, StartingCol: 1}
		} else {
			c.Position = Position{StartingRow: 2, StartingCol: 2}
		}
	} else {
		last_component := self.Components[len(self.Components)-1]
		lastRow := last_component.StartingRow + last_component.Height - 1
		// can be 1 or 2
		lastCol := last_component.StartingCol

		// TODO: COMBINE WITH BORDER
		// if(newRow + some shit I havent figured out yet)

		newRow := lastRow + c.Height
		if newRow > self.Height {
			panic("Component height will not fit, do a math, you dumbass")
		}

		if lastCol+c.Width > self.Width {
			panic("Component width will not fit, do a math, you dumbass")
		}

		c.Position = Position{StartingRow: newRow, StartingCol: lastCol}
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

func assert_dimentions(w *Window) {
	if w.Styles.Width <= 0 || w.Styles.Height <= 0 {
		panic("Width and height should be bigger than 0")
	}

	if w.Styles.Border.Color != "" && w.Height < 3 {
		panic("Min height with border should be 3")
	}
}

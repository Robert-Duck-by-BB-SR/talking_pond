package duckdom

import (
	"fmt"
	"strings"
)

type Border struct {
	Style BorderStyle
	Color string
}

type BorderStyle struct {
	Top          string
	Bottom       string
	Left         string
	Right        string
	TopLeft      string
	TopRight     string
	BottomLeft   string
	BottomRight  string
	MiddleLeft   string
	MiddleRight  string
	Middle       string
	MiddleTop    string
	MiddleBottom string
}

var (
	NormalBorder = BorderStyle{
		Top:          "─",
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "┌",
		TopRight:     "┐",
		BottomLeft:   "└",
		BottomRight:  "┘",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}

	BoldBorder = BorderStyle{
		Top:          "━",
		Bottom:       "━",
		Left:         "┃",
		Right:        "┃",
		TopLeft:      "┏",
		TopRight:     "┓",
		BottomLeft:   "┗",
		BottomRight:  "┛",
		MiddleLeft:   "┣",
		MiddleRight:  "┫",
		Middle:       "╋",
		MiddleTop:    "┳",
		MiddleBottom: "┻",
	}

	RoundedBorder = BorderStyle{
		Top:          "─",
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "╭",
		TopRight:     "╮",
		BottomLeft:   "╰",
		BottomRight:  "╯",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}
)

func render_border(position Position, styles *Styles) string {
	// box-sizing: border-box;
	var border_builder strings.Builder

	middle := strings.Repeat(styles.Border.Style.Bottom, styles.Width-2)
	top := styles.Border.Style.TopLeft + middle + styles.Border.Style.TopRight
	bottom := styles.Border.Style.BottomLeft + middle + styles.Border.Style.BottomRight

	border_builder.WriteString(RESET_STYLES)
	border_builder.WriteString(styles.Border.Color)
	border_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, position.StartingRow, position.StartingCol))
	border_builder.WriteString(top)

	for i := 1; uint(i) < uint(styles.Height); i += 1 {
		left_wall := fmt.Sprintf(MOVE_CURSOR_TO_POSITION, position.StartingRow + uint(i), position.StartingCol)
		right_wall := fmt.Sprintf(MOVE_CURSOR_TO_POSITION, position.StartingRow + uint(i), position.StartingCol+uint(styles.Width)-1)
		wall := left_wall + styles.Border.Style.Left + right_wall + styles.Border.Style.Right
		border_builder.WriteString(wall)
	}
	border_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, styles.Height + int(position.StartingRow), position.StartingCol))
	border_builder.WriteString(bottom)
	border_builder.WriteString(RESET_STYLES)

	return border_builder.String()
}

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
	var bor_builder strings.Builder

	middle := strings.Repeat(styles.Border.Style.Top, styles.Width-2)
	top := styles.Border.Style.TopLeft + middle + styles.Border.Style.TopRight
	bottom := styles.Border.Style.BottomLeft + middle + styles.Border.Style.BottomRight

	bor_builder.WriteString(RESET_STYLES)
	bor_builder.WriteString(styles.Border.Color)
	bor_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, position.StartingRow, position.StartingCol))
	bor_builder.WriteString(top)
	for row := position.StartingRow + 1; row < uint(styles.Height); row += 1 {
		left_wall := fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, position.StartingCol)
		right_wall := fmt.Sprintf(MOVE_CURSOR_TO_POSITION, row, position.StartingCol+uint(styles.Width)-1)
		wall := left_wall + styles.Border.Style.Left + right_wall + styles.Border.Style.Right
		bor_builder.WriteString(wall)
	}
	bor_builder.WriteString(fmt.Sprintf(MOVE_CURSOR_TO_POSITION, styles.Height, position.StartingCol))
	bor_builder.WriteString(bottom)
	bor_builder.WriteString(RESET_STYLES)

	return bor_builder.String()
}
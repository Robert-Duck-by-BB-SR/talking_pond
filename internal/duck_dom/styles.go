package duckdom

import "fmt"

// oh boy CSS is comming

type Styles struct {
	Width      int
	Height     int
	Paddding   int
	Maaargin   int
	Background string
	Border
}

type BorderStyle int

const (
	Solid BorderStyle = iota
	Bold
)

type Border struct {
	Width int
	Style BorderStyle
	Color string
}

type BorderParts struct {
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
	normalBorder = BorderParts{
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

	boldBorder = BorderParts{
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
)

func MakeRGBBackground(r, g, b int) string {
	return BG_KEY + fmt.Sprintf("%d;%d;%dm", r, g, b)
}

func MakeRGBTextColor(r, g, b int) string {
	return FG_KEY + fmt.Sprintf("%d;%d;%dm", r, g, b)
}

func (self *Styles) SetWidth(w int) *Styles {
	self.Width = w
	return self
}

func (self *Styles) SetHeight(h int) *Styles {
	self.Height = h
	return self
}

func (self *Styles) SetPadding(p int) *Styles {
	self.Paddding = p
	return self
}

func (self *Styles) SetMargin(m int) *Styles {
	self.Maaargin = m
	return self
}

func (self *Styles) SetBackground(b string) *Styles {
	self.Background = b
	return self
}

func (self *Styles) SetBorder(b Border) *Styles {
	self.Border = b
	return self
}

func (self *Styles) Compiled() string {
	return ""
}

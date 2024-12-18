package duckdom

import (
	"fmt"
	"strings"
)

// oh boy CSS is comming

type Styles struct {
	Width      int
	Height     int
	Paddding   int
	Maaargin   int
	Background string
	TextColor  string
	Border
}

func MakeRGBBackground(r, g, b int) string {
	return BG_KEY + fmt.Sprintf(RGB, r, g, b)
}

func MakeRGBTextColor(r, g, b int) string {
	return FG_KEY + fmt.Sprintf(RGB, r, g, b)
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

func (self *Styles) SetTextColor(tc string) *Styles {
	self.TextColor = tc
	return self
}

func (self *Styles) SetBorder(b Border) *Styles {
	self.Border = b
	return self
}

func (self *Styles) Compile() string {
	var styles_builder strings.Builder
	styles_builder.WriteString(self.Background)
	styles_builder.WriteString(self.TextColor)
	return styles_builder.String()
}

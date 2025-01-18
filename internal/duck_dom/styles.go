package duckdom

import (
	"fmt"
	"strings"
)

// oh boy CSS is comming
type Direction int

const (
	BLOCK Direction = iota
	INLINE
)

type Styles struct {
	Border
	BorderBackground string
	Background       string
	TextColor        string
	Width            int
	MinWidth         int
	MaxWidth         int
	Height           int
	MaxHeight        int
	Paddding         int
	Direction
}

func MakeRGBBackground(r, g, b int) string {
	return BG_KEY + fmt.Sprintf(RGB, r, g, b)
}

func MakeRGBTextColor(r, g, b int) string {
	return FG_KEY + fmt.Sprintf(RGB, r, g, b)
}

func (self *Styles) Compile(styles_builder *strings.Builder) {
	styles_builder.WriteString(self.Background)
	styles_builder.WriteString(self.TextColor)
}

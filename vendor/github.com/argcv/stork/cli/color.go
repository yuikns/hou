package cli

import (
	"bytes"
	"fmt"
)

type Color int
type Attr int

const (
	Inherit Color = iota
	NoColor
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite
	Bold Attr = iota
	Underscore
	Blink
	ReverseVideo
	Concealed
)

// ref: https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
var (
	prefix = "\033["
	fgCode = map[Color]int{
		NoColor:       0,
		Black:         30,
		Red:           31,
		Green:         32,
		Yellow:        33,
		Blue:          34,
		Magenta:       35,
		Cyan:          36,
		White:         37,
		BrightBlack:   90,
		BrightRed:     91,
		BrightGreen:   92,
		BrightYellow:  93,
		BrightBlue:    94,
		BrightMagenta: 95,
		BrightCyan:    96,
		BrightWhite:   97,
	}
	bgCode = map[Color]int{
		NoColor:       0,
		Black:         40,
		Red:           41,
		Green:         42,
		Yellow:        43,
		Blue:          44,
		Magenta:       45,
		Cyan:          46,
		White:         47,
		BrightBlack:   100,
		BrightRed:     101,
		BrightGreen:   102,
		BrightYellow:  103,
		BrightBlue:    104,
		BrightMagenta: 105,
		BrightCyan:    106,
		BrightWhite:   107,
	}
)

type ColoredText struct {
	Text  string
	FG    Color
	BG    Color
	Attrs map[int]bool
	Clear bool
}

func NewColoredText(text string) *ColoredText {
	ct := ColoredText{
		Text:  text,
		Attrs: map[int]bool{},
	}
	return ct.Init()
}

func (c *ColoredText) Reset() *ColoredText {
	c.FG = Inherit
	c.BG = Inherit
	c.Clear = true
	return c
}

func (c *ColoredText) Init() *ColoredText {
	c.FG = Inherit
	c.BG = Inherit
	c.Clear = false
	c.Attrs = map[int]bool{}
	return c
}

func (c *ColoredText) SetFG(color Color) *ColoredText {
	c.FG = color
	return c
}

func (c *ColoredText) SetBG(color Color) *ColoredText {
	c.BG = color
	return c
}

func (c *ColoredText) SetAttr(attr Attr, on bool) *ColoredText {
	op := func(val int) {
		if on {
			c.Attrs[val] = true
		} else {
			delete(c.Attrs, val)
		}
	}
	switch attr {
	case Bold:
		op(1) // bold on/off
	case Underscore:
		op(4) // Underscore on/off
	case Blink:
		op(5) // Blink on/off
	case ReverseVideo:
		op(7) // Reverse Video on/off
	case Concealed:
		op(8) // Concealed on/off
	}
	return c
}

func (c *ColoredText) String() string {
	var buff []string
	if c.Clear {
		// Reset
		buff = append(buff, prefix)
		buff = append(buff, "0m")
		buff = append(buff, c.Text)
	} else if c.FG == Inherit && c.BG == Inherit && len(c.Attrs) == 0 {
		// Just Pass
		buff = append(buff, c.Text)
	} else {
		buff = append(buff, prefix)
		cnt := 0
		checkPrefix := func() {
			if cnt == 0 {
				cnt++
			} else {
				cnt++
				buff = append(buff, ";")
			}
		}
		if c.FG != Inherit {
			buff = append(buff, fmt.Sprint(fgCode[c.FG]))
			cnt++
		}

		if c.BG != Inherit {
			checkPrefix()
			buff = append(buff, fmt.Sprint(bgCode[c.BG]))
		}

		for key, _ := range c.Attrs {
			checkPrefix()
			buff = append(buff, fmt.Sprint(key))
		}

		buff = append(buff, "m")

		buff = append(buff, c.Text)
		// clear all
		buff = append(buff, prefix)
		buff = append(buff, "0m")

	}
	out := bytes.NewBufferString("")
	for _, buf := range buff {
		//fmt.Printf("append: %T%v\n\n", buf, buf)
		out.WriteString(buf)
	}
	return out.String()
}

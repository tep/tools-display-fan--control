package main

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"toolman.org/base/toolman/v2"
)

type configuration struct {
	fontID int

	ohBase string

	onString  string
	offString string
	chgString string

	width  int
	spacer string

	onColor  string
	offColor string
	chgColor string
	bgColor  string
}

func defaults() *configuration {
	return &configuration{
		fontID:    9,
		ohBase:    defaultOpenhabBase,
		onString:  "",
		offString: "ﴛ",
		chgString: "ﴛ",
		width:     3,
		spacer:    "",
		onColor:   "#eee0e0",
		offColor:  "#555050",
		chgColor:  "#888080",
		bgColor:   "#0a0703",
	}
}

func (c *configuration) flags() *toolman.InitOption {
	fs := pflag.NewFlagSet(program, pflag.ExitOnError)

	fs.IntVar(&c.fontID, "font-id", c.fontID, "Polybar Font ID")
	fs.IntVar(&c.width, "width", c.width, "Total character width")

	fs.StringVar(&c.ohBase, "openhab-base-url", c.ohBase, "Openhab base URL")

	fs.StringVar(&c.onString, "on-string", c.onString, "String to display when fan is ON")
	fs.StringVar(&c.offString, "off-string", c.offString, "String to display when fan is OFF")
	fs.StringVar(&c.chgString, "change-string", c.chgString, "String to display when fan is changing state")
	fs.StringVar(&c.spacer, "spacer", c.spacer, "String to use as a spacer (uses BG color as FG)")
	fs.StringVar(&c.onColor, "on-color", c.onColor, "Hex `color` value when fan is ON")
	fs.StringVar(&c.offColor, "off-color", c.offColor, "Hex `color` value when fan is OFF")
	fs.StringVar(&c.chgColor, "change-color", c.chgColor, "Hex `color` value when fan is changing state")
	fs.StringVar(&c.bgColor, "bg-color", c.bgColor, "Hex background `color` value")

	return toolman.FlagSet(fs)
}

func (c *configuration) stateOutputMap() map[fanState]string {
	output := func(clr, str string) string {
		sw := (c.width - len([]rune(str))) / len([]rune(c.spacer))
		if sw < 0 {
			sw = 0
		}

		return fmt.Sprintf("%%{T%d}%%{B%s}%%{F%s}%s%%{F%s}%s%%{F-}%%{B-}%%{T-}",
			c.fontID, c.bgColor, clr, str, c.bgColor, strings.Repeat(c.spacer, sw))
	}

	return map[fanState]string{
		fsOff:    output(c.offColor, c.offString),
		fsOn:     output(c.onColor, c.onString),
		fsChange: output(c.chgColor, c.chgString),
	}
}

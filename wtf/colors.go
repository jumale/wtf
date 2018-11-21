package wtf

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
)

// Color is a color label (e.g. "black", "red"m etc),
// which is also convertable from/to tcell constants
type Color string

// ToTcell converts color label to a corresponding
// tcell color constant
func (c Color) ToTcell() tcell.Color {
	return tcell.GetColor(string(c))
}

// FromTcell creates a new color label from
// corresponding tcell constant
func (Color) FromTcell(c tcell.Color) Color {
	for label, color := range tcell.ColorNames {
		if color == c {
			return Color(label)
		}
	}
	return Color("default") // in case if the constant is not in the color map
}

// Converts all ASCII colors in the text to corresponding "tview" tags.
func ASCIItoTviewColors(text string) string {
	boldRegExp := regexp.MustCompile(`\033\[1m`)
	fgColorRegExp := regexp.MustCompile(`\033\[38;5;(?P<color>\d+);*\d*m`)
	resColorRegExp := regexp.MustCompile(`\033\[0m`)

	return resColorRegExp.ReplaceAllString(
		boldRegExp.ReplaceAllString(
			fgColorRegExp.ReplaceAllStringFunc(
				text, replaceWithHexColorString), `[::b]`), `[-]`)
}

/* -------------------- Unexported Functions -------------------- */

func replaceWithHexColorString(substring string) string {
	colorID, err := strconv.Atoi(
		strings.Trim(
			strings.Split(substring, ";")[2],
			"m",
		),
	)
	if err != nil {
		return substring
	}
	tcellColor := tcell.Color(colorID)
	hexColor := strconv.FormatInt(int64(tcellColor), 16)

	return "[#" + hexColor + "]"
}

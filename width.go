package main

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

func stringWidth(s string) int {
	return runewidth.StringWidth(s)
}

func fillRight(s string, width int) string {
	w := stringWidth(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

func expandTabs(s string, tabWidth int) string {
	var buf strings.Builder
	col := 0
	for _, r := range s {
		if r == '\t' {
			spaces := tabWidth - (col % tabWidth)
			buf.WriteString(strings.Repeat(" ", spaces))
			col += spaces
		} else {
			buf.WriteRune(r)
			col += runewidth.RuneWidth(r)
		}
	}
	return buf.String()
}

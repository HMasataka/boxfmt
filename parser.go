package main

import (
	"strings"
)

type lineType int

const (
	linePlain lineType = iota
	lineTopBorder
	lineBottomBorder
	lineDivider
	lineContent
)

type classifiedLine struct {
	raw      string
	typ      lineType
	indent   string
	trimmed  string
	isASCII  bool
}

type boxRegion struct {
	startIdx int
	endIdx   int
	lines    []classifiedLine
	indent   string
}

func isUnicodeHorizontal(r rune) bool {
	return r == '─'
}

func isASCIIHorizontal(r rune) bool {
	return r == '-'
}

func isHorizontal(r rune) bool {
	return isUnicodeHorizontal(r) || isASCIIHorizontal(r)
}

func isUnicodeVertical(r rune) bool {
	return r == '│'
}

func isASCIIVertical(r rune) bool {
	return r == '|'
}

func isVertical(r rune) bool {
	return isUnicodeVertical(r) || isASCIIVertical(r)
}

func isBoxDrawing(r rune) bool {
	switch r {
	case '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼', '─', '│':
		return true
	case '+', '-', '|':
		return true
	}
	return false
}

func firstNonSpace(s string) (rune, int) {
	for i, r := range s {
		if r != ' ' && r != '\t' {
			return r, i
		}
	}
	return 0, -1
}

func lastNonSpace(s string) rune {
	runes := []rune(strings.TrimRight(s, " \t"))
	if len(runes) == 0 {
		return 0
	}
	return runes[len(runes)-1]
}

func getIndent(s string) string {
	for i, r := range s {
		if r != ' ' && r != '\t' {
			return s[:i]
		}
	}
	return s
}

func isBorderLine(trimmed string, leftCorner, rightCorner, midJunction rune, horizontal rune) bool {
	runes := []rune(trimmed)
	if len(runes) < 2 {
		return false
	}
	if runes[0] != leftCorner {
		return false
	}
	if runes[len(runes)-1] != rightCorner {
		return false
	}
	for _, r := range runes[1 : len(runes)-1] {
		if r != horizontal && r != midJunction {
			return false
		}
	}
	return true
}

func classifyLine(line string) classifiedLine {
	indent := getIndent(line)
	trimmed := strings.TrimSpace(line)

	if trimmed == "" {
		return classifiedLine{raw: line, typ: linePlain, indent: indent, trimmed: trimmed}
	}

	firstR, _ := firstNonSpace(line)
	lastR := lastNonSpace(line)

	// Unicode TopBorder: ┌...┐
	if isBorderLine(trimmed, '┌', '┐', '┬', '─') {
		return classifiedLine{raw: line, typ: lineTopBorder, indent: indent, trimmed: trimmed, isASCII: false}
	}

	// Unicode BottomBorder: └...┘
	if isBorderLine(trimmed, '└', '┘', '┴', '─') {
		return classifiedLine{raw: line, typ: lineBottomBorder, indent: indent, trimmed: trimmed, isASCII: false}
	}

	// Unicode Divider: ├...┤
	if isBorderLine(trimmed, '├', '┤', '┼', '─') {
		return classifiedLine{raw: line, typ: lineDivider, indent: indent, trimmed: trimmed, isASCII: false}
	}

	// ASCII border lines: +---+---+
	if isASCIIBorderLine(trimmed) {
		typ := classifyASCIIBorder(trimmed)
		return classifiedLine{raw: line, typ: typ, indent: indent, trimmed: trimmed, isASCII: true}
	}

	// Content line: │...│ or |...|
	if (isVertical(firstR)) && (isVertical(lastR)) {
		ascii := isASCIIVertical(firstR)
		return classifiedLine{raw: line, typ: lineContent, indent: indent, trimmed: trimmed, isASCII: ascii}
	}

	return classifiedLine{raw: line, typ: linePlain, indent: indent, trimmed: trimmed}
}

func isASCIIBorderLine(trimmed string) bool {
	runes := []rune(trimmed)
	if len(runes) < 2 {
		return false
	}
	if runes[0] != '+' || runes[len(runes)-1] != '+' {
		return false
	}
	for _, r := range runes[1 : len(runes)-1] {
		if r != '-' && r != '+' {
			return false
		}
	}
	return true
}

func classifyASCIIBorder(trimmed string) lineType {
	// ASCII borders are ambiguous. We'll need context to determine
	// if they're top, bottom, or divider. Default to TopBorder;
	// the grouping logic will reclassify.
	return lineTopBorder
}

func classifyLines(lines []string) []classifiedLine {
	classified := make([]classifiedLine, len(lines))
	for i, line := range lines {
		classified[i] = classifyLine(line)
	}
	return classified
}

func detectBoxRegions(classified []classifiedLine) []boxRegion {
	var regions []boxRegion
	n := len(classified)
	i := 0

	for i < n {
		if classified[i].typ == linePlain {
			i++
			continue
		}

		// Found a non-plain line, collect consecutive non-plain lines
		start := i
		for i < n && classified[i].typ != linePlain {
			i++
		}
		end := i // exclusive

		group := classified[start:end]

		// Reclassify ASCII borders based on position
		reclassifyASCIIBorders(group)

		// Check if this forms a valid box
		if isValidBox(group) {
			indent := commonIndent(group)
			region := boxRegion{
				startIdx: start,
				endIdx:   end,
				lines:    make([]classifiedLine, len(group)),
				indent:   indent,
			}
			copy(region.lines, group)
			regions = append(regions, region)
		}
	}

	return regions
}

func reclassifyASCIIBorders(group []classifiedLine) {
	if len(group) < 2 {
		return
	}

	for i := range group {
		if !group[i].isASCII {
			continue
		}
		if group[i].typ != lineTopBorder && group[i].typ != lineBottomBorder && group[i].typ != lineDivider {
			continue
		}
		// Only reclassify ASCII border lines (those matched by isASCIIBorderLine)
		if i == 0 {
			group[i].typ = lineTopBorder
		} else if i == len(group)-1 {
			group[i].typ = lineBottomBorder
		} else {
			group[i].typ = lineDivider
		}
	}
}

func isValidBox(group []classifiedLine) bool {
	if len(group) < 2 {
		return false
	}

	if group[0].typ != lineTopBorder {
		return false
	}
	if group[len(group)-1].typ != lineBottomBorder {
		return false
	}

	for _, cl := range group[1 : len(group)-1] {
		if cl.typ != lineContent && cl.typ != lineDivider {
			return false
		}
	}

	return true
}

func commonIndent(group []classifiedLine) string {
	if len(group) == 0 {
		return ""
	}
	indent := group[0].indent
	for _, cl := range group[1:] {
		indent = shorterIndent(indent, cl.indent)
	}
	return indent
}

func shorterIndent(a, b string) string {
	if len(a) <= len(b) {
		return a
	}
	return b
}

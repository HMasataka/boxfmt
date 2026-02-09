package main

import (
	"strings"
)

func fixBoxRegion(region boxRegion) []string {
	columns := detectColumns(region)
	if len(columns) == 0 {
		return fixSingleColumnBox(region)
	}
	return fixMultiColumnBox(region, columns)
}

// detectColumns returns column separator positions from the first border line.
// Returns nil or single-element slice for single-column boxes.
func detectColumns(region boxRegion) []int {
	if len(region.lines) == 0 {
		return nil
	}

	border := region.lines[0].trimmed
	runes := []rune(border)

	var separators []int
	for i, r := range runes {
		if i == 0 || i == len(runes)-1 {
			continue
		}
		if r == '┬' || r == '┴' || r == '┼' {
			separators = append(separators, i)
		}
		if r == '+' && region.lines[0].isASCII {
			separators = append(separators, i)
		}
	}

	return separators
}

func fixSingleColumnBox(region boxRegion) []string {
	// Extract content texts and compute max width
	var contentTexts []string
	var contentIndices []int

	for i, cl := range region.lines {
		if cl.typ == lineContent {
			text := extractContentText(cl.trimmed)
			contentTexts = append(contentTexts, text)
			contentIndices = append(contentIndices, i)
		}
	}

	maxWidth := 0
	for _, text := range contentTexts {
		w := stringWidth(text)
		if w > maxWidth {
			maxWidth = w
		}
	}

	// Rebuild lines
	result := make([]string, len(region.lines))
	for i, cl := range region.lines {
		switch cl.typ {
		case lineTopBorder:
			result[i] = region.indent + buildBorderLine(cl, maxWidth)
		case lineBottomBorder:
			result[i] = region.indent + buildBorderLine(cl, maxWidth)
		case lineDivider:
			result[i] = region.indent + buildBorderLine(cl, maxWidth)
		case lineContent:
			text := extractContentText(cl.trimmed)
			leftV, rightV := getVerticalChars(cl)
			padded := fillRight(text, maxWidth)
			result[i] = region.indent + string(leftV) + " " + padded + " " + string(rightV)
		}
	}

	return result
}

func fixMultiColumnBox(region boxRegion, separators []int) []string {
	// Parse columns from content lines
	numCols := len(separators) + 1
	colContents := make([][]string, numCols)
	for i := range colContents {
		colContents[i] = []string{}
	}

	var contentIndices []int
	var contentCols [][]string

	for i, cl := range region.lines {
		if cl.typ == lineContent {
			cols := splitContentColumns(cl.trimmed, numCols)
			contentCols = append(contentCols, cols)
			contentIndices = append(contentIndices, i)
			for c := 0; c < numCols; c++ {
				if c < len(cols) {
					colContents[c] = append(colContents[c], cols[c])
				}
			}
		}
	}

	// Compute max width for each column
	maxWidths := make([]int, numCols)
	for c := 0; c < numCols; c++ {
		for _, text := range colContents[c] {
			w := stringWidth(text)
			if w > maxWidths[c] {
				maxWidths[c] = w
			}
		}
	}

	// Rebuild lines
	result := make([]string, len(region.lines))
	contentIdx := 0

	for i, cl := range region.lines {
		switch cl.typ {
		case lineTopBorder:
			result[i] = region.indent + buildMultiColBorderLine(cl, maxWidths, getTopBorderChars(cl))
		case lineBottomBorder:
			result[i] = region.indent + buildMultiColBorderLine(cl, maxWidths, getBottomBorderChars(cl))
		case lineDivider:
			result[i] = region.indent + buildMultiColBorderLine(cl, maxWidths, getDividerChars(cl))
		case lineContent:
			cols := contentCols[contentIdx]
			contentIdx++
			leftV, rightV := getVerticalChars(cl)
			var buf strings.Builder
			buf.WriteRune(leftV)
			for c := 0; c < numCols; c++ {
				text := ""
				if c < len(cols) {
					text = cols[c]
				}
				padded := fillRight(text, maxWidths[c])
				buf.WriteString(" " + padded + " ")
				if c < numCols-1 {
					// Use inner vertical separator
					buf.WriteRune(getInnerVertical(cl))
				}
			}
			buf.WriteRune(rightV)
			result[i] = region.indent + buf.String()
		}
	}

	return result
}

func extractContentText(trimmed string) string {
	runes := []rune(trimmed)
	if len(runes) < 2 {
		return ""
	}
	// Remove left vertical and right vertical
	inner := string(runes[1 : len(runes)-1])

	// Remove exactly one leading space and one trailing space if present
	if len(inner) > 0 && inner[0] == ' ' {
		inner = inner[1:]
	}
	if len(inner) > 0 && inner[len(inner)-1] == ' ' {
		inner = inner[:len(inner)-1]
	}

	return inner
}

func splitContentColumns(trimmed string, numCols int) []string {
	runes := []rune(trimmed)
	if len(runes) < 2 {
		return nil
	}

	// Remove outer verticals
	inner := runes[1 : len(runes)-1]

	// Split by inner vertical characters (│ or |)
	var cols []string
	var current []rune

	for _, r := range inner {
		if isVertical(r) {
			text := strings.TrimSpace(string(current))
			cols = append(cols, text)
			current = nil
		} else {
			current = append(current, r)
		}
	}
	// Last column
	text := strings.TrimSpace(string(current))
	cols = append(cols, text)

	return cols
}

func getVerticalChars(cl classifiedLine) (rune, rune) {
	runes := []rune(cl.trimmed)
	if len(runes) < 2 {
		if cl.isASCII {
			return '|', '|'
		}
		return '│', '│'
	}
	return runes[0], runes[len(runes)-1]
}

func getInnerVertical(cl classifiedLine) rune {
	if cl.isASCII {
		return '|'
	}
	return '│'
}

type borderChars struct {
	left       rune
	right      rune
	horizontal rune
	junction   rune
}

func getTopBorderChars(cl classifiedLine) borderChars {
	if cl.isASCII {
		return borderChars{'+', '+', '-', '+'}
	}
	return borderChars{'┌', '┐', '─', '┬'}
}

func getBottomBorderChars(cl classifiedLine) borderChars {
	if cl.isASCII {
		return borderChars{'+', '+', '-', '+'}
	}
	return borderChars{'└', '┘', '─', '┴'}
}

func getDividerChars(cl classifiedLine) borderChars {
	if cl.isASCII {
		return borderChars{'+', '+', '-', '+'}
	}
	return borderChars{'├', '┤', '─', '┼'}
}

func buildBorderLine(cl classifiedLine, contentWidth int) string {
	runes := []rune(cl.trimmed)
	if len(runes) < 2 {
		return cl.trimmed
	}

	left := runes[0]
	right := runes[len(runes)-1]
	horiz := getHorizontalChar(cl)

	return string(left) + strings.Repeat(string(horiz), contentWidth+2) + string(right)
}

func buildMultiColBorderLine(cl classifiedLine, maxWidths []int, chars borderChars) string {
	var buf strings.Builder
	buf.WriteRune(chars.left)
	for c, w := range maxWidths {
		buf.WriteString(strings.Repeat(string(chars.horizontal), w+2))
		if c < len(maxWidths)-1 {
			buf.WriteRune(chars.junction)
		}
	}
	buf.WriteRune(chars.right)
	return buf.String()
}

func getHorizontalChar(cl classifiedLine) rune {
	if cl.isASCII {
		return '-'
	}
	return '─'
}

func processFile(content string) string {
	// Preserve trailing newline state
	hasTrailingNewline := len(content) > 0 && content[len(content)-1] == '\n'

	lines := strings.Split(content, "\n")

	// Remove last empty element from Split if file ends with newline
	if hasTrailingNewline && len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	// Expand tabs
	for i, line := range lines {
		lines[i] = expandTabs(line, 4)
	}

	// Classify lines
	classified := classifyLines(lines)

	// Detect box regions
	regions := detectBoxRegions(classified)

	// Apply fixes (process in reverse to preserve indices)
	for i := len(regions) - 1; i >= 0; i-- {
		region := regions[i]
		fixed := fixBoxRegion(region)

		// Replace lines in-place
		newLines := make([]string, 0, len(lines)-region.endIdx+region.startIdx+len(fixed))
		newLines = append(newLines, lines[:region.startIdx]...)
		newLines = append(newLines, fixed...)
		newLines = append(newLines, lines[region.endIdx:]...)
		lines = newLines
	}

	result := strings.Join(lines, "\n")
	if hasTrailingNewline {
		result += "\n"
	}

	return result
}

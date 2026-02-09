package main

import (
	"testing"
)

func TestClassifyLine(t *testing.T) {
	tests := []struct {
		input string
		want  lineType
	}{
		// Unicode borders
		{"┌────────┐", lineTopBorder},
		{"└────────┘", lineBottomBorder},
		{"├────────┤", lineDivider},
		{"│ text   │", lineContent},

		// Unicode multi-column borders
		{"┌────┬────┐", lineTopBorder},
		{"└────┴────┘", lineBottomBorder},
		{"├────┼────┤", lineDivider},

		// ASCII borders
		{"+--------+", lineTopBorder},
		{"+----+----+", lineTopBorder},

		// ASCII content
		{"| text   |", lineContent},

		// Indented lines
		{"  ┌────────┐", lineTopBorder},
		{"  │ text   │", lineContent},

		// Plain lines
		{"hello world", linePlain},
		{"", linePlain},
		{"  some text  ", linePlain},

		// Markdown table (not a box)
		{"| col1 | col2 |", lineContent},
	}
	for _, tt := range tests {
		cl := classifyLine(tt.input)
		if cl.typ != tt.want {
			t.Errorf("classifyLine(%q).typ = %v, want %v", tt.input, cl.typ, tt.want)
		}
	}
}

func TestClassifyLineASCIIFlag(t *testing.T) {
	cl := classifyLine("+--------+")
	if !cl.isASCII {
		t.Error("expected isASCII=true for ASCII border")
	}

	cl = classifyLine("┌────────┐")
	if cl.isASCII {
		t.Error("expected isASCII=false for Unicode border")
	}

	cl = classifyLine("| text |")
	if !cl.isASCII {
		t.Error("expected isASCII=true for ASCII content")
	}

	cl = classifyLine("│ text │")
	if cl.isASCII {
		t.Error("expected isASCII=false for Unicode content")
	}
}

func TestDetectBoxRegions(t *testing.T) {
	lines := []string{
		"some text",
		"┌────────┐",
		"│ hello  │",
		"│ world  │",
		"└────────┘",
		"more text",
	}
	classified := classifyLines(lines)
	regions := detectBoxRegions(classified)

	if len(regions) != 1 {
		t.Fatalf("expected 1 region, got %d", len(regions))
	}

	r := regions[0]
	if r.startIdx != 1 || r.endIdx != 5 {
		t.Errorf("region bounds = [%d, %d), want [1, 5)", r.startIdx, r.endIdx)
	}
}

func TestDetectBoxRegionsMultiple(t *testing.T) {
	lines := []string{
		"┌──┐",
		"│ A│",
		"└──┘",
		"text",
		"┌──┐",
		"│ B│",
		"└──┘",
	}
	classified := classifyLines(lines)
	regions := detectBoxRegions(classified)

	if len(regions) != 2 {
		t.Fatalf("expected 2 regions, got %d", len(regions))
	}
}

func TestDetectBoxRegionsInvalid(t *testing.T) {
	// Missing bottom border
	lines := []string{
		"┌──┐",
		"│ A│",
		"plain text",
	}
	classified := classifyLines(lines)
	regions := detectBoxRegions(classified)

	if len(regions) != 0 {
		t.Fatalf("expected 0 regions, got %d", len(regions))
	}
}

func TestCommonIndent(t *testing.T) {
	group := []classifiedLine{
		{indent: "  "},
		{indent: "    "},
		{indent: "  "},
	}
	got := commonIndent(group)
	if got != "  " {
		t.Errorf("commonIndent = %q, want %q", got, "  ")
	}
}

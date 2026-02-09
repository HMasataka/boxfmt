package main

import (
	"strings"
	"testing"
)

func TestExtractContentText(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"│ hello │", "hello"},
		{"│ 日本語 │", "日本語"},
		{"| text |", "text"},
		{"│  │", ""},
		{"│ hello world │", "hello world"},
	}
	for _, tt := range tests {
		got := extractContentText(tt.input)
		if got != tt.want {
			t.Errorf("extractContentText(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestSplitContentColumns(t *testing.T) {
	tests := []struct {
		input   string
		numCols int
		want    []string
	}{
		{"│ a │ b │", 2, []string{"a", "b"}},
		{"│ hello │ world │", 2, []string{"hello", "world"}},
		{"| a | b | c |", 3, []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		got := splitContentColumns(tt.input, tt.numCols)
		if len(got) != len(tt.want) {
			t.Errorf("splitContentColumns(%q, %d) len = %d, want %d", tt.input, tt.numCols, len(got), len(tt.want))
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("splitContentColumns(%q, %d)[%d] = %q, want %q", tt.input, tt.numCols, i, got[i], tt.want[i])
			}
		}
	}
}

func TestFixSingleColumnUnicode(t *testing.T) {
	input := strings.Join([]string{
		"┌──────┐",
		"│ hi │",
		"│ 日本語 │",
		"└──────┘",
	}, "\n")

	result := processFile(input)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	// maxContentWidth should be 6 (日本語 = 6 display width)
	// border = left + (6+2)*horizontal + right = 1+8+1
	expectedBorder := "┌────────┐"
	if lines[0] != expectedBorder {
		t.Errorf("top border = %q, want %q", lines[0], expectedBorder)
	}

	// Content "hi" should be padded to width 6
	expectedContent := "│ hi     │"
	if lines[1] != expectedContent {
		t.Errorf("content line 1 = %q, want %q", lines[1], expectedContent)
	}

	expectedContent2 := "│ 日本語 │"
	if lines[2] != expectedContent2 {
		t.Errorf("content line 2 = %q, want %q", lines[2], expectedContent2)
	}
}

func TestFixSingleColumnASCII(t *testing.T) {
	input := strings.Join([]string{
		"+------+",
		"| hi |",
		"| 日本語 |",
		"+------+",
	}, "\n")

	result := processFile(input)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	expectedBorder := "+--------+"
	if lines[0] != expectedBorder {
		t.Errorf("top border = %q, want %q", lines[0], expectedBorder)
	}

	expectedContent := "| hi     |"
	if lines[1] != expectedContent {
		t.Errorf("content line 1 = %q, want %q", lines[1], expectedContent)
	}
}

func TestFixMultiColumn(t *testing.T) {
	input := strings.Join([]string{
		"┌──┬──┐",
		"│ A│ B│",
		"│ 日本│ CD│",
		"└──┴──┘",
	}, "\n")

	result := processFile(input)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	// Column 1: max width = 4 (日本), Column 2: max width = 2 (CD or B)
	expectedBorder := "┌──────┬────┐"
	if lines[0] != expectedBorder {
		t.Errorf("top border = %q, want %q", lines[0], expectedBorder)
	}
}

func TestFixPreservesIndent(t *testing.T) {
	input := strings.Join([]string{
		"  ┌──────┐",
		"  │ hi │",
		"  │ 日本語 │",
		"  └──────┘",
	}, "\n")

	result := processFile(input)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	for _, line := range lines {
		if !strings.HasPrefix(line, "  ") {
			t.Errorf("line %q does not preserve indent", line)
		}
	}
}

func TestFixWithDivider(t *testing.T) {
	input := strings.Join([]string{
		"┌──────┐",
		"│ hi │",
		"├──────┤",
		"│ 日本語 │",
		"└──────┘",
	}, "\n")

	result := processFile(input)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")

	expectedDivider := "├────────┤"
	if lines[2] != expectedDivider {
		t.Errorf("divider = %q, want %q", lines[2], expectedDivider)
	}
}

func TestProcessFilePreservesPlainLines(t *testing.T) {
	input := "hello\nworld\n"
	result := processFile(input)
	if result != input {
		t.Errorf("processFile(%q) = %q, want %q", input, result, input)
	}
}

func TestProcessFileEmptyContent(t *testing.T) {
	result := processFile("")
	if result != "" {
		t.Errorf("processFile(\"\") = %q, want \"\"", result)
	}
}

func TestPreserveTrailingNewline(t *testing.T) {
	input := "┌──┐\n│ A│\n└──┘\n"
	result := processFile(input)
	if result[len(result)-1] != '\n' {
		t.Error("trailing newline not preserved")
	}

	input2 := "┌──┐\n│ A│\n└──┘"
	result2 := processFile(input2)
	if result2[len(result2)-1] == '\n' {
		t.Error("no-trailing-newline not preserved")
	}
}

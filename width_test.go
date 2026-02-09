package main

import (
	"testing"
)

func TestStringWidth(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"abc", 3},
		{"日本語", 6},
		{"abc日本語", 9},
		{"─", 1},
		{"│", 1},
		{"┌", 1},
		{"", 0},
		{" ", 1},
		{"hello world", 11},
		{"テスト", 6},
		{"A全B角C", 7},
	}
	for _, tt := range tests {
		got := stringWidth(tt.input)
		if got != tt.want {
			t.Errorf("stringWidth(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestFillRight(t *testing.T) {
	tests := []struct {
		input string
		width int
		want  string
	}{
		{"abc", 5, "abc  "},
		{"日本語", 8, "日本語  "},
		{"abc", 3, "abc"},
		{"abc", 2, "abc"},
		{"テスト", 6, "テスト"},
		{"テスト", 10, "テスト    "},
		{"A全角", 6, "A全角 "},
	}
	for _, tt := range tests {
		got := fillRight(tt.input, tt.width)
		if got != tt.want {
			t.Errorf("fillRight(%q, %d) = %q, want %q", tt.input, tt.width, got, tt.want)
		}
	}
}

func TestExpandTabs(t *testing.T) {
	tests := []struct {
		input    string
		tabWidth int
		want     string
	}{
		{"\thello", 4, "    hello"},
		{"ab\tcd", 4, "ab  cd"},
		{"abcd\tef", 4, "abcd    ef"},
		{"no tabs", 4, "no tabs"},
		{"\t\t", 4, "        "},
	}
	for _, tt := range tests {
		got := expandTabs(tt.input, tt.tabWidth)
		if got != tt.want {
			t.Errorf("expandTabs(%q, %d) = %q, want %q", tt.input, tt.tabWidth, got, tt.want)
		}
	}
}

package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "update golden files")

func TestGoldenFiles(t *testing.T) {
	entries, err := filepath.Glob("testdata/*.input.md")
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) == 0 {
		t.Fatal("no testdata files found")
	}

	for _, inputPath := range entries {
		name := strings.TrimSuffix(filepath.Base(inputPath), ".input.md")
		expectedPath := strings.Replace(inputPath, ".input.md", ".expected.md", 1)

		t.Run(name, func(t *testing.T) {
			inputData, err := os.ReadFile(inputPath)
			if err != nil {
				t.Fatal(err)
			}

			result := processFile(string(inputData))

			if *update {
				if err := os.WriteFile(expectedPath, []byte(result), 0644); err != nil {
					t.Fatal(err)
				}
				t.Logf("updated %s", expectedPath)
				return
			}

			expectedData, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Fatal(err)
			}

			expected := string(expectedData)
			if result != expected {
				t.Errorf("output mismatch for %s\n--- got ---\n%s\n--- want ---\n%s", name, result, expected)
			}
		})
	}
}

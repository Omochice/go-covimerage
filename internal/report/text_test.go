package report

import (
	"bytes"
	"strings"
	"testing"

	"github.com/omochice/go-covimerage/internal/coverage"
)

func TestTextReporter_Generate(t *testing.T) {
	coverageData := &coverage.CoverageData{
		Scripts: map[string]*coverage.ScriptCoverage{
			"/path/to/script.vim": {
				Path:          "/path/to/script.vim",
				ExecutedLines: map[int]int{1: 5, 2: 3, 4: 2},
				MissingLines:  []int{3, 5},
				TotalLines:    5,
				CoveredLines:  3,
			},
		},
		Version: "1.0",
	}

	reporter := NewTextReporter(coverageData)
	var buf bytes.Buffer
	err := reporter.Generate(&buf)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()

	// Check that output contains script path
	if !strings.Contains(output, "/path/to/script.vim") {
		t.Error("Expected output to contain script path")
	}

	// Check that output contains coverage percentage
	if !strings.Contains(output, "60") { // 3/5 = 60%
		t.Error("Expected output to contain coverage percentage")
	}
}

func TestTextReporter_Generate_ShowMissing(t *testing.T) {
	coverageData := &coverage.CoverageData{
		Scripts: map[string]*coverage.ScriptCoverage{
			"/path/to/script.vim": {
				Path:          "/path/to/script.vim",
				ExecutedLines: map[int]int{1: 5, 2: 3},
				MissingLines:  []int{3, 4, 5},
				TotalLines:    5,
				CoveredLines:  2,
			},
		},
		Version: "1.0",
	}

	reporter := NewTextReporter(coverageData)
	reporter.SetShowMissing(true)

	var buf bytes.Buffer
	err := reporter.Generate(&buf)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()

	// Check that missing lines are shown
	if !strings.Contains(output, "3") || !strings.Contains(output, "4") || !strings.Contains(output, "5") {
		t.Error("Expected output to contain missing line numbers")
	}
}

func TestTextReporter_Generate_SkipCovered(t *testing.T) {
	coverageData := &coverage.CoverageData{
		Scripts: map[string]*coverage.ScriptCoverage{
			"/path/to/full.vim": {
				Path:          "/path/to/full.vim",
				ExecutedLines: map[int]int{1: 5, 2: 3, 3: 2},
				MissingLines:  []int{},
				TotalLines:    3,
				CoveredLines:  3,
			},
			"/path/to/partial.vim": {
				Path:          "/path/to/partial.vim",
				ExecutedLines: map[int]int{1: 5},
				MissingLines:  []int{2, 3},
				TotalLines:    3,
				CoveredLines:  1,
			},
		},
		Version: "1.0",
	}

	reporter := NewTextReporter(coverageData)
	reporter.SetSkipCovered(true)

	var buf bytes.Buffer
	err := reporter.Generate(&buf)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()

	// Should not contain the fully covered file
	if strings.Contains(output, "full.vim") {
		t.Error("Expected output to skip fully covered files")
	}

	// Should contain the partially covered file
	if !strings.Contains(output, "partial.vim") {
		t.Error("Expected output to contain partially covered file")
	}
}

func TestTextReporter_Generate_Summary(t *testing.T) {
	coverageData := &coverage.CoverageData{
		Scripts: map[string]*coverage.ScriptCoverage{
			"/path/to/first.vim": {
				Path:          "/path/to/first.vim",
				ExecutedLines: map[int]int{1: 5, 2: 3},
				MissingLines:  []int{3},
				TotalLines:    3,
				CoveredLines:  2,
			},
			"/path/to/second.vim": {
				Path:          "/path/to/second.vim",
				ExecutedLines: map[int]int{1: 2, 2: 1},
				MissingLines:  []int{},
				TotalLines:    2,
				CoveredLines:  2,
			},
		},
		Version: "1.0",
	}

	reporter := NewTextReporter(coverageData)
	var buf bytes.Buffer
	err := reporter.Generate(&buf)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()

	// Check that summary line exists
	if !strings.Contains(output, "TOTAL") {
		t.Error("Expected output to contain TOTAL summary line")
	}

	// Overall coverage should be 4/5 = 80%
	if !strings.Contains(output, "80") {
		t.Error("Expected output to contain overall coverage percentage")
	}
}

func TestTextReporter_WithFilters(t *testing.T) {
	coverageData := &coverage.CoverageData{
		Scripts: map[string]*coverage.ScriptCoverage{
			"/path/to/include.vim": {
				Path:          "/path/to/include.vim",
				ExecutedLines: map[int]int{1: 5},
				TotalLines:    5,
				CoveredLines:  1,
			},
			"/path/to/exclude.vim": {
				Path:          "/path/to/exclude.vim",
				ExecutedLines: map[int]int{1: 3},
				TotalLines:    3,
				CoveredLines:  1,
			},
		},
		Version: "1.0",
	}

	reporter := NewTextReporter(coverageData)
	reporter.SetIncludePatterns([]string{"*/include.vim"})

	var buf bytes.Buffer
	err := reporter.Generate(&buf)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()

	// Should include the include.vim file
	if !strings.Contains(output, "include.vim") {
		t.Error("Expected output to include include.vim")
	}

	// Should not include the exclude.vim file
	if strings.Contains(output, "exclude.vim") {
		t.Error("Expected output to exclude exclude.vim")
	}
}

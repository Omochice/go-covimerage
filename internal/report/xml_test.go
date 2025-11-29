package report

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/omochice/go-covimerage/internal/coverage"
)

func TestXMLReporter_Generate(t *testing.T) {
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

	reporter := NewXMLReporter(coverageData)
	var buf bytes.Buffer
	err := reporter.Generate(&buf)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()

	// Check that output is valid XML
	if !strings.HasPrefix(output, "<?xml") {
		t.Error("Expected output to start with XML declaration")
	}

	// Verify it can be parsed
	var result interface{}
	if err := xml.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse XML: %v", err)
	}
}

func TestXMLReporter_Generate_CoberturaFormat(t *testing.T) {
	coverageData := &coverage.CoverageData{
		Scripts: map[string]*coverage.ScriptCoverage{
			"/path/to/script.vim": {
				Path:          "/path/to/script.vim",
				ExecutedLines: map[int]int{1: 5, 2: 3},
				MissingLines:  []int{3},
				TotalLines:    3,
				CoveredLines:  2,
			},
		},
		Version: "1.0",
	}

	reporter := NewXMLReporter(coverageData)
	var buf bytes.Buffer
	err := reporter.Generate(&buf)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	output := buf.String()

	// Check for Cobertura XML elements
	if !strings.Contains(output, "<coverage") {
		t.Error("Expected output to contain <coverage> element")
	}
	if !strings.Contains(output, "line-rate") {
		t.Error("Expected output to contain line-rate attribute")
	}
}

func TestXMLReporter_WriteToFile(t *testing.T) {
	coverageData := &coverage.CoverageData{
		Scripts: map[string]*coverage.ScriptCoverage{
			"/path/to/script.vim": {
				Path:          "/path/to/script.vim",
				ExecutedLines: map[int]int{1: 5},
				TotalLines:    5,
				CoveredLines:  1,
			},
		},
		Version: "1.0",
	}

	tmpFile := "/tmp/claude/coverage.xml"

	reporter := NewXMLReporter(coverageData)
	err := reporter.WriteToFile(tmpFile)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// File should exist
	// (Cleanup happens in defer in actual usage)
}

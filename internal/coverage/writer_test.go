package coverage

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func TestWriter_WriteJSON(t *testing.T) {
	coverageData := &CoverageData{
		Scripts: map[string]*ScriptCoverage{
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

	writer := NewWriter(coverageData)
	var buf bytes.Buffer
	err := writer.WriteJSON(&buf)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify JSON is valid
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	// Check version
	if version, ok := result["version"].(string); !ok || version != "1.0" {
		t.Errorf("Expected version 1.0, got %v", result["version"])
	}

	// Check scripts exist
	if _, ok := result["scripts"]; !ok {
		t.Error("Expected scripts field in JSON")
	}
}

func TestWriter_WriteJSONToFile(t *testing.T) {
	coverageData := &CoverageData{
		Scripts: map[string]*ScriptCoverage{
			"/path/to/script.vim": {
				Path:          "/path/to/script.vim",
				ExecutedLines: map[int]int{1: 5},
				TotalLines:    5,
				CoveredLines:  1,
			},
		},
		Version: "1.0",
	}

	tmpFile := "/tmp/claude/test_coverage.json"
	defer os.Remove(tmpFile)

	writer := NewWriter(coverageData)
	err := writer.WriteJSONToFile(tmpFile)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify file exists and is valid JSON
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to parse JSON from file: %v", err)
	}
}

func TestWriter_WriteCoverage(t *testing.T) {
	coverageData := &CoverageData{
		Scripts: map[string]*ScriptCoverage{
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

	tmpFile := "/tmp/claude/test.coverage"
	defer os.Remove(tmpFile)

	writer := NewWriter(coverageData)
	err := writer.WriteCoverage(tmpFile, false)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tmpFile); os.IsNotExist(err) {
		t.Error("Expected coverage file to exist")
	}
}

func TestWriter_WriteCoverage_Append(t *testing.T) {
	// Create initial coverage data
	initialData := &CoverageData{
		Scripts: map[string]*ScriptCoverage{
			"/path/to/first.vim": {
				Path:          "/path/to/first.vim",
				ExecutedLines: map[int]int{1: 2},
				TotalLines:    2,
				CoveredLines:  1,
			},
		},
		Version: "1.0",
	}

	tmpFile := "/tmp/claude/test_append.coverage"
	defer os.Remove(tmpFile)

	// Write initial data
	writer1 := NewWriter(initialData)
	if err := writer1.WriteCoverage(tmpFile, false); err != nil {
		t.Fatalf("Failed to write initial data: %v", err)
	}

	// Append new data
	newData := &CoverageData{
		Scripts: map[string]*ScriptCoverage{
			"/path/to/second.vim": {
				Path:          "/path/to/second.vim",
				ExecutedLines: map[int]int{1: 3},
				TotalLines:    3,
				CoveredLines:  1,
			},
		},
		Version: "1.0",
	}

	writer2 := NewWriter(newData)
	if err := writer2.WriteCoverage(tmpFile, true); err != nil {
		t.Fatalf("Failed to append data: %v", err)
	}

	// Load and verify merged data
	merged, err := LoadCoverageData(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load coverage data: %v", err)
	}

	if len(merged.Scripts) != 2 {
		t.Errorf("Expected 2 scripts after append, got %d", len(merged.Scripts))
	}
}

func TestLoadCoverageData(t *testing.T) {
	// Create test coverage data
	original := &CoverageData{
		Scripts: map[string]*ScriptCoverage{
			"/path/to/script.vim": {
				Path:          "/path/to/script.vim",
				ExecutedLines: map[int]int{1: 5},
				TotalLines:    5,
				CoveredLines:  1,
			},
		},
		Version: "1.0",
	}

	tmpFile := "/tmp/claude/test_load.coverage"
	defer os.Remove(tmpFile)

	// Write data
	writer := NewWriter(original)
	if err := writer.WriteCoverage(tmpFile, false); err != nil {
		t.Fatalf("Failed to write data: %v", err)
	}

	// Load data
	loaded, err := LoadCoverageData(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	if len(loaded.Scripts) != len(original.Scripts) {
		t.Errorf("Expected %d scripts, got %d", len(original.Scripts), len(loaded.Scripts))
	}
}

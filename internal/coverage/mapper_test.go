package coverage

import (
	"testing"
	"time"

	"github.com/omochice/go-covimerage/internal/parser"
)

func TestMapper_MapFunctionToScript(t *testing.T) {
	// Create a function with execution data
	fn := &parser.Function{
		Name:      "TestFunc",
		Defined:   "/path/to/script.vim",
		StartLine: 5,
		Count:     2,
		Lines: []parser.Line{
			{Number: 1, Count: 2, SelfTime: 10 * time.Microsecond}, // Line 6 in script (StartLine + 1)
			{Number: 2, Count: 2, SelfTime: 20 * time.Microsecond}, // Line 7 in script
		},
	}

	// Create a script coverage
	sc := &ScriptCoverage{
		Path:          "/path/to/script.vim",
		ExecutedLines: make(map[int]int),
		MissingLines:  []int{},
	}

	mapper := NewMapper()
	err := mapper.MapFunctionToScript(fn, sc)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check that function lines are mapped to script lines
	// Function line 1 -> Script line 6 (StartLine + lineNum)
	if count, exists := sc.ExecutedLines[6]; !exists || count != 2 {
		t.Errorf("Expected script line 6 to have count 2, got %d (exists: %v)", count, exists)
	}

	// Function line 2 -> Script line 7
	if count, exists := sc.ExecutedLines[7]; !exists || count != 2 {
		t.Errorf("Expected script line 7 to have count 2, got %d (exists: %v)", count, exists)
	}
}

func TestMapper_MapFunctionToScript_AccumulateCounts(t *testing.T) {
	// Test that calling a function multiple times accumulates counts
	fn := &parser.Function{
		Name:      "TestFunc",
		Defined:   "/path/to/script.vim",
		StartLine: 5,
		Count:     3,
		Lines: []parser.Line{
			{Number: 1, Count: 3, SelfTime: 10 * time.Microsecond},
		},
	}

	sc := &ScriptCoverage{
		Path:          "/path/to/script.vim",
		ExecutedLines: make(map[int]int),
	}

	// Pre-populate with existing count
	sc.ExecutedLines[6] = 5

	mapper := NewMapper()
	err := mapper.MapFunctionToScript(fn, sc)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should accumulate: 5 (existing) + 3 (function) = 8
	if count := sc.ExecutedLines[6]; count != 8 {
		t.Errorf("Expected accumulated count 8, got %d", count)
	}
}

func TestMapper_MapFunctionToScript_NoStartLine(t *testing.T) {
	// Function without StartLine information (e.g., not from "Defined:" line)
	fn := &parser.Function{
		Name:      "TestFunc",
		Defined:   "",
		StartLine: 0,
		Count:     1,
		Lines: []parser.Line{
			{Number: 1, Count: 1, SelfTime: 10 * time.Microsecond},
		},
	}

	sc := &ScriptCoverage{
		Path:          "/path/to/script.vim",
		ExecutedLines: make(map[int]int),
	}

	mapper := NewMapper()
	err := mapper.MapFunctionToScript(fn, sc)

	// Should handle gracefully (skip mapping or use alternative strategy)
	if err != nil {
		// It's okay to return an error for unmappable functions
		return
	}

	// Or it might skip mapping - either is acceptable
	if len(sc.ExecutedLines) > 0 {
		t.Error("Expected no lines to be mapped without StartLine information")
	}
}

func TestMapper_CalculateMissingLines(t *testing.T) {
	sc := &ScriptCoverage{
		Path: "/path/to/script.vim",
		ExecutedLines: map[int]int{
			1: 1,
			2: 5,
			4: 2,
			5: 3,
		},
		TotalLines: 5,
	}

	mapper := NewMapper()
	mapper.CalculateMissingLines(sc)

	// Line 3 is missing (not in ExecutedLines)
	expectedMissing := []int{3}
	if len(sc.MissingLines) != len(expectedMissing) {
		t.Errorf("Expected %d missing lines, got %d", len(expectedMissing), len(sc.MissingLines))
	}

	for _, line := range expectedMissing {
		found := false
		for _, missing := range sc.MissingLines {
			if missing == line {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected line %d to be in missing lines", line)
		}
	}
}

func TestMapper_CalculateCoveredLines(t *testing.T) {
	sc := &ScriptCoverage{
		Path: "/path/to/script.vim",
		ExecutedLines: map[int]int{
			1: 1,
			2: 5,
			3: 0, // Not executed (count 0)
			4: 2,
		},
		TotalLines: 5,
	}

	mapper := NewMapper()
	mapper.CalculateCoveredLines(sc)

	// Only lines with count > 0 are covered
	expectedCovered := 3 // lines 1, 2, 4
	if sc.CoveredLines != expectedCovered {
		t.Errorf("Expected %d covered lines, got %d", expectedCovered, sc.CoveredLines)
	}
}

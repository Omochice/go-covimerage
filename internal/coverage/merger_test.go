package coverage

import (
	"testing"
	"time"

	"github.com/omochice/go-covimerage/internal/parser"
)

func TestMerger_Merge_SingleProfile(t *testing.T) {
	// Create a simple profile
	profile := &parser.Profile{
		Scripts: map[string]*parser.Script{
			"/path/to/script.vim": {
				Path: "/path/to/script.vim",
				Lines: map[int]*parser.Line{
					1: {Number: 1, Count: 5, SelfTime: 10 * time.Microsecond},
					2: {Number: 2, Count: 3, SelfTime: 20 * time.Microsecond},
				},
			},
		},
		Functions: []*parser.Function{},
	}

	merger := NewMerger([]*parser.Profile{profile})
	coverageData, err := merger.Merge()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(coverageData.Scripts) != 1 {
		t.Errorf("Expected 1 script, got %d", len(coverageData.Scripts))
	}

	sc := coverageData.Scripts["/path/to/script.vim"]
	if sc == nil {
		t.Fatal("Expected script coverage to exist")
	}

	if sc.ExecutedLines[1] != 5 {
		t.Errorf("Expected line 1 count 5, got %d", sc.ExecutedLines[1])
	}
	if sc.ExecutedLines[2] != 3 {
		t.Errorf("Expected line 2 count 3, got %d", sc.ExecutedLines[2])
	}
}

func TestMerger_Merge_MultipleProfiles(t *testing.T) {
	// Create two profiles with overlapping scripts
	profile1 := &parser.Profile{
		Scripts: map[string]*parser.Script{
			"/path/to/script.vim": {
				Path: "/path/to/script.vim",
				Lines: map[int]*parser.Line{
					1: {Number: 1, Count: 3, SelfTime: 10 * time.Microsecond},
					2: {Number: 2, Count: 2, SelfTime: 15 * time.Microsecond},
				},
			},
		},
		Functions: []*parser.Function{},
	}

	profile2 := &parser.Profile{
		Scripts: map[string]*parser.Script{
			"/path/to/script.vim": {
				Path: "/path/to/script.vim",
				Lines: map[int]*parser.Line{
					1: {Number: 1, Count: 2, SelfTime: 5 * time.Microsecond},
					3: {Number: 3, Count: 5, SelfTime: 20 * time.Microsecond},
				},
			},
		},
		Functions: []*parser.Function{},
	}

	merger := NewMerger([]*parser.Profile{profile1, profile2})
	coverageData, err := merger.Merge()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	sc := coverageData.Scripts["/path/to/script.vim"]
	if sc == nil {
		t.Fatal("Expected script coverage to exist")
	}

	// Counts should be merged (added together)
	if sc.ExecutedLines[1] != 5 { // 3 + 2
		t.Errorf("Expected line 1 count 5, got %d", sc.ExecutedLines[1])
	}
	if sc.ExecutedLines[2] != 2 {
		t.Errorf("Expected line 2 count 2, got %d", sc.ExecutedLines[2])
	}
	if sc.ExecutedLines[3] != 5 {
		t.Errorf("Expected line 3 count 5, got %d", sc.ExecutedLines[3])
	}
}

func TestMerger_Merge_WithFunctions(t *testing.T) {
	// Create a profile with functions
	profile := &parser.Profile{
		Scripts: map[string]*parser.Script{
			"/path/to/script.vim": {
				Path:  "/path/to/script.vim",
				Lines: map[int]*parser.Line{},
			},
		},
		Functions: []*parser.Function{
			{
				Name:      "TestFunc",
				Defined:   "/path/to/script.vim",
				StartLine: 10,
				Count:     2,
				Lines: []parser.Line{
					{Number: 1, Count: 2, SelfTime: 10 * time.Microsecond},
					{Number: 2, Count: 2, SelfTime: 15 * time.Microsecond},
				},
			},
		},
	}

	merger := NewMerger([]*parser.Profile{profile})
	coverageData, err := merger.Merge()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	sc := coverageData.Scripts["/path/to/script.vim"]
	if sc == nil {
		t.Fatal("Expected script coverage to exist")
	}

	// Function lines should be mapped to script lines
	// Function line 1 -> Script line 11 (StartLine + 1)
	// Function line 2 -> Script line 12 (StartLine + 2)
	if sc.ExecutedLines[11] != 2 {
		t.Errorf("Expected line 11 count 2, got %d", sc.ExecutedLines[11])
	}
	if sc.ExecutedLines[12] != 2 {
		t.Errorf("Expected line 12 count 2, got %d", sc.ExecutedLines[12])
	}
}

func TestMerger_Merge_EmptyProfiles(t *testing.T) {
	merger := NewMerger([]*parser.Profile{})
	coverageData, err := merger.Merge()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(coverageData.Scripts) != 0 {
		t.Errorf("Expected 0 scripts, got %d", len(coverageData.Scripts))
	}
}

func TestMerger_WithSourceFilter(t *testing.T) {
	profile := &parser.Profile{
		Scripts: map[string]*parser.Script{
			"/path/to/included.vim": {
				Path:  "/path/to/included.vim",
				Lines: map[int]*parser.Line{1: {Number: 1, Count: 1}},
			},
			"/other/excluded.vim": {
				Path:  "/other/excluded.vim",
				Lines: map[int]*parser.Line{1: {Number: 1, Count: 1}},
			},
		},
		Functions: []*parser.Function{},
	}

	merger := NewMerger([]*parser.Profile{profile})
	merger.SetSourceFilter([]string{"/path/to"})

	coverageData, err := merger.Merge()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Should only include scripts matching the filter
	if len(coverageData.Scripts) != 1 {
		t.Errorf("Expected 1 script after filtering, got %d", len(coverageData.Scripts))
	}

	if _, exists := coverageData.Scripts["/path/to/included.vim"]; !exists {
		t.Error("Expected /path/to/included.vim to be included")
	}
	if _, exists := coverageData.Scripts["/other/excluded.vim"]; exists {
		t.Error("Expected /other/excluded.vim to be excluded")
	}
}

package coverage

import (
	"testing"
)

func TestScriptCoverage_New(t *testing.T) {
	sc := &ScriptCoverage{
		Path:          "/path/to/script.vim",
		ExecutedLines: make(map[int]int),
		MissingLines:  []int{},
		TotalLines:    10,
		CoveredLines:  7,
	}

	sc.ExecutedLines[1] = 5
	sc.ExecutedLines[2] = 3
	sc.MissingLines = []int{8, 9, 10}

	if sc.Path != "/path/to/script.vim" {
		t.Errorf("Expected path /path/to/script.vim, got %s", sc.Path)
	}
	if len(sc.ExecutedLines) != 2 {
		t.Errorf("Expected 2 executed lines, got %d", len(sc.ExecutedLines))
	}
	if len(sc.MissingLines) != 3 {
		t.Errorf("Expected 3 missing lines, got %d", len(sc.MissingLines))
	}
}

func TestScriptCoverage_CoveragePercentage(t *testing.T) {
	tests := []struct {
		name         string
		totalLines   int
		coveredLines int
		want         float64
	}{
		{"full coverage", 10, 10, 100.0},
		{"partial coverage", 10, 7, 70.0},
		{"no coverage", 10, 0, 0.0},
		{"zero lines", 0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := &ScriptCoverage{
				TotalLines:   tt.totalLines,
				CoveredLines: tt.coveredLines,
			}
			got := sc.CoveragePercentage()
			if got != tt.want {
				t.Errorf("CoveragePercentage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCoverageData_New(t *testing.T) {
	cd := &CoverageData{
		Scripts: make(map[string]*ScriptCoverage),
		Version: "1.0",
	}

	sc := &ScriptCoverage{Path: "/path/to/script.vim"}
	cd.Scripts[sc.Path] = sc

	if cd.Version != "1.0" {
		t.Errorf("Expected version 1.0, got %s", cd.Version)
	}
	if len(cd.Scripts) != 1 {
		t.Errorf("Expected 1 script, got %d", len(cd.Scripts))
	}
}

func TestCoverageData_OverallCoverage(t *testing.T) {
	cd := &CoverageData{
		Scripts: map[string]*ScriptCoverage{
			"/path/to/first.vim": {
				TotalLines:   10,
				CoveredLines: 8,
			},
			"/path/to/second.vim": {
				TotalLines:   20,
				CoveredLines: 15,
			},
		},
	}

	total, covered, percentage := cd.OverallCoverage()

	if total != 30 {
		t.Errorf("Expected total 30, got %d", total)
	}
	if covered != 23 {
		t.Errorf("Expected covered 23, got %d", covered)
	}
	expectedPercentage := (23.0 / 30.0) * 100.0
	if percentage != expectedPercentage {
		t.Errorf("Expected percentage %.2f, got %.2f", expectedPercentage, percentage)
	}
}

func TestCoverageData_OverallCoverage_Empty(t *testing.T) {
	cd := &CoverageData{
		Scripts: make(map[string]*ScriptCoverage),
	}

	total, covered, percentage := cd.OverallCoverage()

	if total != 0 || covered != 0 || percentage != 0.0 {
		t.Errorf("Expected all zeros, got total=%d, covered=%d, percentage=%.2f", total, covered, percentage)
	}
}

package coverage

import (
	"strings"

	"github.com/omochice/go-covimerage/internal/parser"
)

// Merger merges multiple profiles into coverage data
type Merger struct {
	profiles     []*parser.Profile
	sourceFilter []string // Source filter paths
}

// NewMerger creates a new Merger instance
func NewMerger(profiles []*parser.Profile) *Merger {
	return &Merger{
		profiles:     profiles,
		sourceFilter: []string{},
	}
}

// SetSourceFilter sets source path filters
func (m *Merger) SetSourceFilter(sources []string) {
	m.sourceFilter = sources
}

// Merge combines profiles and generates coverage data
func (m *Merger) Merge() (*CoverageData, error) {
	coverageData := &CoverageData{
		Scripts: make(map[string]*ScriptCoverage),
		Version: "1.0",
	}

	mapper := NewMapper()

	// Process each profile
	for _, profile := range m.profiles {
		// Process scripts
		for path, script := range profile.Scripts {
			// Apply source filter if set
			if !m.matchesSourceFilter(path) {
				continue
			}

			// Get or create script coverage
			sc, exists := coverageData.Scripts[path]
			if !exists {
				sc = &ScriptCoverage{
					Path:          path,
					ExecutedLines: make(map[int]int),
					MissingLines:  []int{},
				}
				coverageData.Scripts[path] = sc
			}

			// Merge script line execution data
			for lineNum, line := range script.Lines {
				if existingCount, ok := sc.ExecutedLines[lineNum]; ok {
					sc.ExecutedLines[lineNum] = existingCount + line.Count
				} else {
					sc.ExecutedLines[lineNum] = line.Count
				}
			}
		}

		// Process functions and map them to scripts
		for _, fn := range profile.Functions {
			// Skip if function doesn't have definition info
			if fn.Defined == "" {
				continue
			}

			// Apply source filter
			if !m.matchesSourceFilter(fn.Defined) {
				continue
			}

			// Get or create script coverage for the function's file
			sc, exists := coverageData.Scripts[fn.Defined]
			if !exists {
				sc = &ScriptCoverage{
					Path:          fn.Defined,
					ExecutedLines: make(map[int]int),
					MissingLines:  []int{},
				}
				coverageData.Scripts[fn.Defined] = sc
			}

			// Map function execution to script lines
			if err := mapper.MapFunctionToScript(fn, sc); err != nil {
				return nil, err
			}
		}
	}

	// Calculate total lines and missing lines for each script
	for _, sc := range coverageData.Scripts {
		// Find the maximum line number to determine total lines
		maxLine := 0
		for lineNum := range sc.ExecutedLines {
			if lineNum > maxLine {
				maxLine = lineNum
			}
		}
		sc.TotalLines = maxLine

		// Calculate covered and missing lines
		mapper.CalculateCoveredLines(sc)
		mapper.CalculateMissingLines(sc)
	}

	return coverageData, nil
}

// matchesSourceFilter checks if a path matches the source filter
func (m *Merger) matchesSourceFilter(path string) bool {
	// If no filter is set, include everything
	if len(m.sourceFilter) == 0 {
		return true
	}

	// Check if path starts with any of the filter prefixes
	for _, filter := range m.sourceFilter {
		if strings.HasPrefix(path, filter) {
			return true
		}
	}

	return false
}

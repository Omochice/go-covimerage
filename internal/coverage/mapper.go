package coverage

import (
	"sort"

	"github.com/omochice/go-covimerage/internal/parser"
)

// Mapper maps function execution data to script coverage
type Mapper struct{}

// NewMapper creates a new Mapper instance
func NewMapper() *Mapper {
	return &Mapper{}
}

// MapFunctionToScript maps function line execution data to script lines
func (m *Mapper) MapFunctionToScript(fn *parser.Function, sc *ScriptCoverage) error {
	// Skip if function doesn't have start line information
	if fn.StartLine == 0 {
		return nil
	}

	// Skip if function is not defined in this script
	if fn.Defined != "" && fn.Defined != sc.Path {
		return nil
	}

	// Map each line in the function to the corresponding script line
	for _, line := range fn.Lines {
		// Calculate the actual line number in the script
		// Function line numbers are relative to the function start
		// Script line = StartLine + function line number
		scriptLine := fn.StartLine + line.Number

		// Accumulate the execution count
		if existingCount, exists := sc.ExecutedLines[scriptLine]; exists {
			sc.ExecutedLines[scriptLine] = existingCount + line.Count
		} else {
			sc.ExecutedLines[scriptLine] = line.Count
		}
	}

	return nil
}

// CalculateMissingLines determines which lines were not executed
func (m *Mapper) CalculateMissingLines(sc *ScriptCoverage) {
	if sc.TotalLines == 0 {
		sc.MissingLines = []int{}
		return
	}

	missing := make([]int, 0)

	// Check each line from 1 to TotalLines
	for i := 1; i <= sc.TotalLines; i++ {
		if _, exists := sc.ExecutedLines[i]; !exists {
			missing = append(missing, i)
		}
	}

	// Sort for consistent output
	sort.Ints(missing)
	sc.MissingLines = missing
}

// CalculateCoveredLines counts how many lines were executed
func (m *Mapper) CalculateCoveredLines(sc *ScriptCoverage) {
	covered := 0

	for _, count := range sc.ExecutedLines {
		if count > 0 {
			covered++
		}
	}

	sc.CoveredLines = covered
}

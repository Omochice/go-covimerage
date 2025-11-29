package coverage

// ScriptCoverage represents coverage for a single script
type ScriptCoverage struct {
	Path          string      // Path to the script file
	ExecutedLines map[int]int // Line number -> execution count
	MissingLines  []int       // Lines not executed
	TotalLines    int         // Total executable lines
	CoveredLines  int         // Lines executed at least once
}

// CoveragePercentage calculates the coverage percentage for this script
func (sc *ScriptCoverage) CoveragePercentage() float64 {
	if sc.TotalLines == 0 {
		return 0.0
	}
	return (float64(sc.CoveredLines) / float64(sc.TotalLines)) * 100.0
}

// CoverageData represents coverage information for all scripts
type CoverageData struct {
	Scripts map[string]*ScriptCoverage // Path -> Coverage
	Version string                      // Format version
}

// OverallCoverage calculates overall coverage statistics
func (cd *CoverageData) OverallCoverage() (totalLines, coveredLines int, percentage float64) {
	for _, script := range cd.Scripts {
		totalLines += script.TotalLines
		coveredLines += script.CoveredLines
	}

	if totalLines == 0 {
		return 0, 0, 0.0
	}

	percentage = (float64(coveredLines) / float64(totalLines)) * 100.0
	return totalLines, coveredLines, percentage
}

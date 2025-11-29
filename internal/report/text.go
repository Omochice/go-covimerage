package report

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	"github.com/omochice/go-covimerage/internal/coverage"
)

// TextReporter generates text-based coverage reports
type TextReporter struct {
	data            *coverage.CoverageData
	showMissing     bool
	skipCovered     bool
	includePatterns []string
	omitPatterns    []string
}

// NewTextReporter creates a new TextReporter instance
func NewTextReporter(data *coverage.CoverageData) *TextReporter {
	return &TextReporter{
		data:            data,
		showMissing:     false,
		skipCovered:     false,
		includePatterns: []string{},
		omitPatterns:    []string{},
	}
}

// SetShowMissing sets whether to show missing line numbers
func (r *TextReporter) SetShowMissing(show bool) {
	r.showMissing = show
}

// SetSkipCovered sets whether to skip 100% covered files
func (r *TextReporter) SetSkipCovered(skip bool) {
	r.skipCovered = skip
}

// SetIncludePatterns sets patterns for files to include
func (r *TextReporter) SetIncludePatterns(patterns []string) {
	r.includePatterns = patterns
}

// SetOmitPatterns sets patterns for files to omit
func (r *TextReporter) SetOmitPatterns(patterns []string) {
	r.omitPatterns = patterns
}

// Generate creates a text report
func (r *TextReporter) Generate(writer io.Writer) error {
	// Get sorted list of script paths
	paths := make([]string, 0, len(r.data.Scripts))
	for path := range r.data.Scripts {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// Write header
	fmt.Fprintf(writer, "%-50s %8s %8s %8s", "Name", "Stmts", "Miss", "Cover")
	if r.showMissing {
		fmt.Fprintf(writer, " %s", "Missing")
	}
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, strings.Repeat("-", 80))

	// Track totals
	totalLines := 0
	totalCovered := 0

	// Process each script
	for _, path := range paths {
		sc := r.data.Scripts[path]

		// Apply filters
		if !r.shouldInclude(path) {
			continue
		}

		// Skip fully covered files if requested
		if r.skipCovered && sc.CoveredLines == sc.TotalLines && sc.TotalLines > 0 {
			continue
		}

		// Calculate stats
		missed := sc.TotalLines - sc.CoveredLines
		percentage := sc.CoveragePercentage()

		// Update totals
		totalLines += sc.TotalLines
		totalCovered += sc.CoveredLines

		// Format output
		fmt.Fprintf(writer, "%-50s %8d %8d %7.0f%%", path, sc.TotalLines, missed, percentage)

		// Add missing lines if requested
		if r.showMissing && len(sc.MissingLines) > 0 {
			missingStr := formatLineNumbers(sc.MissingLines)
			fmt.Fprintf(writer, " %s", missingStr)
		}

		fmt.Fprintln(writer)
	}

	// Write summary line
	fmt.Fprintln(writer, strings.Repeat("-", 80))
	totalMissed := totalLines - totalCovered
	totalPercentage := 0.0
	if totalLines > 0 {
		totalPercentage = (float64(totalCovered) / float64(totalLines)) * 100.0
	}

	fmt.Fprintf(writer, "%-50s %8d %8d %7.0f%%\n", "TOTAL", totalLines, totalMissed, totalPercentage)

	return nil
}

// shouldInclude checks if a path should be included based on filters
func (r *TextReporter) shouldInclude(path string) bool {
	// Check omit patterns first
	for _, pattern := range r.omitPatterns {
		if matchPattern(pattern, path) {
			return false
		}
	}

	// If include patterns are specified, check them
	if len(r.includePatterns) > 0 {
		for _, pattern := range r.includePatterns {
			if matchPattern(pattern, path) {
				return true
			}
		}
		return false
	}

	return true
}

// matchPattern checks if a pattern matches a path
func matchPattern(pattern, path string) bool {
	// Try matching against full path
	if matched, _ := filepath.Match(pattern, path); matched {
		return true
	}

	// Try matching against basename
	if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
		return true
	}

	// Handle patterns like */file.vim by checking if the path ends with the pattern suffix
	if strings.HasPrefix(pattern, "*/") {
		suffix := pattern[2:] // Remove */
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}

	return false
}

// formatLineNumbers formats a list of line numbers for display
func formatLineNumbers(lines []int) string {
	if len(lines) == 0 {
		return ""
	}

	// Group consecutive lines into ranges
	ranges := []string{}
	start := lines[0]
	end := lines[0]

	for i := 1; i < len(lines); i++ {
		if lines[i] == end+1 {
			end = lines[i]
		} else {
			// Add the range
			if start == end {
				ranges = append(ranges, fmt.Sprintf("%d", start))
			} else {
				ranges = append(ranges, fmt.Sprintf("%d-%d", start, end))
			}
			start = lines[i]
			end = lines[i]
		}
	}

	// Add the last range
	if start == end {
		ranges = append(ranges, fmt.Sprintf("%d", start))
	} else {
		ranges = append(ranges, fmt.Sprintf("%d-%d", start, end))
	}

	return strings.Join(ranges, ", ")
}

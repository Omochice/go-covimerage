package parser

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

// Parser parses Vim profile output files
type Parser struct {
	profiles []*Profile
}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{
		profiles: make([]*Profile, 0),
	}
}

// parseScriptHeader parses the SCRIPT header section
func (p *Parser) parseScriptHeader(reader io.Reader) (*Script, error) {
	scanner := bufio.NewScanner(reader)
	var scriptPath string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "SCRIPT") {
			scriptPath = p.extractScriptPath(line)
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	script := &Script{
		Path:      scriptPath,
		Lines:     make(map[int]*Line),
		Functions: make(map[string]*Function),
	}

	return script, nil
}

// extractScriptPath extracts the script path from a SCRIPT line
func (p *Parser) extractScriptPath(line string) string {
	// Format: "SCRIPT  /path/to/script.vim" or "SCRIPT /path/to/script.vim"
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// parseLineData parses a line of execution data
// Format: "    count   total (s)   self (s)  code"
func (p *Parser) parseLineData(line string, lineNumber int) (*Line, error) {
	// Skip header line
	if strings.Contains(line, "count  total (s)   self (s)") {
		return nil, errors.New("header line")
	}

	result := &Line{
		Number: lineNumber,
	}

	// The line has fixed column positions
	// Count is at the beginning, then we need to check which time columns are present
	// We need to parse based on the actual spacing in the line

	// If the line is too short or doesn't start with spaces, it has no execution data
	if len(line) < 20 {
		return result, nil
	}

	// Extract fields - the format is space-separated but we need to be careful
	fields := strings.Fields(line)

	// If line has no execution data (empty or just code)
	if len(fields) == 0 || !isNumeric(fields[0]) {
		return result, nil
	}

	// Parse count (first field if numeric)
	result.Count = p.parseCount(fields[0])

	// Try to parse timing information
	// Check positions in the original line to determine which time column is present
	// Total time appears after count with some spacing
	// Self time appears after total time (or after count if no total)

	// Simple heuristic: if we have 2 time values, first is total, second is self
	// If we have 1 time value and it appears after significant spacing, it might be self
	timeFields := make([]string, 0)
	for i := 1; i < len(fields) && len(timeFields) < 2; i++ {
		if isTimeValue(fields[i]) {
			timeFields = append(timeFields, fields[i])
		}
	}

	if len(timeFields) == 2 {
		// Both total and self present
		result.TotalTime = p.parseTime(timeFields[0])
		result.SelfTime = p.parseTime(timeFields[1])
	} else if len(timeFields) == 1 {
		// Only one time value
		// Check position in original line to determine if it's total or self
		// The total time column is around position 8-16
		// The self time column is around position 20+
		timeIdx := strings.Index(line, timeFields[0])
		if timeIdx >= 8 && timeIdx < 18 {
			// Appears in total time column position
			result.TotalTime = p.parseTime(timeFields[0])
		} else {
			// Appears in self time column position
			result.SelfTime = p.parseTime(timeFields[0])
		}
	}

	return result, nil
}

// parseTime parses a time string in seconds to Duration
func (p *Parser) parseTime(timeStr string) time.Duration {
	if timeStr == "" {
		return 0
	}

	seconds, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		return 0
	}

	// Convert seconds to microseconds for better precision
	microseconds := int64(seconds * 1000000)
	return time.Duration(microseconds) * time.Microsecond
}

// parseCount parses a count string to int
func (p *Parser) parseCount(countStr string) int {
	countStr = strings.TrimSpace(countStr)
	if countStr == "" {
		return 0
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return 0
	}

	return count
}

// isNumeric checks if a string represents a number
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// isTimeValue checks if a string looks like a time value
func isTimeValue(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil && strings.Contains(s, ".")
}

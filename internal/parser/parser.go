package parser

import (
	"bufio"
	"errors"
	"io"
	"os"
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

// parseFunction parses a FUNCTION section
func (p *Parser) parseFunction(reader io.Reader) (*Function, error) {
	scanner := bufio.NewScanner(reader)
	fn := &Function{
		Lines: make([]Line, 0),
	}

	lineNum := 0
	inLines := false

	for scanner.Scan() {
		line := scanner.Text()

		// Parse FUNCTION header
		if strings.HasPrefix(line, "FUNCTION") {
			fn.Name = p.extractFunctionName(line)
			continue
		}

		// Parse "Defined:" line
		if strings.Contains(line, "Defined:") {
			fn.Defined, fn.StartLine = p.parseFunctionDefinition(line)
			continue
		}

		// Parse "Called N times" line
		if strings.HasPrefix(line, "Called") {
			fn.Count = p.parseCalledCount(line)
			continue
		}

		// Parse "Total time:" line
		if strings.HasPrefix(line, "Total time:") {
			fn.TotalTime = p.parseFunctionTime(line)
			continue
		}

		// Parse " Self time:" line
		if strings.Contains(line, "Self time:") {
			fn.SelfTime = p.parseFunctionTime(line)
			continue
		}

		// Skip header line
		if strings.Contains(line, "count  total (s)   self (s)") {
			inLines = true
			continue
		}

		// Parse line execution data
		if inLines && strings.TrimSpace(line) != "" {
			lineNum++
			lineData, err := p.parseLineData(line, lineNum)
			if err == nil && lineData != nil {
				fn.Lines = append(fn.Lines, *lineData)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return fn, nil
}

// extractFunctionName extracts the function name from a FUNCTION line
func (p *Parser) extractFunctionName(line string) string {
	// Format: "FUNCTION  FuncName()" or "FUNCTION FuncName()"
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}

// parseFunctionDefinition parses the "Defined:" line
// Format: "    Defined: /path/to/file.vim line 10"
func (p *Parser) parseFunctionDefinition(line string) (string, int) {
	// Remove leading spaces and "Defined:" prefix
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "Defined:")
	line = strings.TrimSpace(line)

	// Split by " line "
	parts := strings.Split(line, " line ")
	if len(parts) != 2 {
		return "", 0
	}

	path := strings.TrimSpace(parts[0])
	lineNum, _ := strconv.Atoi(strings.TrimSpace(parts[1]))

	return path, lineNum
}

// parseCalledCount parses the "Called N times" line
func (p *Parser) parseCalledCount(line string) int {
	// Format: "Called 2 times" or "Called 1 time"
	fields := strings.Fields(line)
	if len(fields) >= 2 {
		count, _ := strconv.Atoi(fields[1])
		return count
	}
	return 0
}

// parseFunctionTime parses time from "Total time:" or "Self time:" lines
// Format: "Total time:   0.000050" or " Self time:   0.000040"
func (p *Parser) parseFunctionTime(line string) time.Duration {
	// Remove prefix and get the time value
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "Total time:") {
		line = strings.TrimPrefix(line, "Total time:")
	} else if strings.HasPrefix(line, "Self time:") {
		line = strings.TrimPrefix(line, "Self time:")
	}

	timeStr := strings.TrimSpace(line)
	return p.parseTime(timeStr)
}

// Parse parses a complete Vim profile from a reader
func (p *Parser) Parse(reader io.Reader) (*Profile, error) {
	profile := &Profile{
		Scripts:   make(map[string]*Script),
		Functions: make([]*Function, 0),
	}

	scanner := bufio.NewScanner(reader)
	var currentScript *Script
	var currentFunction *Function
	var inScriptLines bool
	var inFunctionLines bool
	scriptLineNum := 0
	functionLineNum := 0

	for scanner.Scan() {
		line := scanner.Text()

		// Check for SCRIPT section
		if strings.HasPrefix(line, "SCRIPT") {
			// Save previous script if exists
			if currentScript != nil {
				profile.Scripts[currentScript.Path] = currentScript
			}

			// Start new script
			currentScript = &Script{
				Path:      p.extractScriptPath(line),
				Lines:     make(map[int]*Line),
				Functions: make(map[string]*Function),
			}
			inScriptLines = false
			inFunctionLines = false
			scriptLineNum = 0
			continue
		}

		// Check for FUNCTION section
		if strings.HasPrefix(line, "FUNCTION") {
			// Save previous function if exists
			if currentFunction != nil {
				profile.Functions = append(profile.Functions, currentFunction)
			}

			// Start new function
			currentFunction = &Function{
				Name:  p.extractFunctionName(line),
				Lines: make([]Line, 0),
			}
			inScriptLines = false
			inFunctionLines = false
			functionLineNum = 0
			continue
		}

		// Process SCRIPT section
		if currentScript != nil && currentFunction == nil {
			// Skip metadata lines
			if strings.HasPrefix(line, "Sourced") ||
				strings.HasPrefix(line, "Total time:") ||
				strings.Contains(line, "Self time:") {
				continue
			}

			// Check for line data header
			if strings.Contains(line, "count  total (s)   self (s)") {
				inScriptLines = true
				continue
			}

			// Parse script lines
			if inScriptLines && strings.TrimSpace(line) != "" {
				scriptLineNum++
				lineData, err := p.parseLineData(line, scriptLineNum)
				if err == nil && lineData != nil {
					currentScript.Lines[scriptLineNum] = lineData
				}
			}
		}

		// Process FUNCTION section
		if currentFunction != nil {
			// Parse "Defined:" line
			if strings.Contains(line, "Defined:") {
				currentFunction.Defined, currentFunction.StartLine = p.parseFunctionDefinition(line)
				continue
			}

			// Parse "Called N times" line
			if strings.HasPrefix(line, "Called") {
				currentFunction.Count = p.parseCalledCount(line)
				continue
			}

			// Parse "Total time:" line
			if strings.HasPrefix(line, "Total time:") {
				currentFunction.TotalTime = p.parseFunctionTime(line)
				continue
			}

			// Parse " Self time:" line
			if strings.Contains(line, "Self time:") {
				currentFunction.SelfTime = p.parseFunctionTime(line)
				continue
			}

			// Check for line data header
			if strings.Contains(line, "count  total (s)   self (s)") {
				inFunctionLines = true
				continue
			}

			// Parse function lines
			if inFunctionLines && strings.TrimSpace(line) != "" {
				functionLineNum++
				lineData, err := p.parseLineData(line, functionLineNum)
				if err == nil && lineData != nil {
					currentFunction.Lines = append(currentFunction.Lines, *lineData)
				}
			}
		}
	}

	// Save last script
	if currentScript != nil {
		profile.Scripts[currentScript.Path] = currentScript
	}

	// Save last function
	if currentFunction != nil {
		profile.Functions = append(profile.Functions, currentFunction)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return profile, nil
}

// ParseFile parses a Vim profile from a file
func (p *Parser) ParseFile(path string) (*Profile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return p.Parse(file)
}

// ParseString parses a Vim profile from a string
func (p *Parser) ParseString(content string) (*Profile, error) {
	return p.Parse(strings.NewReader(content))
}

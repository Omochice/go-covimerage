package parser

import "time"

// Line represents a single line in a Vim script with execution data
type Line struct {
	Number    int           // Line number (1-based)
	Count     int           // Execution count
	TotalTime time.Duration // Total execution time
	SelfTime  time.Duration // Self execution time (excluding called functions)
}

// Function represents a function definition in Vim script
type Function struct {
	Name      string        // Function name
	Defined   string        // File where defined
	StartLine int           // Starting line number
	EndLine   int           // Ending line number
	Count     int           // Call count
	TotalTime time.Duration // Total execution time
	SelfTime  time.Duration // Self execution time
	Lines     []Line        // Line-level execution data
}

// Script represents a Vim script file with coverage data
type Script struct {
	Path      string               // Absolute path to the script
	Lines     map[int]*Line        // Line number -> Line data
	Functions map[string]*Function // Function name -> Function data
}

// Profile represents parsed Vim profile output
type Profile struct {
	Scripts   map[string]*Script // Path -> Script
	Functions []*Function        // All functions (for mapping)
}

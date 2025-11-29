# Design: Go Implementation of covimerage (go-covimerage)

## Project Structure

```
go-covimerage/
├── cmd/
│   └── go-covimerage/
│       └── main.go              # CLI entry point
├── internal/
│   ├── parser/
│   │   ├── parser.go            # Profile parser
│   │   ├── parser_test.go
│   │   ├── line.go              # Line representation
│   │   └── script.go            # Script representation
│   ├── coverage/
│   │   ├── coverage.go          # Coverage data structure
│   │   ├── coverage_test.go
│   │   ├── merger.go            # Profile merger
│   │   └── writer.go            # Coverage data writer
│   ├── report/
│   │   ├── text.go              # Text report generator
│   │   ├── text_test.go
│   │   ├── xml.go               # Cobertura XML generator
│   │   └── xml_test.go
│   └── runner/
│       ├── runner.go            # Vim execution wrapper
│       └── runner_test.go
├── pkg/
│   └── covimerage/
│       └── covimerage.go        # Public API (if needed)
├── go.mod
├── go.sum
├── README.md
├── LICENSE
└── Makefile
```

## Core Data Structures

### 1. Profile Package (`internal/parser`)

```go
// Line represents a single line in a Vim script with execution data
type Line struct {
    Number    int       // Line number (1-based)
    Count     int       // Execution count
    TotalTime Duration  // Total execution time
    SelfTime  Duration  // Self execution time (excluding called functions)
}

// Function represents a function definition in Vim script
type Function struct {
    Name       string    // Function name
    Defined    string    // File where defined
    StartLine  int       // Starting line number
    EndLine    int       // Ending line number
    Count      int       // Call count
    TotalTime  Duration  // Total execution time
    SelfTime   Duration  // Self execution time
    Lines      []Line    // Line-level execution data
}

// Script represents a Vim script file with coverage data
type Script struct {
    Path       string              // Absolute path to the script
    Lines      map[int]*Line       // Line number -> Line data
    Functions  map[string]*Function // Function name -> Function data
}

// Profile represents parsed Vim profile output
type Profile struct {
    Scripts   map[string]*Script  // Path -> Script
    Functions []*Function         // All functions (for mapping)
}

// Parser parses Vim profile output files
type Parser struct {
    profiles []*Profile
}

// Parse parses a single profile file
func (p *Parser) Parse(reader io.Reader) (*Profile, error)

// ParseFile parses a profile file by path
func (p *Parser) ParseFile(path string) (*Profile, error)
```

### 2. Coverage Package (`internal/coverage`)

```go
// CoverageData represents coverage information for all scripts
type CoverageData struct {
    Scripts map[string]*ScriptCoverage // Path -> Coverage
    Version string                       // Format version
}

// ScriptCoverage represents coverage for a single script
type ScriptCoverage struct {
    Path          string
    ExecutedLines map[int]int  // Line number -> execution count
    MissingLines  []int         // Lines not executed
    TotalLines    int           // Total executable lines
    CoveredLines  int           // Lines executed at least once
}

// Merger merges multiple profiles into coverage data
type Merger struct {
    profiles []*parser.Profile
    sources  []string  // Source filter paths
}

// Merge combines profiles and generates coverage data
func (m *Merger) Merge() (*CoverageData, error)

// Writer writes coverage data to various formats
type Writer struct {
    data *CoverageData
}

// WriteJSON writes coverage data in JSON format
func (w *Writer) WriteJSON(writer io.Writer) error

// WriteCoverage writes in Coverage.py-compatible format (if feasible)
// OR custom format that can be converted
func (w *Writer) WriteCoverage(path string, append bool) error
```

### 3. Report Package (`internal/report`)

```go
// TextReporter generates text-based coverage reports
type TextReporter struct {
    data         *coverage.CoverageData
    showMissing  bool
    skipCovered  bool
    includeGlobs []string
    omitGlobs    []string
}

// Generate creates a text report
func (r *TextReporter) Generate(writer io.Writer) error

// XMLReporter generates Cobertura XML reports
type XMLReporter struct {
    data         *coverage.CoverageData
    includeGlobs []string
    omitGlobs    []string
    ignoreErrors bool
}

// Generate creates an XML report
func (r *XMLReporter) Generate(writer io.Writer) error

// Cobertura XML structure
type Coverage struct {
    XMLName       xml.Name  `xml:"coverage"`
    LineRate      float64   `xml:"line-rate,attr"`
    BranchRate    float64   `xml:"branch-rate,attr"`
    Version       string    `xml:"version,attr"`
    Timestamp     int64     `xml:"timestamp,attr"`
    LinesCovered  int       `xml:"lines-covered,attr"`
    LinesValid    int       `xml:"lines-valid,attr"`
    BranchesValid int       `xml:"branches-valid,attr"`
    Packages      []Package `xml:"packages>package"`
}

type Package struct {
    Name      string  `xml:"name,attr"`
    LineRate  float64 `xml:"line-rate,attr"`
    Classes   []Class `xml:"classes>class"`
}

type Class struct {
    Name     string  `xml:"name,attr"`
    Filename string  `xml:"filename,attr"`
    LineRate float64 `xml:"line-rate,attr"`
    Lines    []Line  `xml:"lines>line"`
}

type Line struct {
    Number int `xml:"number,attr"`
    Hits   int `xml:"hits,attr"`
}
```

### 4. Runner Package (`internal/runner`)

```go
// Runner wraps Vim execution with profiling
type Runner struct {
    vimPath      string    // Path to vim/nvim executable
    args         []string  // Arguments to pass to vim
    profileFile  string    // Output profile file path
    wrapProfile  bool      // Whether to add profiling commands
}

// Run executes Vim with profiling enabled
func (r *Runner) Run() error

// generateProfileCommands creates Vim commands for profiling
func (r *Runner) generateProfileCommands() string
```

## Processing Flow

### Flow 1: write-coverage Command

```
1. Parse command-line arguments
   ↓
2. For each profile file:
   - parser.ParseFile(profilePath)
   ↓
3. Create Merger with all profiles
   ↓
4. merger.Merge() → CoverageData
   ↓
5. Apply source filters
   ↓
6. If append flag:
   - Load existing coverage data
   - Merge with new data
   ↓
7. writer.WriteCoverage(dataFile)
```

### Flow 2: run Command

```
1. Parse command-line arguments
   ↓
2. Create Runner with Vim path and args
   ↓
3. If wrapProfile:
   - runner.Run() with profiling enabled
   Else:
   - Execute Vim directly
   ↓
4. If writeData:
   - Parse generated profile file
   - Generate coverage data
   - Write coverage data
   ↓
5. If report:
   - Load coverage data
   - Generate and display report
```

### Flow 3: report Command

```
1. Parse command-line arguments
   ↓
2. If profile files provided:
   - Parse profiles
   - Generate coverage data
   Else:
   - Load coverage data from file
   ↓
3. Create TextReporter with options
   ↓
4. reporter.Generate(outputWriter)
```

### Flow 4: xml Command

```
1. Parse command-line arguments
   ↓
2. Load coverage data from file
   ↓
3. Create XMLReporter with options
   ↓
4. reporter.Generate(outputFile)
```

## CLI Design

Using `cobra` library for CLI (or `flag` package for simplicity):

```go
// Command structure
go-covimerage [global-flags] <command> [command-flags] [args]

// Global flags
--verbose, -v          Increase verbosity
--quiet, -q            Decrease verbosity
--loglevel             Set log level (error/warn/info/debug)
--version, -V          Show version

// Commands
write-coverage         Parse profile and write coverage data
run                    Execute Vim with profiling
report                 Generate text coverage report
xml                    Generate XML coverage report
```

### Command Implementations

```go
package main

import (
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "go-covimerage",
    Short: "Generate code coverage for Vim scripts",
    Long:  `A tool to generate code coverage information for Vim scripts`,
}

var writeCoverageCmd = &cobra.Command{
    Use:   "write-coverage [profile-files...]",
    Short: "Parse profiles and write coverage data",
    Args:  cobra.MinimumNArgs(1),
    RunE:  runWriteCoverage,
}

var runCmd = &cobra.Command{
    Use:   "run [vim-args...]",
    Short: "Execute Vim with profiling",
    Args:  cobra.MinimumNArgs(1),
    RunE:  runVim,
}

var reportCmd = &cobra.Command{
    Use:   "report [profile-files...]",
    Short: "Generate coverage report",
    RunE:  runReport,
}

var xmlCmd = &cobra.Command{
    Use:   "xml",
    Short: "Generate XML coverage report",
    RunE:  runXML,
}
```

## Configuration File Format

Support a simple TOML configuration file (optional):

```toml
# .go-covimerage.toml

[coverage]
data_file = ".coverage"
source = ["src", "autoload"]

[report]
show_missing = true
skip_covered = false
omit = ["*/test/*", "*/tests/*"]

[xml]
output = "coverage.xml"
```

## Error Handling Strategy

1. **Parser errors**: Return wrapped errors with line context
2. **File I/O errors**: Return errors with file path information
3. **Vim execution errors**: Capture stderr and include in error message
4. **Invalid data**: Validate and return descriptive errors

```go
// Custom error types
type ParseError struct {
    File string
    Line int
    Msg  string
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("parse error in %s:%d: %s", e.File, e.Line, e.Msg)
}
```

## Testing Strategy

### Unit Tests
- Test each parser function with sample Vim profile output
- Test merger with multiple profiles
- Test coverage calculation accuracy
- Test report generation with known data

### Integration Tests
- Parse real Vim profile files from test fixtures
- Generate coverage for sample Vim scripts
- Verify XML output against schema
- Test CLI commands end-to-end

### Test Data
- Create fixture directory with:
  - Sample Vim scripts
  - Corresponding profile outputs
  - Expected coverage data
  - Expected reports (text and XML)

## Dependencies

### Required
- Go 1.21+ (standard library)
- `github.com/spf13/cobra` - CLI framework (optional, could use stdlib)
- `github.com/pelletier/go-toml/v2` - TOML config parsing (optional)

### Development
- `github.com/stretchr/testify` - Testing utilities
- Standard Go testing tools

## Build and Distribution

### Makefile targets
```makefile
build:          Build binary
test:           Run all tests
test-coverage:  Run tests with coverage
install:        Install binary
clean:          Clean build artifacts
lint:           Run linters
```

### Release process
- Tag versions (v0.1.0, v0.2.0, etc.)
- Use GitHub Actions or similar for cross-compilation
- Distribute via:
  - GitHub Releases (binaries)
  - Homebrew tap (macOS)
  - go install (source)

## Migration Notes

### Differences from Python version

1. **No Coverage.py integration**: Go version is standalone
2. **Custom data format**: May use JSON instead of Coverage.py's SQLite format
3. **Single binary**: No pip installation needed
4. **Configuration**: TOML instead of .coveragerc (different syntax)

### Compatibility

- Profile parsing: 100% compatible (same input format)
- XML output: 100% compatible (Cobertura standard)
- Text report: Similar format, may have minor differences
- Data file: Different format (provide migration tool if needed)

## Performance Considerations

1. **Streaming parsing**: Parse large files without loading entire file in memory
2. **Concurrent processing**: Parse multiple profile files in parallel
3. **Efficient data structures**: Use maps for O(1) lookups
4. **Minimal allocations**: Reuse buffers where possible

## Security Considerations

1. **Path traversal**: Validate file paths to prevent directory traversal
2. **Command injection**: Sanitize Vim arguments in runner
3. **Resource limits**: Set reasonable limits for file sizes and memory usage

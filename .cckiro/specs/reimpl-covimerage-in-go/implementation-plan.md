# Implementation Plan: Go Covimerage

## Overview

This implementation follows Test-Driven Development (TDD) with red-green-refactor cycle.
Each cycle will have 3 commits:
1. **RED**: Write failing test
2. **GREEN**: Implement minimum code to pass test
3. **REFACTOR**: Improve code quality (if needed)

Commit messages follow conventional-commit format with types based on user-facing impact.

## Phase 1: Project Setup

### Step 1.1: Initialize Go Module
**Type**: `build`
- Create go.mod
- Set up directory structure
- Add .gitignore

**Commits**:
1. `build: initialize Go module and project structure`

### Step 1.2: Add Testing Infrastructure
**Type**: `test`
- Set up test fixtures directory
- Create sample Vim profile output for testing
- Add Makefile with test targets

**Commits**:
1. `test: add testing infrastructure and fixtures`

## Phase 2: Parser Implementation

### Step 2.1: Basic Data Structures (TDD)
**Type**: `feat` (enables parsing functionality)

**RED Commit**: `test: add tests for basic parser data structures`
- Write tests for Line struct
- Write tests for Function struct
- Write tests for Script struct
- Write tests for Profile struct

**GREEN Commit**: `feat: implement basic parser data structures`
- Implement Line type
- Implement Function type
- Implement Script type
- Implement Profile type

**REFACTOR Commit** (if needed): `refactor: improve parser data structure design`

### Step 2.2: Profile Header Parsing (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for profile header parsing`
- Test parsing SCRIPT section headers
- Test parsing script path extraction
- Test parsing line count information

**GREEN Commit**: `feat: implement profile header parsing`
- Parse SCRIPT markers
- Extract script paths
- Parse "Sourced N times" information

**REFACTOR Commit** (if needed): `refactor: clean up header parsing logic`

### Step 2.3: Line Execution Data Parsing (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for line execution data parsing`
- Test parsing line counts
- Test parsing execution times
- Test handling various number formats

**GREEN Commit**: `feat: implement line execution data parsing`
- Parse "count total (s) self (s)" format
- Handle line number prefixes
- Parse timing information

**REFACTOR Commit** (if needed): `refactor: optimize line data parsing`

### Step 2.4: Function Parsing (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for function parsing`
- Test FUNCTION section parsing
- Test function definition location
- Test function timing data
- Test function line execution data

**GREEN Commit**: `feat: implement function parsing`
- Parse FUNCTION markers
- Extract function names and locations
- Parse function timing
- Parse function body lines

**REFACTOR Commit** (if needed): `refactor: improve function parsing structure`

### Step 2.5: Complete Profile Parsing (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for complete profile parsing`
- Test parsing full profile files
- Test handling multiple scripts
- Test handling multiple functions
- Test error cases (malformed input)

**GREEN Commit**: `feat: implement complete profile parser`
- Implement Parser.Parse() method
- Implement Parser.ParseFile() method
- Add error handling
- Handle edge cases

**REFACTOR Commit** (if needed): `refactor: refactor parser for better maintainability`

## Phase 3: Coverage Data Generation

### Step 3.1: Coverage Data Structures (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for coverage data structures`
- Test ScriptCoverage creation
- Test CoverageData creation
- Test coverage metrics calculation

**GREEN Commit**: `feat: implement coverage data structures`
- Implement ScriptCoverage type
- Implement CoverageData type
- Implement coverage calculations

**REFACTOR Commit** (if needed): `refactor: optimize coverage data structures`

### Step 3.2: Function Mapping (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for function-to-script mapping`
- Test mapping function lines to script lines
- Test handling overlapping functions
- Test handling re-defined functions

**GREEN Commit**: `feat: implement function-to-script mapping`
- Map function executions to script lines
- Handle function call counts
- Merge function data into script coverage

**REFACTOR Commit** (if needed): `refactor: improve mapping algorithm`

### Step 3.3: Profile Merging (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for profile merging`
- Test merging multiple profiles
- Test combining execution counts
- Test handling same file from multiple profiles

**GREEN Commit**: `feat: implement profile merger`
- Implement Merger.Merge() method
- Combine execution counts
- Handle duplicate scripts

**REFACTOR Commit** (if needed): `refactor: optimize merge performance`

### Step 3.4: Source Filtering (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for source filtering`
- Test filtering by source directories
- Test glob pattern matching
- Test include/omit patterns

**GREEN Commit**: `feat: implement source filtering`
- Implement path filtering
- Add glob pattern support
- Apply include/omit rules

**REFACTOR Commit** (if needed): `refactor: simplify filtering logic`

## Phase 4: Coverage Data Writer

### Step 4.1: JSON Writer (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for JSON coverage writer`
- Test JSON serialization
- Test writing to file
- Test append mode

**GREEN Commit**: `feat: implement JSON coverage writer`
- Implement WriteJSON method
- Handle file I/O
- Support append mode

**REFACTOR Commit** (if needed): `refactor: improve JSON writer error handling`

### Step 4.2: Coverage Format Writer (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for coverage format writer`
- Test custom coverage format
- Test data persistence
- Test loading existing data

**GREEN Commit**: `feat: implement coverage format writer`
- Define custom format
- Implement WriteCoverage method
- Support append operations

**REFACTOR Commit** (if needed): `refactor: optimize coverage file format`

## Phase 5: Report Generation

### Step 5.1: Text Report (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for text report generation`
- Test report formatting
- Test coverage percentage display
- Test missing lines display
- Test filtering options

**GREEN Commit**: `feat: implement text report generator`
- Implement TextReporter.Generate()
- Format coverage statistics
- Display missing lines
- Apply filters

**REFACTOR Commit** (if needed): `refactor: improve report formatting`

### Step 5.2: XML Report (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for XML report generation`
- Test Cobertura XML structure
- Test XML serialization
- Test coverage metrics in XML

**GREEN Commit**: `feat: implement XML report generator`
- Implement XMLReporter.Generate()
- Create Cobertura-compatible XML
- Calculate line rates
- Handle packages and classes

**REFACTOR Commit** (if needed): `refactor: optimize XML generation`

## Phase 6: Vim Runner

### Step 6.1: Profile Command Generation (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for profile command generation`
- Test generating Vim profiling commands
- Test profile file path handling

**GREEN Commit**: `feat: implement profile command generation`
- Generate :profile start commands
- Generate :profile func/file commands
- Handle profile file paths

**REFACTOR Commit** (if needed): `refactor: clean up command generation`

### Step 6.2: Vim Execution (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for Vim execution wrapper`
- Test Vim process execution
- Test passing arguments
- Test capturing output/errors

**GREEN Commit**: `feat: implement Vim execution wrapper`
- Implement Runner.Run() method
- Execute Vim with profiling
- Capture stderr/stdout
- Handle exit codes

**REFACTOR Commit** (if needed): `refactor: improve Vim execution error handling`

## Phase 7: CLI Implementation

### Step 7.1: CLI Framework Setup
**Type**: `build`

**Commits**:
1. `build: add cobra CLI framework dependency`
2. `build: set up CLI command structure`

### Step 7.2: write-coverage Command (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for write-coverage command`
- Test command parsing
- Test flag handling
- Test end-to-end execution

**GREEN Commit**: `feat: implement write-coverage command`
- Implement command handler
- Parse flags
- Integrate parser and writer

**REFACTOR Commit** (if needed): `refactor: simplify write-coverage command logic`

### Step 7.3: run Command (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for run command`
- Test Vim execution
- Test profile generation
- Test automatic reporting

**GREEN Commit**: `feat: implement run command`
- Implement command handler
- Integrate runner
- Add auto-report feature

**REFACTOR Commit** (if needed): `refactor: improve run command flow`

### Step 7.4: report Command (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for report command`
- Test report generation
- Test filtering options
- Test output formatting

**GREEN Commit**: `feat: implement report command`
- Implement command handler
- Integrate text reporter
- Handle filtering

**REFACTOR Commit** (if needed): `refactor: clean up report command`

### Step 7.5: xml Command (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for xml command`
- Test XML generation
- Test output file handling

**GREEN Commit**: `feat: implement xml command`
- Implement command handler
- Integrate XML reporter
- Handle output

**REFACTOR Commit** (if needed): `refactor: optimize xml command`

### Step 7.6: Global Flags and Options
**Type**: `feat`

**RED Commit**: `test: add tests for global flags`
- Test verbose/quiet flags
- Test log level setting
- Test version display

**GREEN Commit**: `feat: implement global flags and options`
- Add verbose/quiet handling
- Add log level configuration
- Add version command

**REFACTOR Commit** (if needed): `refactor: consolidate flag handling`

## Phase 8: Configuration File Support

### Step 8.1: Config File Parsing (TDD)
**Type**: `feat`

**RED Commit**: `test: add tests for config file parsing`
- Test TOML parsing
- Test config merging with CLI flags

**GREEN Commit**: `feat: implement config file support`
- Parse TOML config
- Merge with CLI flags
- Handle default values

**REFACTOR Commit** (if needed): `refactor: simplify config handling`

## Phase 9: Integration and Polish

### Step 9.1: Integration Tests
**Type**: `test`

**Commits**:
1. `test: add end-to-end integration tests`

### Step 9.2: Documentation
**Type**: `docs`

**Commits**:
1. `docs: add README with installation and usage`
2. `docs: add code documentation and examples`

### Step 9.3: Build System
**Type**: `build`

**Commits**:
1. `build: add Makefile with build targets`
2. `build: add cross-compilation support`

### Step 9.4: Error Messages and UX
**Type**: `feat`

**RED Commit**: `test: add tests for error messages`
**GREEN Commit**: `feat: improve error messages and user experience`
**REFACTOR Commit** (if needed): `refactor: enhance error handling`

## Phase 10: Performance Optimization

### Step 10.1: Profiling and Benchmarks
**Type**: `perf`

**Commits**:
1. `test: add performance benchmarks`
2. `perf: optimize parser for large files`
3. `perf: add concurrent profile parsing`

## Phase 11: Final Polish

### Step 11.1: Code Quality
**Type**: `refactor`

**Commits**:
1. `refactor: apply linter suggestions`
2. `refactor: improve code documentation`

### Step 11.2: Release Preparation
**Type**: `build`

**Commits**:
1. `build: add release workflow`
2. `build: prepare v0.1.0 release`

## Commit Type Guidelines

Based on user-facing impact:

- **feat**: New features or commands visible to users
- **fix**: Bug fixes that affect user experience
- **perf**: Performance improvements users can notice
- **docs**: Documentation changes
- **test**: Adding or updating tests (no user-facing changes)
- **build**: Build system, dependencies, CI/CD
- **refactor**: Code improvements without behavior changes
- **style**: Code formatting (rarely used, prefer refactor)

## Testing Strategy per Phase

Each phase includes:
1. **Unit tests**: Test individual functions/methods
2. **Integration tests**: Test component interactions
3. **Fixtures**: Real-world Vim profile samples
4. **Edge cases**: Error conditions, malformed input, boundary values

## Success Metrics

- [ ] All tests pass
- [ ] Test coverage > 80%
- [ ] Can parse real Vim profile output
- [ ] Generated reports match expected format
- [ ] All CLI commands work end-to-end
- [ ] Documentation is complete
- [ ] Binary can be easily distributed

## Estimated Implementation Order

1. **Phase 1-2**: Core parsing (foundation)
2. **Phase 3-4**: Coverage generation (core functionality)
3. **Phase 5**: Reporting (user-visible output)
4. **Phase 6**: Vim runner (convenience feature)
5. **Phase 7**: CLI (user interface)
6. **Phase 8**: Config (nice-to-have)
7. **Phase 9-11**: Polish and release

## Notes

- Each TDD cycle should be small and focused
- Commit early and often
- Keep commits atomic and well-described
- Run tests before each commit
- Update documentation as features are added
- Seek feedback at major milestones

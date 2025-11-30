# go-covimerage

Generate code coverage information for Vim scripts.

This is a Go reimplementation of [covimerage](https://github.com/Vimjas/covimerage), providing a fast, standalone binary for coverage analysis of Vim script code.

## Features

- ✅ Parse Vim profile output (`:profile`)
- ✅ Generate coverage data in JSON format
- ✅ Merge multiple profile files
- ✅ Text and XML (Cobertura) report generation
- ✅ Fast, single-binary distribution
- ✅ No Python dependencies

## Installation

```bash
go install github.com/omochice/go-covimerage/cmd/go-covimerage@latest
```

Or build from source:

```bash
git clone https://github.com/omochice/go-covimerage
cd go-covimerage
make build
```

## Usage

### Basic Workflow

1. Generate Vim profile output:

```vim
:profile start profile.txt
:profile file */path/to/*.vim
:profile func *
" Run your tests here
:qall!
```

2. Generate coverage data:

```bash
go-covimerage write-coverage profile.txt
```

This creates a `.coverage` file containing coverage information.

### Commands

#### write-coverage

Parse Vim profile files and generate coverage data.

```bash
# Basic usage
go-covimerage write-coverage profile.txt

# Multiple profile files
go-covimerage write-coverage profile1.txt profile2.txt

# Custom output file
go-covimerage write-coverage --data-file=my-coverage.json profile.txt

# Append to existing coverage
go-covimerage write-coverage --append profile.txt
```

**Options:**
- `--data-file` - Output file path (default: `.coverage`)
- `--append` - Append to existing coverage data

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for a specific package
go test ./internal/parser/
```

### Building

```bash
# Build binary
make build

# Install locally
make install
```

### Project Structure

```
go-covimerage/
├── cmd/go-covimerage/     # CLI entry point
├── internal/
│   ├── parser/            # Vim profile parser
│   ├── coverage/          # Coverage data generation
│   └── report/            # Report generators (text, XML)
├── testdata/              # Test fixtures
└── Makefile
```

## Implementation Status

**Completed:**
- ✅ Vim profile parser
- ✅ Coverage data generation and merging
- ✅ JSON coverage writer
- ✅ Text report generator
- ✅ XML (Cobertura) report generator
- ✅ CLI with `write-coverage` command

**Future Work:**
- ⏳ `run` command (wrap Vim execution with profiling)
- ⏳ `report` command (generate text reports from coverage data)
- ⏳ `xml` command (generate XML reports from coverage data)
- ⏳ Configuration file support
- ⏳ Advanced filtering and source mapping

## Contributing

This project follows Test-Driven Development (TDD) with red-green-refactor cycles. All contributions should:

1. Include tests for new functionality
2. Follow Go best practices and idioms
3. Use conventional commit messages
4. Maintain >80% test coverage

## License

[zlib](./LICENSE)

<details>
    <summary>The original covimerage license is as follows:</summary>

<https://github.com/Vimjas/covimerage/blob/51e59b232015fe5fb8ea87da4098aaefab1c9631/LICENSE>

```txt
Copyright 2017 Daniel Hahler

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

</details>

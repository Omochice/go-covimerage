# go-covimerage

Generate code coverage information for Vim scripts.

This is a Go reimplementation of [covimerage](https://github.com/Vimjas/covimerage), providing a fast, standalone binary for coverage analysis of Vim script code.

## Installation

```bash
go install github.com/omochice/go-covimerage/cmd/go-covimerage@latest
```

## Usage

```bash
# Run Vim with profiling and generate coverage
go-covimerage run vim -u test/vimrc -c 'Vader! test/*.vader'

# Parse profile and write coverage data
go-covimerage write-coverage profile.txt

# Generate text report
go-covimerage report

# Generate XML report
go-covimerage xml
```

## Development

```bash
# Run tests
make test

# Build
make build

# Run tests with coverage
make test-coverage
```

## License

MIT

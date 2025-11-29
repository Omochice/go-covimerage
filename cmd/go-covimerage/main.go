package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return showUsage()
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "write-coverage":
		return writeCoverageCommand(commandArgs)
	case "version", "--version", "-v":
		fmt.Printf("go-covimerage version %s\n", version)
		return nil
	case "help", "--help", "-h":
		return showUsage()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func showUsage() error {
	usage := `go-covimerage - Generate code coverage for Vim scripts

Usage:
  go-covimerage <command> [options]

Commands:
  write-coverage  Parse profile and write coverage data
  version         Show version information
  help            Show this help message

Examples:
  go-covimerage write-coverage profile.txt
  go-covimerage write-coverage --data-file=.coverage profile1.txt profile2.txt
  go-covimerage write-coverage --append profile.txt

For more information, see: https://github.com/omochice/go-covimerage
`
	fmt.Print(usage)
	return nil
}

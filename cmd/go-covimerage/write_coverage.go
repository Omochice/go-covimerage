package main

import (
	"flag"
	"fmt"

	"github.com/omochice/go-covimerage/internal/coverage"
	"github.com/omochice/go-covimerage/internal/parser"
)

// writeCoverageCommand implements the write-coverage command
func writeCoverageCommand(args []string) error {
	fs := flag.NewFlagSet("write-coverage", flag.ExitOnError)
	dataFile := fs.String("data-file", ".coverage", "Output data file")
	appendMode := fs.Bool("append", false, "Append to existing data file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	profileFiles := fs.Args()
	if len(profileFiles) == 0 {
		return fmt.Errorf("at least one profile file is required")
	}

	// Parse all profile files
	profiles := make([]*parser.Profile, 0, len(profileFiles))
	p := parser.NewParser()

	for _, profileFile := range profileFiles {
		profile, err := p.ParseFile(profileFile)
		if err != nil {
			return fmt.Errorf("failed to parse %s: %w", profileFile, err)
		}
		profiles = append(profiles, profile)
	}

	// Merge profiles into coverage data
	merger := coverage.NewMerger(profiles)
	coverageData, err := merger.Merge()
	if err != nil {
		return fmt.Errorf("failed to merge profiles: %w", err)
	}

	// Write coverage data
	writer := coverage.NewWriter(coverageData)
	if err := writer.WriteCoverage(*dataFile, *appendMode); err != nil {
		return fmt.Errorf("failed to write coverage data: %w", err)
	}

	fmt.Printf("Coverage data written to %s\n", *dataFile)
	return nil
}

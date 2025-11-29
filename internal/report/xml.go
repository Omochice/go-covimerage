package report

import (
	"encoding/xml"
	"io"
	"os"
	"sort"
	"time"

	"github.com/omochice/go-covimerage/internal/coverage"
)

// XMLReporter generates Cobertura XML reports
type XMLReporter struct {
	data *coverage.CoverageData
}

// NewXMLReporter creates a new XMLReporter instance
func NewXMLReporter(data *coverage.CoverageData) *XMLReporter {
	return &XMLReporter{data: data}
}

// Coverage represents the root Cobertura XML element
type Coverage struct {
	XMLName       xml.Name  `xml:"coverage"`
	LineRate      float64   `xml:"line-rate,attr"`
	BranchRate    float64   `xml:"branch-rate,attr"`
	Version       string    `xml:"version,attr"`
	Timestamp     int64     `xml:"timestamp,attr"`
	LinesCovered  int       `xml:"lines-covered,attr"`
	LinesValid    int       `xml:"lines-valid,attr"`
	BranchesCovered int     `xml:"branches-covered,attr"`
	BranchesValid int       `xml:"branches-valid,attr"`
	Packages      []Package `xml:"packages>package"`
}

// Package represents a package in Cobertura XML
type Package struct {
	Name      string  `xml:"name,attr"`
	LineRate  float64 `xml:"line-rate,attr"`
	BranchRate float64 `xml:"branch-rate,attr"`
	Complexity float64 `xml:"complexity,attr"`
	Classes   []Class `xml:"classes>class"`
}

// Class represents a class (file) in Cobertura XML
type Class struct {
	Name       string  `xml:"name,attr"`
	Filename   string  `xml:"filename,attr"`
	LineRate   float64 `xml:"line-rate,attr"`
	BranchRate float64 `xml:"branch-rate,attr"`
	Complexity float64 `xml:"complexity,attr"`
	Lines      []Line  `xml:"lines>line"`
}

// Line represents a line in Cobertura XML
type Line struct {
	Number int  `xml:"number,attr"`
	Hits   int  `xml:"hits,attr"`
	Branch bool `xml:"branch,attr"`
}

// Generate creates an XML report
func (r *XMLReporter) Generate(writer io.Writer) error {
	// Calculate overall statistics
	totalLines, coveredLines, _ := r.data.OverallCoverage()
	lineRate := 0.0
	if totalLines > 0 {
		lineRate = float64(coveredLines) / float64(totalLines)
	}

	// Create coverage structure
	cov := Coverage{
		LineRate:        lineRate,
		BranchRate:      0.0, // Not tracking branches
		Version:         r.data.Version,
		Timestamp:       time.Now().Unix(),
		LinesCovered:    coveredLines,
		LinesValid:      totalLines,
		BranchesCovered: 0,
		BranchesValid:   0,
		Packages:        []Package{},
	}

	// Group scripts into a single package
	pkg := Package{
		Name:       "vim-scripts",
		LineRate:   lineRate,
		BranchRate: 0.0,
		Complexity: 0.0,
		Classes:    []Class{},
	}

	// Get sorted paths
	paths := make([]string, 0, len(r.data.Scripts))
	for path := range r.data.Scripts {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	// Process each script
	for _, path := range paths {
		sc := r.data.Scripts[path]

		// Calculate line rate for this script
		scriptLineRate := 0.0
		if sc.TotalLines > 0 {
			scriptLineRate = float64(sc.CoveredLines) / float64(sc.TotalLines)
		}

		// Create class (file) entry
		class := Class{
			Name:       path,
			Filename:   path,
			LineRate:   scriptLineRate,
			BranchRate: 0.0,
			Complexity: 0.0,
			Lines:      []Line{},
		}

		// Add line entries
		lineNumbers := make([]int, 0, len(sc.ExecutedLines))
		for lineNum := range sc.ExecutedLines {
			lineNumbers = append(lineNumbers, lineNum)
		}
		sort.Ints(lineNumbers)

		for _, lineNum := range lineNumbers {
			hits := sc.ExecutedLines[lineNum]
			class.Lines = append(class.Lines, Line{
				Number: lineNum,
				Hits:   hits,
				Branch: false,
			})
		}

		// Add missing lines with 0 hits
		for _, lineNum := range sc.MissingLines {
			class.Lines = append(class.Lines, Line{
				Number: lineNum,
				Hits:   0,
				Branch: false,
			})
		}

		pkg.Classes = append(pkg.Classes, class)
	}

	cov.Packages = append(cov.Packages, pkg)

	// Write XML with declaration
	writer.Write([]byte(xml.Header))
	encoder := xml.NewEncoder(writer)
	encoder.Indent("", "  ")
	return encoder.Encode(cov)
}

// WriteToFile writes the XML report to a file
func (r *XMLReporter) WriteToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return r.Generate(file)
}

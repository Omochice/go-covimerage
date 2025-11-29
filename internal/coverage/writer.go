package coverage

import (
	"encoding/json"
	"io"
	"os"
)

// Writer writes coverage data to various formats
type Writer struct {
	data *CoverageData
}

// NewWriter creates a new Writer instance
func NewWriter(data *CoverageData) *Writer {
	return &Writer{data: data}
}

// WriteJSON writes coverage data in JSON format
func (w *Writer) WriteJSON(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(w.data)
}

// WriteJSONToFile writes coverage data to a JSON file
func (w *Writer) WriteJSONToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return w.WriteJSON(file)
}

// WriteCoverage writes in custom coverage format
// If append is true, merges with existing data
func (w *Writer) WriteCoverage(path string, appendMode bool) error {
	var mergedData *CoverageData

	if appendMode {
		// Load existing data if it exists
		existing, err := LoadCoverageData(path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if existing != nil {
			// Merge existing and new data
			mergedData = mergeCoverageData(existing, w.data)
		} else {
			mergedData = w.data
		}
	} else {
		mergedData = w.data
	}

	// Write as JSON for now (can be customized later)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(mergedData)
}

// LoadCoverageData loads coverage data from a file
func LoadCoverageData(path string) (*CoverageData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data CoverageData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

// mergeCoverageData merges two coverage data sets
func mergeCoverageData(existing, new *CoverageData) *CoverageData {
	merged := &CoverageData{
		Scripts: make(map[string]*ScriptCoverage),
		Version: new.Version,
	}

	// Copy existing scripts
	for path, sc := range existing.Scripts {
		merged.Scripts[path] = &ScriptCoverage{
			Path:          sc.Path,
			ExecutedLines: make(map[int]int),
			MissingLines:  append([]int{}, sc.MissingLines...),
			TotalLines:    sc.TotalLines,
			CoveredLines:  sc.CoveredLines,
		}
		for line, count := range sc.ExecutedLines {
			merged.Scripts[path].ExecutedLines[line] = count
		}
	}

	// Merge new scripts
	for path, sc := range new.Scripts {
		if existingSc, exists := merged.Scripts[path]; exists {
			// Merge execution counts
			for line, count := range sc.ExecutedLines {
				if existingCount, ok := existingSc.ExecutedLines[line]; ok {
					existingSc.ExecutedLines[line] = existingCount + count
				} else {
					existingSc.ExecutedLines[line] = count
				}
			}

			// Recalculate statistics
			mapper := NewMapper()
			mapper.CalculateCoveredLines(existingSc)
			mapper.CalculateMissingLines(existingSc)
		} else {
			// Add new script
			merged.Scripts[path] = &ScriptCoverage{
				Path:          sc.Path,
				ExecutedLines: make(map[int]int),
				MissingLines:  append([]int{}, sc.MissingLines...),
				TotalLines:    sc.TotalLines,
				CoveredLines:  sc.CoveredLines,
			}
			for line, count := range sc.ExecutedLines {
				merged.Scripts[path].ExecutedLines[line] = count
			}
		}
	}

	return merged
}

package parser

import (
	"strings"
	"testing"
)

func TestParser_ParseScriptHeader(t *testing.T) {
	input := `SCRIPT  /path/to/script.vim
Sourced 1 time
Total time:   0.000100
 Self time:   0.000080`

	parser := NewParser()
	script, err := parser.parseScriptHeader(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if script == nil {
		t.Fatal("Expected script to be non-nil")
	}
	if script.Path != "/path/to/script.vim" {
		t.Errorf("Expected path /path/to/script.vim, got %s", script.Path)
	}
}

func TestParser_ParseScriptHeader_MultipleSourced(t *testing.T) {
	input := `SCRIPT  /path/to/script.vim
Sourced 5 times
Total time:   0.000500
 Self time:   0.000400`

	parser := NewParser()
	script, err := parser.parseScriptHeader(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if script.Path != "/path/to/script.vim" {
		t.Errorf("Expected path /path/to/script.vim, got %s", script.Path)
	}
}

func TestParser_ExtractScriptPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"SCRIPT  /path/to/script.vim", "/path/to/script.vim"},
		{"SCRIPT /another/path.vim", "/another/path.vim"},
		{"SCRIPT  ~/user/script.vim", "~/user/script.vim"},
	}

	parser := NewParser()
	for _, tt := range tests {
		result := parser.extractScriptPath(tt.input)
		if result != tt.expected {
			t.Errorf("extractScriptPath(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

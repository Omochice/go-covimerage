package parser

import (
	"strings"
	"testing"
	"time"
)

func TestParser_ParseFunction(t *testing.T) {
	input := `FUNCTION  SimpleFunc()
Called 2 times
Total time:   0.000050
 Self time:   0.000040

count  total (s)   self (s)
    2              0.000020   echo "hello"
    2              0.000020   return 42`

	parser := NewParser()
	fn, err := parser.parseFunction(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if fn == nil {
		t.Fatal("Expected function to be non-nil")
	}
	if fn.Name != "SimpleFunc()" {
		t.Errorf("Expected name SimpleFunc(), got %s", fn.Name)
	}
	if fn.Count != 2 {
		t.Errorf("Expected count 2, got %d", fn.Count)
	}
	if fn.TotalTime != 50*time.Microsecond {
		t.Errorf("Expected total time 50µs, got %v", fn.TotalTime)
	}
	if fn.SelfTime != 40*time.Microsecond {
		t.Errorf("Expected self time 40µs, got %v", fn.SelfTime)
	}
	if len(fn.Lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(fn.Lines))
	}
}

func TestParser_ParseFunction_DictFunction(t *testing.T) {
	input := `FUNCTION  obj.dictFunc()
    Defined: /path/to/script.vim line 10
Called 1 time
Total time:   0.000010
 Self time:   0.000005

count  total (s)   self (s)
    1              0.000005   return 1`

	parser := NewParser()
	fn, err := parser.parseFunction(strings.NewReader(input))

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if fn.Name != "obj.dictFunc()" {
		t.Errorf("Expected name obj.dictFunc(), got %s", fn.Name)
	}
	if fn.Defined != "/path/to/script.vim" {
		t.Errorf("Expected defined /path/to/script.vim, got %s", fn.Defined)
	}
	if fn.StartLine != 10 {
		t.Errorf("Expected start line 10, got %d", fn.StartLine)
	}
}

func TestParser_ExtractFunctionName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"FUNCTION  SimpleFunc()", "SimpleFunc()"},
		{"FUNCTION TestFunc()", "TestFunc()"},
		{"FUNCTION  <SNR>123_InternalFunc()", "<SNR>123_InternalFunc()"},
		{"FUNCTION  obj.dictFunc()", "obj.dictFunc()"},
	}

	parser := NewParser()
	for _, tt := range tests {
		result := parser.extractFunctionName(tt.input)
		if result != tt.expected {
			t.Errorf("extractFunctionName(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestParser_ParseFunctionDefinition(t *testing.T) {
	tests := []struct {
		input         string
		wantPath      string
		wantStartLine int
	}{
		{
			input:         "    Defined: /path/to/script.vim line 10",
			wantPath:      "/path/to/script.vim",
			wantStartLine: 10,
		},
		{
			input:         "    Defined: ~/user/test.vim line 42",
			wantPath:      "~/user/test.vim",
			wantStartLine: 42,
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			path, line := parser.parseFunctionDefinition(tt.input)
			if path != tt.wantPath {
				t.Errorf("path = %q, want %q", path, tt.wantPath)
			}
			if line != tt.wantStartLine {
				t.Errorf("line = %d, want %d", line, tt.wantStartLine)
			}
		})
	}
}

func TestParser_ParseCalledCount(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"Called 1 time", 1},
		{"Called 2 times", 2},
		{"Called 100 times", 100},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parser.parseCalledCount(tt.input)
			if result != tt.want {
				t.Errorf("parseCalledCount(%q) = %d, want %d", tt.input, result, tt.want)
			}
		})
	}
}

package parser

import (
	"os"
	"testing"
)

func TestParser_Parse_SimpleProfile(t *testing.T) {
	file, err := os.Open("../../testdata/fixtures/simple.profile")
	if err != nil {
		t.Fatalf("Failed to open test file: %v", err)
	}
	defer file.Close()

	parser := NewParser()
	profile, err := parser.Parse(file)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if profile == nil {
		t.Fatal("Expected profile to be non-nil")
	}

	// Check scripts
	if len(profile.Scripts) != 1 {
		t.Errorf("Expected 1 script, got %d", len(profile.Scripts))
	}

	script, exists := profile.Scripts["/path/to/simple.vim"]
	if !exists {
		t.Fatal("Expected script /path/to/simple.vim to exist")
	}

	// Check script has lines
	if len(script.Lines) == 0 {
		t.Error("Expected script to have lines")
	}

	// Check functions
	if len(profile.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(profile.Functions))
	}

	fn := profile.Functions[0]
	if fn.Name != "SimpleFunc()" {
		t.Errorf("Expected function name SimpleFunc(), got %s", fn.Name)
	}
	if fn.Count != 2 {
		t.Errorf("Expected function count 2, got %d", fn.Count)
	}
}

func TestParser_ParseFile(t *testing.T) {
	parser := NewParser()
	profile, err := parser.ParseFile("../../testdata/fixtures/simple.profile")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if profile == nil {
		t.Fatal("Expected profile to be non-nil")
	}

	if len(profile.Scripts) == 0 {
		t.Error("Expected at least one script")
	}
}

func TestParser_ParseFile_NonExistent(t *testing.T) {
	parser := NewParser()
	_, err := parser.ParseFile("nonexistent.profile")

	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestParser_Parse_MultipleScripts(t *testing.T) {
	// Test with profile containing multiple scripts
	input := `SCRIPT  /path/to/first.vim
Sourced 1 time
Total time:   0.000010
 Self time:   0.000010

count  total (s)   self (s)
    1              0.000005 let g:x = 1

SCRIPT  /path/to/second.vim
Sourced 1 time
Total time:   0.000020
 Self time:   0.000020

count  total (s)   self (s)
    1              0.000010 let g:y = 2`

	parser := NewParser()
	profile, err := parser.ParseString(input)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(profile.Scripts) != 2 {
		t.Errorf("Expected 2 scripts, got %d", len(profile.Scripts))
	}

	if _, exists := profile.Scripts["/path/to/first.vim"]; !exists {
		t.Error("Expected first.vim to exist")
	}
	if _, exists := profile.Scripts["/path/to/second.vim"]; !exists {
		t.Error("Expected second.vim to exist")
	}
}

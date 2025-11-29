package parser

import (
	"testing"
	"time"
)

func TestLine_New(t *testing.T) {
	line := &Line{
		Number:    10,
		Count:     5,
		TotalTime: 100 * time.Microsecond,
		SelfTime:  50 * time.Microsecond,
	}

	if line.Number != 10 {
		t.Errorf("Expected Number to be 10, got %d", line.Number)
	}
	if line.Count != 5 {
		t.Errorf("Expected Count to be 5, got %d", line.Count)
	}
	if line.TotalTime != 100*time.Microsecond {
		t.Errorf("Expected TotalTime to be 100µs, got %v", line.TotalTime)
	}
	if line.SelfTime != 50*time.Microsecond {
		t.Errorf("Expected SelfTime to be 50µs, got %v", line.SelfTime)
	}
}

func TestFunction_New(t *testing.T) {
	lines := []Line{
		{Number: 1, Count: 2, SelfTime: 10 * time.Microsecond},
		{Number: 2, Count: 2, SelfTime: 20 * time.Microsecond},
	}

	fn := &Function{
		Name:      "TestFunc",
		Defined:   "/path/to/test.vim",
		StartLine: 10,
		EndLine:   12,
		Count:     2,
		TotalTime: 100 * time.Microsecond,
		SelfTime:  50 * time.Microsecond,
		Lines:     lines,
	}

	if fn.Name != "TestFunc" {
		t.Errorf("Expected Name to be TestFunc, got %s", fn.Name)
	}
	if fn.Count != 2 {
		t.Errorf("Expected Count to be 2, got %d", fn.Count)
	}
	if len(fn.Lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(fn.Lines))
	}
}

func TestScript_New(t *testing.T) {
	script := &Script{
		Path:      "/path/to/script.vim",
		Lines:     make(map[int]*Line),
		Functions: make(map[string]*Function),
	}

	script.Lines[1] = &Line{Number: 1, Count: 5}
	script.Lines[2] = &Line{Number: 2, Count: 3}
	script.Functions["TestFunc"] = &Function{Name: "TestFunc"}

	if script.Path != "/path/to/script.vim" {
		t.Errorf("Expected Path to be /path/to/script.vim, got %s", script.Path)
	}
	if len(script.Lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(script.Lines))
	}
	if len(script.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(script.Functions))
	}
}

func TestProfile_New(t *testing.T) {
	profile := &Profile{
		Scripts:   make(map[string]*Script),
		Functions: []*Function{},
	}

	script := &Script{Path: "/path/to/script.vim"}
	profile.Scripts[script.Path] = script
	profile.Functions = append(profile.Functions, &Function{Name: "Func1"})

	if len(profile.Scripts) != 1 {
		t.Errorf("Expected 1 script, got %d", len(profile.Scripts))
	}
	if len(profile.Functions) != 1 {
		t.Errorf("Expected 1 function, got %d", len(profile.Functions))
	}
}

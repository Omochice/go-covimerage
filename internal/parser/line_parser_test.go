package parser

import (
	"testing"
	"time"
)

func TestParser_ParseLineData(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantNumber    int
		wantCount     int
		wantTotalTime time.Duration
		wantSelfTime  time.Duration
		wantError     bool
	}{
		{
			name:         "line with count and self time",
			input:        "    5              0.000020 let g:count = 0",
			wantNumber:   -1, // Will be set separately
			wantCount:    5,
			wantTotalTime: 0,
			wantSelfTime: 20 * time.Microsecond,
			wantError:    false,
		},
		{
			name:          "line with count and total time",
			input:         "    5   0.000030             for i in range(5)",
			wantNumber:    -1,
			wantCount:     5,
			wantTotalTime: 30 * time.Microsecond,
			wantSelfTime:  0,
			wantError:     false,
		},
		{
			name:          "line with count, total and self time",
			input:         "    2   0.000015   0.000010   call SomeFunc()",
			wantNumber:    -1,
			wantCount:     2,
			wantTotalTime: 15 * time.Microsecond,
			wantSelfTime:  10 * time.Microsecond,
			wantError:     false,
		},
		{
			name:          "line not executed",
			input:         "                            \" Comment line",
			wantNumber:    -1,
			wantCount:     0,
			wantTotalTime: 0,
			wantSelfTime:  0,
			wantError:     false,
		},
		{
			name:      "header line",
			input:     "count  total (s)   self (s)",
			wantError: true, // Should skip header
		},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line, err := parser.parseLineData(tt.input, 1)

			if tt.wantError {
				if err == nil {
					t.Error("Expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if line.Count != tt.wantCount {
				t.Errorf("Count = %d, want %d", line.Count, tt.wantCount)
			}
			if line.TotalTime != tt.wantTotalTime {
				t.Errorf("TotalTime = %v, want %v", line.TotalTime, tt.wantTotalTime)
			}
			if line.SelfTime != tt.wantSelfTime {
				t.Errorf("SelfTime = %v, want %v", line.SelfTime, tt.wantSelfTime)
			}
		})
	}
}

func TestParser_ParseTime(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
	}{
		{"0.000010", 10 * time.Microsecond},
		{"0.000100", 100 * time.Microsecond},
		{"0.001000", 1000 * time.Microsecond},
		{"0.010000", 10000 * time.Microsecond},
		{"1.000000", 1 * time.Second},
		{"", 0},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parser.parseTime(tt.input)
			if result != tt.want {
				t.Errorf("parseTime(%q) = %v, want %v", tt.input, result, tt.want)
			}
		})
	}
}

func TestParser_ParseCount(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"1", 1},
		{"5", 5},
		{"100", 100},
		{"", 0},
		{"   ", 0},
	}

	parser := NewParser()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parser.parseCount(tt.input)
			if result != tt.want {
				t.Errorf("parseCount(%q) = %d, want %d", tt.input, result, tt.want)
			}
		})
	}
}

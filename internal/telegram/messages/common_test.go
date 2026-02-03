package messages

import (
	"testing"
)

func TestProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		current  int
		total    int
		width    int
		expected string
	}{
		{
			name:     "empty progress",
			current:  0,
			total:    10,
			width:    10,
			expected: "░░░░░░░░░░",
		},
		{
			name:     "half progress",
			current:  5,
			total:    10,
			width:    10,
			expected: "▓▓▓▓▓░░░░░",
		},
		{
			name:     "full progress",
			current:  10,
			total:    10,
			width:    10,
			expected: "▓▓▓▓▓▓▓▓▓▓",
		},
		{
			name:     "one third",
			current:  1,
			total:    3,
			width:    9,
			expected: "▓▓▓░░░░░░",
		},
		{
			name:     "two thirds",
			current:  2,
			total:    3,
			width:    9,
			expected: "▓▓▓▓▓▓░░░",
		},
		{
			name:     "negative current clamps to zero",
			current:  -5,
			total:    10,
			width:    10,
			expected: "░░░░░░░░░░",
		},
		{
			name:     "current exceeds total clamps to full",
			current:  15,
			total:    10,
			width:    10,
			expected: "▓▓▓▓▓▓▓▓▓▓",
		},
		{
			name:     "zero total returns empty string",
			current:  5,
			total:    0,
			width:    10,
			expected: "",
		},
		{
			name:     "negative total returns empty string",
			current:  5,
			total:    -5,
			width:    10,
			expected: "",
		},
		{
			name:     "zero width returns empty string",
			current:  5,
			total:    10,
			width:    0,
			expected: "",
		},
		{
			name:     "small width",
			current:  3,
			total:    4,
			width:    4,
			expected: "▓▓▓░",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProgressBar(tt.current, tt.total, tt.width)
			if result != tt.expected {
				t.Errorf("ProgressBar(%d, %d, %d) = %q, want %q",
					tt.current, tt.total, tt.width, result, tt.expected)
			}
		})
	}
}

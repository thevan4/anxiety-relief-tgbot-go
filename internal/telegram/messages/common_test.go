package messages

import "testing"

func TestTimerCountdown(t *testing.T) {
	tests := []struct {
		name     string
		elapsed  int
		total    int
		expected string
	}{
		{"start", 0, 5, "⏱️ 5"},
		{"middle", 2, 5, "⏱️ 3"},
		{"almost done", 4, 5, "⏱️ 1"},
		{"done", 5, 5, "✨"},
		{"over", 6, 5, "✨"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TimerCountdown(tt.elapsed, tt.total)
			if result != tt.expected {
				t.Errorf("TimerCountdown(%d, %d) = %q, want %q",
					tt.elapsed, tt.total, result, tt.expected)
			}
		})
	}
}

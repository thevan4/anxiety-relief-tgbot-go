package messages

import "testing"

func TestTimerCountdown(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := TimerCountdown(tt.elapsed, tt.total)
			if result != tt.expected {
				t.Errorf("TimerCountdown(%d, %d) = %q, want %q",
					tt.elapsed, tt.total, result, tt.expected)
			}
		})
	}
}

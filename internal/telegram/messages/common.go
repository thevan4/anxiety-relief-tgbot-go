// Package messages provides text formatting functions for bot responses.
package messages

import "fmt"

// TimerCountdown returns emoji timer with remaining seconds.
// elapsed - seconds passed
// total - total seconds
// Returns format like "⏱️ 4" or "✨" when done.
func TimerCountdown(elapsed, total int) string {
	remaining := total - elapsed
	if remaining <= 0 {
		return "✨"
	}
	return fmt.Sprintf("⏱️ %d", remaining)
}

// Package sessioncore provides platform-agnostic session management primitives.
// No external dependencies (Telegram, Redis, etc.) - pure Go stdlib.
package sessioncore

import "time"

// CleanupRetry stores retry information for cleanup decisions.
type CleanupRetry struct {
	State   string
	Attempt int
}

// CleanupDecision represents the result of cleanup logic evaluation.
type CleanupDecision int

const (
	// DecisionDelete means the resource should be deleted.
	DecisionDelete CleanupDecision = iota
	// DecisionKeep means the resource should be kept (user became active).
	DecisionKeep
	// DecisionRetry means we should check again later.
	DecisionRetry
)

// CleanupResult contains the decision and optional retry info.
type CleanupResult struct {
	Decision  CleanupDecision
	Retry     *CleanupRetry // set if Decision == DecisionRetry
	NextCheck time.Time     // when to check next (for DecisionRetry)
	Reason    string        // human-readable reason for logging
}

// Config holds timing configuration for session management.
type Config struct {
	// StateTTL is how long state lives without refresh (activity indicator).
	StateTTL time.Duration

	// ResourceTTL is max lifetime of managed resource (e.g., 48h for Telegram messages).
	ResourceTTL time.Duration

	// CleanupDelay is time before first cleanup check after resource creation.
	CleanupDelay time.Duration

	// RetryInterval is time between retry checks.
	RetryInterval time.Duration

	// MaxRetries is maximum retry attempts before force cleanup.
	MaxRetries int
}

// DefaultConfig returns sensible defaults for Telegram-like platforms.
func DefaultConfig() Config {
	return Config{
		StateTTL:      5 * time.Minute,
		ResourceTTL:   48 * time.Hour,
		CleanupDelay:  47 * time.Hour,
		RetryInterval: 5 * time.Minute,
		MaxRetries:    6,
	}
}

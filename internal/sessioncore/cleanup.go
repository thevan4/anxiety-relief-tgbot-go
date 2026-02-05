package sessioncore

import "time"

// CleanupProcessor contains the core cleanup decision logic.
// It decides whether to delete, keep, or retry based on user state.
// No platform dependencies - pure logic.
type CleanupProcessor struct {
	config Config
}

// NewCleanupProcessor creates a new cleanup processor with given config.
func NewCleanupProcessor(config Config) *CleanupProcessor {
	return &CleanupProcessor{config: config}
}

// EvaluateCleanup determines what to do with a user's resource.
// Arguments:
//   - resourceCreatedAt: when the resource was created
//   - currentState: user's current state (empty string = inactive)
//   - retry: previous retry info (nil = first check)
//   - now: current time
//
// Returns CleanupResult with decision and metadata.
func (p *CleanupProcessor) EvaluateCleanup(
	resourceCreatedAt time.Time,
	currentState string,
	retry *CleanupRetry,
	now time.Time,
) CleanupResult {
	// Check 1: Resource too old (past platform limit)
	if !resourceCreatedAt.IsZero() {
		age := now.Sub(resourceCreatedAt)
		if age > p.config.ResourceTTL {
			return CleanupResult{
				Decision: DecisionDelete,
				Reason:   "resource exceeded TTL",
			}
		}
	}

	// Check 2: No resource creation time - can't determine age, clean up
	if resourceCreatedAt.IsZero() {
		return CleanupResult{
			Decision: DecisionDelete,
			Reason:   "no resource creation timestamp",
		}
	}

	// Check 3: First check (no retry exists)
	if retry == nil {
		return p.handleFirstCheck(currentState, now)
	}

	// Check 4: Retry check
	return p.handleRetryCheck(currentState, retry, now)
}

// handleFirstCheck handles the first cleanup evaluation for a user.
func (p *CleanupProcessor) handleFirstCheck(currentState string, now time.Time) CleanupResult {
	// User NOT active (state empty)
	if currentState == "" {
		return CleanupResult{
			Decision: DecisionDelete,
			Reason:   "user inactive on first check",
		}
	}

	// User IS active — schedule retry
	return CleanupResult{
		Decision: DecisionRetry,
		Retry: &CleanupRetry{
			State:   currentState,
			Attempt: 1,
		},
		NextCheck: now.Add(p.config.RetryInterval),
		Reason:    "user active, scheduling retry",
	}
}

// handleRetryCheck handles retry cleanup evaluations.
func (p *CleanupProcessor) handleRetryCheck(currentState string, retry *CleanupRetry, now time.Time) CleanupResult {
	// User EXITED activity (state became empty)
	if currentState == "" {
		// User was active, finished — DON'T delete, they might still use it
		return CleanupResult{
			Decision: DecisionKeep,
			Reason:   "user finished activity, keeping resource",
		}
	}

	// User in SAME state as before
	if currentState == retry.State {
		// Stuck in same state for too long — delete
		return CleanupResult{
			Decision: DecisionDelete,
			Reason:   "user stuck in same state",
		}
	}

	// User in NEW/DIFFERENT state (actively using)
	if retry.Attempt < p.config.MaxRetries {
		// Still active, reschedule
		return CleanupResult{
			Decision: DecisionRetry,
			Retry: &CleanupRetry{
				State:   currentState,
				Attempt: retry.Attempt + 1,
			},
			NextCheck: now.Add(p.config.RetryInterval),
			Reason:    "user switched states, scheduling retry",
		}
	}

	// Max retries reached — force delete
	return CleanupResult{
		Decision: DecisionDelete,
		Reason:   "max retries exceeded",
	}
}

// ScheduleInitialCleanup returns the time when first cleanup check should happen.
func (p *CleanupProcessor) ScheduleInitialCleanup(resourceCreatedAt time.Time) time.Time {
	return resourceCreatedAt.Add(p.config.CleanupDelay)
}

// Config returns the current configuration.
func (p *CleanupProcessor) Config() Config {
	return p.config
}

package session

import "time"

// NeedsRecreate checks if resource is too old to edit (48h Telegram limit).
// Returns true if resource must be deleted and recreated instead of edited.
func NeedsRecreate(createdAt time.Time, ttl time.Duration) bool {
	if createdAt.IsZero() {
		return true // No creation time — assume old, recreate
	}
	return time.Since(createdAt) >= ttl
}

// GetTransitionOp determines the operation for transitioning between screens.
// This encapsulates the business logic of when to edit vs recreate.
//
// Rules:
//   - Unknown source: Send (no existing message to modify)
//   - From/to Holder: Recreate (holder is a separate message structure)
//   - Everything else: Edit (same message, different content)
func GetTransitionOp(from, to Screen) TransitionOp {
	if from == ScreenUnknown {
		return OpSend
	}
	if from == ScreenHolder || to == ScreenHolder {
		return OpRecreate
	}
	return OpEdit
}

// ShouldRecreateForAge checks if edit should become recreate due to age.
// Use this after GetTransitionOp returns OpEdit to check if age forces recreate.
func ShouldRecreateForAge(op TransitionOp, createdAt time.Time, ttl time.Duration) TransitionOp {
	if op == OpEdit && NeedsRecreate(createdAt, ttl) {
		return OpRecreate
	}
	return op
}

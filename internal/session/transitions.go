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
//   - Holder → Menu: always Recreate (different message structure)
//   - Menu → ExerciseIntro: Edit (same message)
//   - Menu → Menu: Edit
//   - ExerciseComplete → Menu: Edit
//   - To LanguageSelect: Edit
//   - From LanguageSelect: Edit
//   - Unknown source: Send
func GetTransitionOp(from, to Screen) TransitionOp {
	// Unknown source — need to send new message
	if from == ScreenUnknown {
		return OpSend
	}

	// Holder to anything else — recreate (holder is separate message)
	if from == ScreenHolder {
		return OpRecreate
	}

	// Same screen — edit
	if from == to {
		return OpEdit
	}

	// Menu transitions
	if from == ScreenMenu {
		switch to {
		case ScreenExerciseIntro, ScreenExerciseRunning, ScreenLanguageSelect:
			return OpEdit
		case ScreenHolder:
			return OpRecreate
		}
	}

	// Exercise transitions
	if from == ScreenExerciseIntro || from == ScreenExerciseRunning || from == ScreenExerciseComplete {
		switch to {
		case ScreenMenu, ScreenExerciseIntro, ScreenExerciseRunning, ScreenExerciseComplete, ScreenLanguageSelect:
			return OpEdit
		case ScreenHolder:
			return OpRecreate
		}
	}

	// Language select transitions
	if from == ScreenLanguageSelect {
		switch to {
		case ScreenMenu, ScreenExerciseIntro, ScreenExerciseRunning, ScreenExerciseComplete:
			return OpEdit
		case ScreenHolder:
			return OpRecreate
		}
	}

	// Error screen — can always be edited to anything
	if from == ScreenError {
		if to == ScreenHolder {
			return OpRecreate
		}
		return OpEdit
	}

	// Default: edit
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

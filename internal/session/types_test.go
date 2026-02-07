package session

import "testing"

func TestStateIsRunning(t *testing.T) {
	t.Parallel()
	tests := []struct {
		state State
		want  bool
	}{
		{StateUnknown, false},
		{StateIdle, false},
		{StateBreathingActive, false},
		{StateBreathingRunning, true},
		{StateGroundingStep1, false},
		{StateGroundingStep2, false},
		{StateGroundingStep3, false},
		{StateGroundingStep4, false},
		{StateGroundingStep5, false},
		{StateGuidedBreathingSelect, false},
		{StateGuidedBreathingRunning, true},
		{StatePMRActive, false},
		{StatePMRRunning, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(string(tt.state), func(t *testing.T) {
			t.Parallel()
			if got := tt.state.IsRunning(); got != tt.want {
				t.Errorf("State(%q).IsRunning() = %v, want %v", tt.state, got, tt.want)
			}
		})
	}
}

func TestStateIsActive(t *testing.T) {
	t.Parallel()
	tests := []struct {
		state State
		want  bool
	}{
		{StateUnknown, false},
		{StateIdle, false},
		{StateBreathingActive, true},
		{StateBreathingRunning, true},
		{StateGroundingStep1, true},
		{StateGroundingStep2, true},
		{StateGroundingStep3, true},
		{StateGroundingStep4, true},
		{StateGroundingStep5, true},
		{StateGuidedBreathingSelect, true},
		{StateGuidedBreathingRunning, true},
		{StatePMRActive, true},
		{StatePMRRunning, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(string(tt.state), func(t *testing.T) {
			t.Parallel()
			if got := tt.state.IsActive(); got != tt.want {
				t.Errorf("State(%q).IsActive() = %v, want %v", tt.state, got, tt.want)
			}
		})
	}
}

func TestStateConstants(t *testing.T) {
	t.Parallel()
	// Verify state constants are unique
	states := []State{
		StateUnknown,
		StateIdle,
		StateBreathingActive,
		StateBreathingRunning,
		StateGroundingStep1,
		StateGroundingStep2,
		StateGroundingStep3,
		StateGroundingStep4,
		StateGroundingStep5,
		StateGuidedBreathingSelect,
		StateGuidedBreathingRunning,
		StatePMRActive,
		StatePMRRunning,
	}

	seen := make(map[State]bool)
	for _, s := range states {
		if s != StateUnknown && seen[s] {
			t.Errorf("duplicate state: %q", s)
		}
		seen[s] = true
	}
}

func TestScreenString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		screen Screen
		want   string
	}{
		{ScreenUnknown, "Unknown"},
		{ScreenHolder, "Holder"},
		{ScreenMenu, "Menu"},
		{ScreenExerciseIntro, "ExerciseIntro"},
		{ScreenExerciseRunning, "ExerciseRunning"},
		{ScreenExerciseComplete, "ExerciseComplete"},
		{ScreenLanguageSelect, "LanguageSelect"},
		{ScreenError, "Error"},
		{Screen(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			if got := tt.screen.String(); got != tt.want {
				t.Errorf("Screen(%d).String() = %q, want %q", tt.screen, got, tt.want)
			}
		})
	}
}

func TestTransitionOpString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		op   TransitionOp
		want string
	}{
		{OpEdit, "Edit"},
		{OpRecreate, "Recreate"},
		{OpSend, "Send"},
		{TransitionOp(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			if got := tt.op.String(); got != tt.want {
				t.Errorf("TransitionOp(%d).String() = %q, want %q", tt.op, got, tt.want)
			}
		})
	}
}

func TestCleanupDecisionString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		decision CleanupDecision
		want     string
	}{
		{DecisionDelete, "Delete"},
		{DecisionKeep, "Keep"},
		{DecisionRetry, "Retry"},
		{CleanupDecision(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()
			if got := tt.decision.String(); got != tt.want {
				t.Errorf("CleanupDecision(%d).String() = %q, want %q", tt.decision, got, tt.want)
			}
		})
	}
}

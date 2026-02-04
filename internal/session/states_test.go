package session

import "testing"

func TestStateIsRunning(t *testing.T) {
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
		{StatePMRTense, false},
		{StatePMRRelax, false},
		{StatePMRRunning, true},
		{StateThoughtLabelingActive, false},
		{StateThoughtLabelingInput, false},
		{StateThoughtLabelingCategory, false},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if got := tt.state.IsRunning(); got != tt.want {
				t.Errorf("State(%q).IsRunning() = %v, want %v", tt.state, got, tt.want)
			}
		})
	}
}

func TestStateIsActive(t *testing.T) {
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
		{StatePMRTense, true},
		{StatePMRRelax, true},
		{StatePMRRunning, true},
		{StateThoughtLabelingActive, true},
		{StateThoughtLabelingInput, true},
		{StateThoughtLabelingCategory, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			if got := tt.state.IsActive(); got != tt.want {
				t.Errorf("State(%q).IsActive() = %v, want %v", tt.state, got, tt.want)
			}
		})
	}
}

func TestStateConstants(t *testing.T) {
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
		StatePMRTense,
		StatePMRRelax,
		StatePMRRunning,
		StateThoughtLabelingActive,
		StateThoughtLabelingInput,
		StateThoughtLabelingCategory,
	}

	seen := make(map[State]bool)
	for _, s := range states {
		if s != StateUnknown && seen[s] {
			t.Errorf("duplicate state: %q", s)
		}
		seen[s] = true
	}
}

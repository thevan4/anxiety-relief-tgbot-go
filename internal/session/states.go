package session

// State represents user's current session state in the bot flow.
type State string

const (
	// StateUnknown represents no active session.
	StateUnknown State = ""
	// StateIdle represents idle state when user is not in any exercise.
	StateIdle State = "idle"
	// StateBreathingActive represents breathing technique menu is active.
	StateBreathingActive State = "breathing_active"
	// StateBreathingRunning represents breathing exercise is currently running.
	StateBreathingRunning State = "breathing_running"
	// StateGroundingStep1 represents grounding step 1 (see 5 things).
	StateGroundingStep1 State = "grounding_step_1"
	// StateGroundingStep2 represents grounding step 2 (touch 4 things).
	StateGroundingStep2 State = "grounding_step_2"
	// StateGroundingStep3 represents grounding step 3 (hear 3 sounds).
	StateGroundingStep3 State = "grounding_step_3"
	// StateGroundingStep4 represents grounding step 4 (smell 2 things).
	StateGroundingStep4 State = "grounding_step_4"
	// StateGroundingStep5 represents grounding step 5 (taste 1 thing).
	StateGroundingStep5 State = "grounding_step_5"

	// StateGuidedBreathingSelect indicates user is selecting breathing pattern.
	StateGuidedBreathingSelect State = "guided_breathing_select"
	// StateGuidedBreathingRunning represents guided breathing exercise is running.
	StateGuidedBreathingRunning State = "guided_breathing_running"

	// StatePMRActive indicates PMR technique menu is active.
	StatePMRActive State = "pmr_active"
	// StatePMRTense represents PMR tense phase.
	StatePMRTense State = "pmr_tense"
	// StatePMRRelax represents PMR relax phase.
	StatePMRRelax State = "pmr_relax"
	// StatePMRRunning represents PMR exercise is running.
	StatePMRRunning State = "pmr_running"

	// StateThoughtLabelingActive indicates thought labeling menu is active.
	StateThoughtLabelingActive State = "thought_labeling_active"
	// StateThoughtLabelingInput represents user is entering thought text.
	StateThoughtLabelingInput State = "thought_labeling_input"
	// StateThoughtLabelingCategory represents user is selecting thought category.
	StateThoughtLabelingCategory State = "thought_labeling_category"

	// StateVisualizationSelect indicates user is selecting visualization scene.
	StateVisualizationSelect State = "visualization_select"
	// StateVisualizationRunning represents visualization exercise is running.
	StateVisualizationRunning State = "visualization_running"
)

// IsRunning returns true if user is in an active exercise that should not be interrupted.
func (s State) IsRunning() bool {
	switch s {
	case StateBreathingRunning,
		StateGuidedBreathingRunning,
		StatePMRRunning,
		StateVisualizationRunning:
		return true
	}
	return false
}

// IsActive returns true if user has any active session (including selection screens).
func (s State) IsActive() bool {
	return s != StateUnknown && s != StateIdle
}
